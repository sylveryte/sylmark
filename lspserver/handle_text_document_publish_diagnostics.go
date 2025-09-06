package lspserver

import (
	"context"
	"encoding/json"
	"sylmark-server/data"
	"sylmark-server/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleDiagnostics(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.DocumentDiagnosticParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	params.TextDocument.URI, _ = data.CleanUpURI(string(params.TextDocument.URI))

	items := h.Store.GetDiagnostics(params.TextDocument.URI, h.parse)

	result = lsp.DiagnosticResult{
		Kind:  lsp.DiagnosticReportFull,
		Items: items,
	}

	return result, nil
}
