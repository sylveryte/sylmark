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

	slog.Info("Got uri " + string(params.TextDocument.URI))
	slog.Info("Got text " + string(params.TextDocument.Text))


	// doc := parseGoldmark(content)
	// slog.Info(fmt.Sprintf("Doc Type %d, Kind %s", doc.Type(), doc.Kind()))

	h.parseTreesitter(content)

	return nil, nil

}
