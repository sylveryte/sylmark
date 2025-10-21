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
	uri, content, trees, err := TreesFromUri(mdDocPath, h.parse)
	if err != nil {
		return
	}
	// using directly ContentFromDocPath to skip caching in store
	defer trees[0].Close()
	defer trees[1].Close()

	id := h.Store.GetIdFromURI(uri)
	h.Store.LoadData(id, content, trees)
}

func (h *LangHandler) onDocCreated(id data.Id, content string) {
	h.onDocOpened(id, content)
	uri, _ := h.Store.GetUri(id)
	docPath, _ := data.PathFromURI(uri)
	h.loadDocData(docPath)
}
func (h *LangHandler) onDocRenamed(param lsp.FileRename) {
	id := h.Store.GetIdFromURI(param.OldUri)
	// replace uri in idstore
	h.Store.IdStore.ReplaceUri(id, param.NewUri)
	oldTarget, _ := data.GetTarget(param.OldUri)
	newTarget, _ := data.GetTarget(param.NewUri)
	h.Store.ReplaceTarget(id, oldTarget, newTarget)
}
func (h *LangHandler) onDocDeleted(id data.Id) {
	docData, ok := h.Store.GetDocMustTree(id, h.parse)
	if ok {
		h.Store.UnloadData(id, string(docData.Content), docData.Trees)
		h.Store.RemoveDoc(id)
	}
}
func (h *LangHandler) onDocOpened(id data.Id, content string) {
	h.Store.UpdateAndReloadDoc(id, content, h.parse)
}

func (h *LangHandler) onDocChanged(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent) {
	id := h.Store.GetIdFromURI(uri)
	h.Store.SyncChangedDocument(id, changes, h.parse)
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

func (h *LangHandler) DocAndNodeFromURIAndPosition(id data.Id, position lsp.Position, parse lsp.ParseFunction) (docData data.DocumentData, node *tree_sitter.Node, ok bool) {
	docData, ok = h.Store.GetDocMustTree(id, parse)
	if !ok {
		slog.Error("Document missing" + string(id))
		return docData, nil, false
	}
	point := lsp.PointFromPosition(position)

	node = docData.Trees.GetMainTree().RootNode().NamedDescendantForPointRange(point, point)
	if node.Kind() == "atx_heading" {
		return
	}
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
