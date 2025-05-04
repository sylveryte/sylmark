package lsp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleInitialize(_ context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	slog.Info("InitializeParams rootUri is" + string(params.RootURI))

	// The rootUri of the workspace. Is null if no folder is open.
	if params.RootURI != "" {
		rootPath, err := fromURI(params.RootURI)
		if err != nil {
			return nil, err
		}
		slog.Warn("RootPath is " + rootPath)
		// h.rootPath = filepath.Clean(rootPath)
		// h.addFolder(rootPath)
	}

	return InitializeResult{
		Capabilities: ServerCapabilities{
			HoverProvider: true,
			TextDocumentSync: TDSKFull,
		},
	}, nil

}
