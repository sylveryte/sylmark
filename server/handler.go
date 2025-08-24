package server

import (
	"context"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"sylmark/data"
	"sylmark/lsp"
	"sylmark/utils"
	"time"

	tree_sitter_sylmark "codeberg.org/sylveryte/tree-sitter-sylmark/bindings/go"
	"github.com/sourcegraph/jsonrpc2"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Config struct {
	RootMarkers   *[]string `yaml:"root-markers" json:"rootMarkers"`
}

func NewConfig() Config {
	rmakers := []string{".sylroot"}
	return Config{
		RootMarkers:   &rmakers,
	}
}

type ServerDebouncers = struct {
	DocumentDidChange   *utils.SylDebouncer
	SemantickTokensFull *utils.SylDebouncer
}

type LangHandler struct {
	Parser     *tree_sitter.Parser
	rootPath   string
	store      data.Store
	Debouncers *ServerDebouncers
	Config     Config
}

func NewHandler() (hanlder *LangHandler) {
	return &LangHandler{
		store:  data.NewStore(),
		Config: NewConfig(),
		Debouncers: &ServerDebouncers{
			DocumentDidChange:   utils.NewSylDebouncer(300 * time.Millisecond),
			SemantickTokensFull: utils.NewSylDebouncer(400 * time.Millisecond),
		},
	}
}

func (h *LangHandler) addRootPathAndLoad(dir string) {
	h.rootPath = dir
	h.loadAllClosedDocsData()
}

func (h *LangHandler) loadAllClosedDocsData() {
	if h.rootPath == "" {
		slog.Error("h.rootPath is empty")
		return
	}

	filepath.WalkDir(h.rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.HasSuffix(path, ".git") {
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
	tree := h.parse(content, nil)
	defer tree.Close()
	h.store.LoadData(uri, content, tree.RootNode())
}

func (h *LangHandler) onDocOpened(uri lsp.DocumentURI, content string) {
	tree := h.parse(content, nil)
	doc := data.Document(content)

	docData := data.NewDocumentData(doc, tree)
	h.store.AddUpdateDoc(uri, docData)
}

func (h *LangHandler) onDocClosed(uri lsp.DocumentURI) {
	// remove data into openedDocs
}

func (h *LangHandler) onDocChanged(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent) {
	h.store.SyncChangedDocument(uri, changes, h.parse)
}

func (h *LangHandler) SetupGrammars() {
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_sylmark.Language())
	parser.SetLanguage(language)

	h.Parser = parser
}

func (h *LangHandler) DocAndNodeFromURIAndPosition(uri lsp.DocumentURI, position lsp.Position) (doc data.Document, node *tree_sitter.Node, ok bool) {
	docData, ok := h.store.GetDoc(uri)
	if !ok {
		slog.Error("Document missing" + string(uri))
		return "", nil, false
	}
	point := lsp.PointFromPosition(position)

	doc = docData.Content
	node = docData.Tree.RootNode().NamedDescendantForPointRange(point, point)

	ok = true
	return
}

func (h *LangHandler) parse(content string, oldTree *tree_sitter.Tree) *tree_sitter.Tree {
	return h.Parser.Parse([]byte(content), oldTree)
}

func (h *LangHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	slog.Info("Received request with Handle method=> " + req.Method)
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, conn, req)
	case "initialized":
		return
	case "shutdown":
		return h.handleShutdown(ctx, conn, req)
	case "textDocument/didOpen":
		return h.handleTextDocumentDidOpen(ctx, conn, req)
	case "textDocument/didClose":
		return h.handleTextDocumentDidClose(ctx, conn, req)
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
	}
	return nil, nil
}
