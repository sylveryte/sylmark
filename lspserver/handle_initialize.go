package lspserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"sylmark/data"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleInitialize(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	// The rootUri of the workspace. Is null if no folder is open or no rootmakers added
	if params.RootURI != "" {
		rootPath, err := data.DirPathFromURI(params.RootURI)
		if err != nil {
			return nil, err
		}
		slog.Warn("RootPath is " + rootPath)
		h.addRootPathAndLoad(rootPath)
	}

	fileOperationRegistrationOptions := lsp.FileOperationRegistrationOptions{
		Filters: []lsp.FileOperationFilter{
			{Scheme: "file", Pattern: lsp.FileOperationPattern{Glob: "**/*.md"}},
		},
	}
	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			HoverProvider:    true,
			TextDocumentSync: lsp.TDSKFull,
			CompletionProvider: &lsp.CompletionProvider{
				ResolveProvider:   true,
				TriggerCharacters: []string{"[[", "|", "#"},
			},
			Workspace: lsp.ServerCapabilitiesWorkspace{
				FileOperations: lsp.FileOperations{
					DidDelete: fileOperationRegistrationOptions,
					DidRename: fileOperationRegistrationOptions,
					DidCreate: fileOperationRegistrationOptions,
				},
			},
			WorkspaceSymbolProvider: lsp.WorkspaceSymbolOptions{
				ResolveProvider: true,
			},
			DefinitionProvider: true,
			ReferencesProvider: true,
			DiagnosticProvider: lsp.DiagnosticOptions{
				InterFileDependencies: true,
				WorkspaceDiagnostics:  false,
			},
			CodeActionProvider: true,
			ExecuteCommandProvider: lsp.ExecuteCommandOptions{
				Commands: []string{"show", "graph"},
			},
			SemanticTokensProvider: lsp.SemanticTokensOptions{
				Legend: lsp.SemanticTokensLegend{
					TokenTypes:     []lsp.SemanticTokenType{lsp.ClassSematicTokenType},
					TokenModifiers: []lsp.SemanticTokenModifier{},
				},
				Full:  true,
				Range: false,
			},
		},
	}, nil

}
