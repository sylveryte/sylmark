package lspserver

import (
	"context"
	"encoding/json"
	"sylmark-server/data"
	"sylmark-server/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentDidClose(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	params.TextDocument.URI, _ = data.CleanUpURI(string(params.TextDocument.URI))
	h.onDocClosed(params.TextDocument.URI)

	return nil, nil

}
