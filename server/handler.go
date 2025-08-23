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

	"github.com/sourcegraph/jsonrpc2"
	tree_sitter_sylmark "codeberg.org/sylveryte/tree-sitter-sylmark/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Config struct {
	RootMarkers   *[]string `yaml:"root-markers" json:"rootMarkers"`
	ExcerptLength int16
}

func NewConfig() Config {
	rmakers := []string{".sylroot"}
	return Config{
		RootMarkers:   &rmakers,
		ExcerptLength: 10,
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
	openedDocs data.DocumentStore
	Debouncers *ServerDebouncers
	Config     Config
}

func NewHandler() (hanlder *LangHandler) {
	return &LangHandler{
		store:      data.NewStore(),
		openedDocs: data.NewDocumentStore(),
		Config:     NewConfig(),
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
	content := data.ContentFromDocPath(mdDocPath)
	uri, err := data.UriFromPath(mdDocPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	tree := h.parse(content)
	defer tree.Close()
	h.store.LoadData(uri, content, tree.RootNode())
}

func (h *LangHandler) onDocOpened(uri lsp.DocumentURI, content string) {
	tree := h.parse(content)
	h.openedDocs.AddDoc(uri, data.Document(content), tree)
	doc := data.Document(content)

	h.openedDocs.AddDoc(uri, doc, tree)
}
func (h *LangHandler) onDocClosed(uri lsp.DocumentURI) {
	// remove data into openedDocs
	_, found := h.openedDocs.RemoveDoc(uri)
	if !found {
		slog.Error("Document not in openedDocs")
		return
	}
}

func (h *LangHandler) onDocChanged(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent) {

	// update data into openedDocs
	updatedDocData, oldDocData, ok := h.openedDocs.UpdateDoc(uri, changes, h.parse)
	if !ok {
		slog.Info("Update doc failed.")
		return
	}

	// update openedDocsStore
	tempStoreOld := data.NewStore()
	tempStoreOld.LoadData(uri, string(oldDocData.Content), oldDocData.Tree.RootNode())

	tempStoreNew := data.NewStore()
	tempStoreNew.LoadData(uri, string(updatedDocData.Content), updatedDocData.Tree.RootNode())

	// syltodo TODO optimze this flow it
	h.store.SubtractStore(&tempStoreOld)
	h.store.MergeStore(&tempStoreNew)

}

func (h *LangHandler) SetupGrammars() {
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_sylmark.Language())
	parser.SetLanguage(language)

	h.Parser = parser

	slog.Info("Grammars are set")
}

func (h *LangHandler) DocAndNodeFromURIAndPosition(uri lsp.DocumentURI, position lsp.Position) (doc data.Document, node *tree_sitter.Node, ok bool) {
	docData, ok := h.openedDocs.DocDataFromURI(uri)
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

func (h *LangHandler) parse(content string) *tree_sitter.Tree {
	return h.Parser.Parse([]byte(content), nil)
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
