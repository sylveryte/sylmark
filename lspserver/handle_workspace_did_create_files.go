package lspserver

import (
	"context"
	"encoding/json"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleWorkspaceDidCreateFiles(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.CreateFilesParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	for _, v := range params.Files {
		h.onDocCreated(v.Uri, "")
	}

	return nil, nil
}
