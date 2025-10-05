package lspserver

import (
	"context"
	"fmt"
	"log/slog"
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
	Connection   *jsonrpc2.Conn
}

func NewHandler() (hanlder *LangHandler) {
	return &LangHandler{
		Store: data.NewStore(),
		Debouncers: &ServerDebouncers{
			DocumentDidChange:   utils.NewSylDebouncer(300 * time.Millisecond),
			SemantickTokensFull: utils.NewSylDebouncer(400 * time.Millisecond),
		},
	}
}

func (h *LangHandler) loadDocData(mdDocPath string) {
	uri, content, trees, err := TreesFromMdDocPath(mdDocPath, h.parse)
	if err != nil {
		return
	}
	// using directly ContentFromDocPath to skip caching in store
	defer trees[0].Close()
	defer trees[1].Close()
	h.Store.LoadData(uri, content, trees)
}

func (h *LangHandler) onDocCreated(uri lsp.DocumentURI, content string) {
	h.onDocOpened(uri, content)
	docPath, _ := data.PathFromURI(uri)
	h.loadDocData(docPath)
}
func (h *LangHandler) onDocRenamed(param lsp.FileRename) {
	docData, ok := h.Store.GetDocMustTree(param.OldUri, h.parse)
	var content string
	if ok {
		content = string(docData.Content)
	}
	h.onDocDeleted(param.OldUri)
	h.onDocCreated(param.NewUri, content)

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

	// slog.Info("First main---------------")
	// lsp.PrintTsTree(*trees.GetMainTree().RootNode(), 0, content)
	// slog.Info("Now inline-------------")
	// lsp.PrintTsTree(*trees.GetInlineTree().RootNode(), 0, content)
	docData := data.NewDocumentData(doc, trees)
	h.Store.AddUpdateDoc(uri, docData)
	docData.Headings = h.Store.GetLoadedDataStore(uri, h.parse)
	docData.FootNotes = h.Store.GetLoadedFootNotesStore(uri, h.parse)
	h.Store.AddUpdateDoc(uri, docData)
}

func (h *LangHandler) onDocChanged(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent) {
	h.Store.SyncChangedDocument(uri, changes, h.parse)
}

func getParsers() [2]*tree_sitter.Parser {
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_markdown.Language())
	parser.SetLanguage(language)

	inlineParser := tree_sitter.NewParser()
	inlineLanguage := tree_sitter.NewLanguage(tree_sitter_markdown.InlineLanguage())
	inlineParser.SetLanguage(inlineLanguage)

	parsers := [2]*tree_sitter.Parser{
		parser,
		inlineParser,
	}

	return parsers
}

func (h *LangHandler) SetupGrammars() {
	parsers := getParsers()

	h.Parser = parsers[0]
	h.InlineParser = parsers[1]
}

func (h *LangHandler) DocAndNodeFromURIAndPosition(uri lsp.DocumentURI, position lsp.Position, parse lsp.ParseFunction) (doc data.Document, node *tree_sitter.Node, ok bool) {
	docData, ok := h.Store.GetDocMustTree(uri, parse)
	if !ok {
		slog.Error("Document missing" + string(uri))
		return "", nil, false
	}
	point := lsp.PointFromPosition(position)

	doc = docData.Content
	node = docData.Trees.GetMainTree().RootNode().NamedDescendantForPointRange(point, point)
	if node.Parent().Kind() == "atx_heading" {
		node = node.Parent()
		return
	}
	if lsp.IsInlineParseNeeded(node) {
		node = docData.Trees.GetInlineTree().RootNode().NamedDescendantForPointRange(point, point)
	}

	ok = true
	return
}

// be sure to close the parsers and trees if able
func getParseFunction(parsers [2]*tree_sitter.Parser) lsp.ParseFunction {
	return func(content string, oldTrees *lsp.Trees) *lsp.Trees {
		var trees lsp.Trees
		trees[0] = parsers[0].Parse([]byte(content), nil)
		trees[1] = parsers[1].Parse([]byte(content), nil)
		return &trees
	}
}

// update getParseFunction too
func (h *LangHandler) parse(content string, oldTrees *lsp.Trees) *lsp.Trees {
	var trees lsp.Trees
	trees[0] = h.Parser.Parse([]byte(content), nil)
	trees[1] = h.InlineParser.Parse([]byte(content), nil)
	return &trees
}

func (h *LangHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	t := time.Now()
	switch req.Method {
	case "initialize":
		result, err = h.handleInitialize(ctx, conn, req)
	case "initialized":
	case "shutdown":
		result, err = h.handleShutdown(ctx, conn, req)
	case "textDocument/didOpen":
		result, err = h.handleTextDocumentDidOpen(ctx, conn, req)
	// case "textDocument/didClose":
	// 	result,err= h.handleTextDocumentDidClose(ctx, conn, req)
	case "textDocument/didChange":
		result, err = h.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/hover":
		result, err = h.handleHover(ctx, conn, req)
	case "textDocument/completion":
		result, err = h.handleTextDocumentCompletion(ctx, conn, req)
	case "textDocument/references":
		result, err = h.handleTextDocumentReferences(ctx, conn, req)
	case "textDocument/definition":
		result, err = h.handleTextDocumentDefinition(ctx, conn, req)
	case "textDocument/semanticTokens/full":
		result, err = h.handleTextDocumentSemanticTokensFull(ctx, conn, req)
	case "textDocument/codeAction":
		result, err = h.handleCodeAction(ctx, conn, req)
	case "textDocument/diagnostic":
		result, err = h.handleDiagnostics(ctx, conn, req)
	case "workspace/executeCommand":
		result, err = h.handleWorkspaceExecuteCommand(ctx, conn, req)
	case "workspace/didDeleteFiles":
		result, err = h.handleWorkspaceDidDeleteFiles(ctx, conn, req)
	case "workspace/didCreateFiles":
		result, err = h.handleWorkspaceDidCreateFiles(ctx, conn, req)
	case "workspace/didRenameFiles":
		result, err = h.handleWorkspaceDidRenameFiles(ctx, conn, req)
	case "workspace/symbol":
		result, err = h.handleWorkspaceSymbol(ctx, conn, req)
	}
	slog.Info(fmt.Sprintf("%dms<==%s", time.Since(t).Milliseconds(), req.Method))
	return result, err
}
