package server

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleShutdown(_ context.Context, conn *jsonrpc2.Conn, _ *jsonrpc2.Request) (result any, err error) {

	return nil, conn.Close()
}
