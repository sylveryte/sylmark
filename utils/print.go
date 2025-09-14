package utils

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"
)

func PrintLocs(locs []lsp.Location) {
	slog.Info(fmt.Sprintf("Total locs => %d", len(locs)))
	for _, l := range locs {
		slog.Info(fmt.Sprintf("\n%s", string(l.URI.GetFileName())))
	}
}
