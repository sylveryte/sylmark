package lspserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"sylmark/data"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleHover(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	params.TextDocument.URI, _ = data.CleanUpURI(string(params.TextDocument.URI))

	doc, node, ok := h.DocAndNodeFromURIAndPosition(params.TextDocument.URI, params.Position, h.parse)
	if !ok {
		return nil, nil
	}

	r := lsp.GetRange(node)
	var content string
	switch node.Kind() {
	case "tag":
		{
			tag := data.GetTag(node, string(doc))
			content = h.Store.GetTagHover(tag)
		}
	case "heading", "title":
		{
			target, ok := data.GetWikilinkTarget(node, string(doc), params.TextDocument.URI)
			if ok {
				content = h.Store.GetGTargetHeadingHover(target)
			} else {
				slog.Warn("Wikilink definition not found" + string(target))
			}
		}
	case "wikilink", "wikitarget":
		{

			target, ok := data.GetWikilinkTarget(node, string(doc), params.TextDocument.URI)
			if ok {
				content = h.Store.GetGTargetWikilinkHover(target)
			} else {
				slog.Warn("Wikilink definition not found" + string(target))
			}
		}
	}

	if len(content) > 0 {
		return lsp.Hover{
			Contents: content,
			Range:    &r,
		}, nil
	}

	return nil, nil
}
