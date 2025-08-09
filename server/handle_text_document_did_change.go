package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentDidChange(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	changes := params.ContentChanges

	slog.Info("Changed got uri " + string(params.TextDocument.URI))
	for _, c := range changes {
		h.onDocChanged(params.TextDocument.URI, c)
	}

	return nil, nil
}
