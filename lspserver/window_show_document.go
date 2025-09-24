package lspserver

import (
	"context"
	"fmt"
	"sylmark/lsp"
)

func (h *LangHandler) ShowDocument(uri lsp.DocumentURI, external bool, rng lsp.Range) error {

	err := h.Connection.Notify(context.Background(), "window/showDocument",
		lsp.ShowDocumentParams{
			URI:       uri,
			External:  external,
			Selection: rng,
			TakeFocus: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to call window/showDocument: %w", err)
	}
	return nil
}
