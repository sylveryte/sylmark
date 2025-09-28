package lspserver

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"sylmark/data"
	"sylmark/lsp"
	"sylmark/utils"
	"time"

	"github.com/sourcegraph/jsonrpc2"
	tree_sitter_markdown "github.com/sylveryte/tree-sitter-markdown/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type ServerDebouncers = struct {
	DocumentDidChange   *utils.SylDebouncer
	SemantickTokensFull *utils.SylDebouncer
}

type LangHandler struct {
	Parser       *tree_sitter.Parser
	InlineParser *tree_sitter.Parser
	Store        data.Store
	Debouncers   *ServerDebouncers
	Config       data.Config
	Connection   *jsonrpc2.Conn
}

func NewHandler() (hanlder *LangHandler) {
	return &LangHandler{
		Store:  data.NewStore(),
		Config: data.NewConfig(),
		Debouncers: &ServerDebouncers{
			DocumentDidChange:   utils.NewSylDebouncer(300 * time.Millisecond),
			SemantickTokensFull: utils.NewSylDebouncer(400 * time.Millisecond),
		},
	}
}

func (h *LangHandler) addRootPathAndLoad(dir string) {
	h.Config.RootPath = dir
	t := time.Now()
	h.loadAllClosedDocsData()
	slog.Info(fmt.Sprintf("=====Load time is %dms",time.Since(t).Milliseconds()))
	h.Config.CreatDirsIfNeeded()
}

func (h *LangHandler) loadAllClosedDocsData() {
	if h.Config.RootPath == "" {
		slog.Error("h.rootPath is empty")
		return
	}

	filepath.WalkDir(h.Config.RootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && (strings.HasSuffix(path, ".git") || strings.HasSuffix(path, "node_modules")) {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			h.loadDocData(path)
		}
		return nil
	})
}

func (h *LangHandler) loadDocData(mdDocPath string) {
	// using directly ContentFromDocPath to skip caching in store
	content := data.ContentFromDocPath(mdDocPath)
	uri, err := data.UriFromPath(mdDocPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	trees := h.parse(content, nil)
	defer trees[0].Close()
	defer trees[1].Close()
	h.Store.LoadData(uri, content, trees)
}

func (h *LangHandler) onDocCreated(uri lsp.DocumentURI, content string) {
	h.onDocOpened(uri, content)
	docPath, _ := data.PathFromURI(uri)
	h.loadDocData(docPath)
}

func (h *LangHandler) onDocDeleted(uri lsp.DocumentURI) {
	docData, ok := h.Store.GetDocMustTree(uri, h.parse)
	if ok {
		h.Store.UnloadData(uri, string(docData.Content), docData.Trees)
		h.Store.RemoveDoc(uri)
	}
}
func (h *LangHandler) onDocOpened(uri lsp.DocumentURI, content string) {
	trees := h.parse(content, nil)
	doc := data.Document(content)

	slog.Info("First main---------------")
	lsp.PrintTsTree(*trees.GetMainTree().RootNode(), 0, content)
	slog.Info("Now inline-------------")
	lsp.PrintTsTree(*trees.GetInlineTree().RootNode(), 0, content)
	docData := data.NewDocumentData(doc, trees)
	h.Store.AddUpdateDoc(uri, docData)
	docData.Headings = h.Store.GetHeadingWithinDataStore(uri, h.parse)
	h.Store.AddUpdateDoc(uri, docData)
}

func (h *LangHandler) onDocChanged(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent) {
	h.Store.SyncChangedDocument(uri, changes, h.parse)
}

func (h *LangHandler) SetupGrammars() {
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_markdown.Language())
	parser.SetLanguage(language)

	inlineParser := tree_sitter.NewParser()
	inlineLanguage := tree_sitter.NewLanguage(tree_sitter_markdown.InlineLanguage())
	inlineParser.SetLanguage(inlineLanguage)

	h.Parser = parser
	h.InlineParser = inlineParser
}

func (h *LangHandler) DocAndNodeFromURIAndPosition(uri lsp.DocumentURI, position lsp.Position, parse lsp.ParseFunction) (doc data.Document, node *tree_sitter.Node, ok bool) {
	docData, ok := h.Store.GetDocMustTree(uri, parse)
	if !ok {
		slog.Error("Document missing" + string(uri))
		return "", nil, false
	}
	point := lsp.PointFromPosition(position)

	doc = docData.Content
	// syltodo handle treees
	node = docData.Trees.GetMainTree().RootNode().NamedDescendantForPointRange(point, point)
	if node.Parent().Kind() == "atx_heading" {
		node = node.Parent()
		return
	}
	if node.Kind() == "inline" || node.Kind() == "paragraph" {
		node = docData.Trees.GetInlineTree().RootNode().NamedDescendantForPointRange(point, point)
	}

	ok = true
	return
}

func (h *LangHandler) parse(content string, oldTrees *lsp.Trees) *lsp.Trees {
	var trees lsp.Trees
	trees[0] = h.Parser.Parse([]byte(content), nil)
	trees[1] = h.InlineParser.Parse([]byte(content), nil)
	return &trees
}

func (h *LangHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	slog.Info("-------------reqmethod=> " + req.Method)
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, conn, req)
	case "initialized":
		return
	case "shutdown":
		return h.handleShutdown(ctx, conn, req)
	case "textDocument/didOpen":
		return h.handleTextDocumentDidOpen(ctx, conn, req)
	// case "textDocument/didClose":
	// 	return h.handleTextDocumentDidClose(ctx, conn, req)
	case "textDocument/didChange":
		return h.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/hover":
		return h.handleHover(ctx, conn, req)
	case "textDocument/completion":
		return h.handleTextDocumentCompletion(ctx, conn, req)
	case "textDocument/references":
		return h.handleTextDocumentReferences(ctx, conn, req)
	case "textDocument/definition":
		return h.handleTextDocumentDefinition(ctx, conn, req)
	case "textDocument/semanticTokens/full":
		return h.handleTextDocumentSemanticTokensFull(ctx, conn, req)
	case "textDocument/codeAction":
		return h.handleCodeAction(ctx, conn, req)
	case "textDocument/diagnostic":
		return h.handleDiagnostics(ctx, conn, req)
	case "workspace/executeCommand":
		return h.handleWorkspaceExecuteCommand(ctx, conn, req)
	case "workspace/didDeleteFiles":
		return h.handleWorkspaceDidDeleteFiles(ctx, conn, req)
	case "workspace/didCreateFiles":
		return h.handleWorkspaceDidCreateFiles(ctx, conn, req)
	case "workspace/didRenameFiles":
		return h.handleWorkspaceDidRenameFiles(ctx, conn, req)
	case "workspace/symbol":
		return h.handleWorkspaceSymbol(ctx, conn, req)
	}
	return nil, nil
}
