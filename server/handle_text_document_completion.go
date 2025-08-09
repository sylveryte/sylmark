package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentCompletion(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	slog.Info("params raw = " + string(*req.Params))
	// if params.CompletionContext.TriggerCharacter != nil {
	// 	char := *params.CompletionContext.TriggerCharacter
	// 	if char == "#" {
	// 		slog.Info("Params TriggerCharacter------------>" + *params.CompletionContext.TriggerCharacter)
	//
	// 	}
	// }

	// completions := []CompletionItem{}

	// Tags
	tagCompletions := h.store.GetTagCompletions()

	// wiklink
	// TODO 🚧

	slog.Info(fmt.Sprintf("Total Completions ------------>%d", len(tagCompletions)))

	return tagCompletions, nil
}
