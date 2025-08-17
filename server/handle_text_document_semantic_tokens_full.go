package server

import (
	"context"
	"encoding/json"
	"sylmark/data"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentSemanticTokensFull(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.SemanticTokensParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	tokens := data.GetSemanticTokens(h.openedDocs, params.TextDocument.URI)

	return tokens, nil
}
