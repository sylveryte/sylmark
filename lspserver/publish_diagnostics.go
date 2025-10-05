package lspserver

import (
	"context"
	"log/slog"
	"sylmark/lsp"
)

func (h *LangHandler) PublishDiagnostics(ctx context.Context, uri lsp.DocumentURI) {
	slog.Info("Publishing uri=" + string(uri))
	h.Connection.Notify(
		context.Background(),
		"textDocument/publishDiagnostics",
		lsp.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: h.Store.GetDiagnostics(uri, h.parse),
		},
	)
}
