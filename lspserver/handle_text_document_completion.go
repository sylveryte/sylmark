package lspserver

import (
	"context"
	"encoding/json"
	"sylmark/data"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentCompletion(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	params.TextDocument.URI, _ = data.CleanUpURI(string(params.TextDocument.URI))

	locs, ok := h.Store.GetCompletions(params)
	return locs, ok

}
