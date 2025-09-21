package lspserver

import (
	"context"
	"encoding/json"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleWorkspaceDidRenameFiles(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.RenameFilesParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	// syltodo
	// for _, v := range params.Files {
	// 	// h.onFileRenamed(v)
	// }

	return nil, nil
}
