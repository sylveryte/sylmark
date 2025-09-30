package lspserver

import (
	"context"
	"encoding/json"
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
	case "wiki_link", "link_destination", "atx_heading", "heading_content", "link_text":
		{
			target, ok := data.GetWikilinkTarget(node, string(doc), params.TextDocument.URI)
			if !ok {
				slog.Warn("No valid gtarget")
			}
			var locs []lsp.Location
			withinTarget, found := target.GetWIthinTarget()
			if found {
				docData, ok := h.Store.GetDoc(params.TextDocument.URI)
				if ok {
					ranges, found := docData.Headings.GetRefs(string(withinTarget))
					if found {
						for _, r := range ranges {
							locs = append(locs, lsp.Location{
								URI:   params.TextDocument.URI,
								Range: r,
							})
						}
					}
				}
			}
			locs = append(locs, h.Store.GetGTargetReferences(target)...)
			if len(locs) > 0 {
				return locs, nil
			}
		}
	default:
		{
			// get files references
			target, ok := data.GetFileGTarget(params.TextDocument.URI)
			if !ok {
				slog.Warn("No valid gtarget")
			}
			locs := h.Store.GetGTargetReferences(data.GTarget(target))
			if len(locs) > 0 {
				return locs, nil
			}
		}
	}

	return nil, nil

}
