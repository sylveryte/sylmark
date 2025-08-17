package server

import (
	"context"
	"encoding/json"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentDidChange(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	// h.Debouncers.DocumentDidChange.Debounce(func() {
	var params lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	changes := params.ContentChanges

	for _, c := range changes {
		h.onDocChanged(params.TextDocument.URI, c)
	}
	// })

	return nil, nil
}
