package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleHover(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	slog.Info("Handling hover uri is = " + string(params.TextDocument.URI))

	doc, node, ok := h.DocAndNodeFromURIAndPosition(params.TextDocument.URI, params.Position)
	if !ok {
		return nil, nil
	}

	switch node.Kind() {
	case "tag":
		{
			r := getRange(node)
			tag := getTag(node, string(doc))
			totalRefs := h.getTagRefs(tag)
			return Hover{
				Contents: fmt.Sprintf("%d refs of %s", totalRefs, tag),
				Range:    &r,
			}, nil
		}
	}

	slog.Info("Node hovered is of kind = " + node.Kind())

	return nil, nil
}
