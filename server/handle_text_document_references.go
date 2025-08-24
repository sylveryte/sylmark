package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sylmark/data"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentReferences(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.ReferencesParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	doc, node, ok := h.DocAndNodeFromURIAndPosition(params.TextDocument.URI, params.Position)
	if !ok {
		return nil, nil
	}

	switch node.Kind() {
	case "tag":
		{
			tag := data.GetTag(node, string(doc))
			locs := h.store.GetTagReferences(tag)
			return locs, nil
		}
	case "wikilink", "wikitarget", "heading", "title":
		{
			target, ok := data.GetWikilinkTarget(node, string(doc), params.TextDocument.URI)
			if !ok {
				slog.Warn("No valid gtarget")
			}
			locs := h.store.GetGTargetReferences(target)
			slog.Info(fmt.Sprintf("Got %d", len(locs)))
			if len(locs) > 0 {
				return locs, nil
			}
		}
	}

	return nil, nil

}
