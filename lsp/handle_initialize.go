package lsp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleInitialize(_ context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	// The rootUri of the workspace. Is null if no folder is open or no rootmakers added
	if params.RootURI != "" {
		rootPath, err := dirPathFromURI(params.RootURI)
		if err != nil {
			return nil, err
		}
		slog.Warn("RootPath is " + rootPath)
		h.addRootPathAndLoad(rootPath)
	}

	return InitializeResult{
		Capabilities: ServerCapabilities{
			HoverProvider:    true,
			TextDocumentSync: TDSKFull,
			CompletionProvider: &CompletionProvider{
				ResolveProvider:   true,
				TriggerCharacters: []string{"#", "[["},
			},
			DefinitionProvider: true,
			ReferencesProvider: true,
			SemanticTokensProvider: SemanticTokensOptions{
				Legend: SemanticTokensLegend{
					TokenTypes:     []SemanticTokenType{ClassSematicTokenType},
					TokenModifiers: []SemanticTokenModifier{},
				},
				Full:  true,
				Range: false,
			},
		},
	}, nil

}
