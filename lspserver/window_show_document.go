package lspserver

import (
	"context"
	"fmt"
	"sylmark/lsp"
)

// returns isSucess
func (h *LangHandler) ShowDocument(uri lsp.DocumentURI, external bool, rng lsp.Range) error {

	result := lsp.ShowDocumentResult{}
	// ctx  := context.WithTimeout(context.Background(), time.Second*3)
	ctx := context.Background()
	err := h.Connection.Call(ctx, "window/showDocument",
		lsp.ShowDocumentParams{
			URI:       uri,
			External:  external,
			Selection: rng,
			TakeFocus: true,
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
