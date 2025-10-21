package lspserver

import (
	"context"
	"encoding/json"
	"sylmark/data"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentDidChange(ctx context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	// h.Debouncers.DocumentDidChange.Debounce(func() {
	var params lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	changes := params.ContentChanges

	// rawUri := params.TextDocument.URI
	params.TextDocument.URI, _ = data.CleanUpURI(string(params.TextDocument.URI))

	// t := time.Now()
	for _, c := range changes {
		h.onDocChanged(params.TextDocument.URI, c)
	}
	// utils.Sprintf("=====>text change [[%dms]]<=====", time.Since(t).Milliseconds())

	// go h.ShowMessage(lsp.MessageTypeLog, "Changed")
	// utils.Sprintf(rawUri)
	// h.PublishDiagnostics(ctx, rawUri)

	return nil, nil
}
