package lsp

import (
	"context"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Config struct {
	RootMarkers *[]string `yaml:"root-markers" json:"rootMarkers"`
}

type LangHandler struct {
	Parser        *tree_sitter.Parser
	rootPath      string
	inactiveStore Store
	openedDocs    DocumentStore
}

func NewHandler() (hanlder *LangHandler) {
	return &LangHandler{
		inactiveStore: newStore(),
		openedDocs:    newDocumentStore(),
	}
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
	case "textDocument/didChange":
		return h.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/hover":
		return h.handleHover(ctx, conn, req)
	case "textDocument/completion":
		// return h.handleHover(ctx, conn, req)
	}
	return nil, nil
}
