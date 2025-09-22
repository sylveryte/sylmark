package lspserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"sylmark/data"
	"sylmark/lsp"

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
				isSubheading := len(target) > 0 && target[0] == '#'
				if isSubheading {
					doc, ok := h.Store.GetDoc(params.TextDocument.URI)
					if ok {

						rng, ok := doc.Headings.GetDef(string(target))
						if ok {
							return lsp.Location{
								URI:   params.TextDocument.URI,
								Range: rng,
							}, nil
						}
					}

				} else {
					return h.Store.GetGTargetDefinition(target), nil
				}
			} else {
				slog.Warn("Wikilink definition not found" + string(target))
			}
		}
	}

	return nil, nil

}
