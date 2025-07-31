package lsp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentSemanticTokensFull(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params SemanticTokensParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	slog.Info("params raw = " + string(*req.Params))

	tokens := h.getSemanticTokens(params.TextDocument.URI)

	return tokens, nil
}
