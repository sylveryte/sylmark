package lspserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"sylmark-server/data"
	"sylmark-server/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentDefinition(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.DefinitionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	params.TextDocument.URI, _ = data.CleanUpURI(string(params.TextDocument.URI))

	doc, node, ok := h.DocAndNodeFromURIAndPosition(params.TextDocument.URI, params.Position, h.parse)
	if !ok {
		return nil, nil
	}

	switch node.Kind() {
	case "tag":
		{
			tag := data.GetTag(node, string(doc))
			locs := h.Store.GetTagReferences(tag)
			return locs, nil
		}
	case "wikilink", "wikitarget":
		{
			target, ok := data.GetWikilinkTarget(node, string(doc), params.TextDocument.URI)
			if ok {
				return h.Store.GetGTargetDefinition(target), nil
			} else {
				slog.Warn("Wikilink definition not found" + string(target))
			}
		}
	}

	return nil, nil

}
