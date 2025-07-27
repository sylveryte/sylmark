package lsp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentDidOpen(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	content := params.TextDocument.Text

	slog.Info("Opened got uri " + string(params.TextDocument.URI))

	h.openedDocs.AddDoc(params.TextDocument.URI, Document(content))

	return nil, nil

}
