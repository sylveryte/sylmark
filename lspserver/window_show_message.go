package lspserver

import (
	"context"
	"fmt"
	"sylmark/lsp"
)

// returns isSucess
func (h *LangHandler) ShowMessage(typ lsp.MessageType, msg string) error {

	result := lsp.ShowDocumentResult{}
	// ctx  := context.WithTimeout(context.Background(), time.Second*3)
	ctx := context.Background()
	err := h.Connection.Call(ctx, "window/showMessage",
		lsp.ShowMessageParams{
			Type:    typ,
			Message: msg,
		},
		&result,
	)
	if err != nil {
		return fmt.Errorf("failed to call window/showDocument: %w", err)
	}
	if !result.Success {
		return fmt.Errorf("client failed to open document")
	}
	return nil
}
