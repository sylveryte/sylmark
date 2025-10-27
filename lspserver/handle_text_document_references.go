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
	uri := params.TextDocument.URI
	id := h.Store.GetIdFromURI(uri)

	doc, node, ok := h.DocAndNodeFromURIAndPosition(id, params.Position, h.parse)
	content := string(doc.Content)
	if !ok {
		return
	}

	var locs []lsp.Location
	var idLocs []data.IdLocation

	switch node.Kind() {
	case "tag":
		{
			tag := data.GetTag(node, content)
			locs := h.Store.GetTagReferences(tag)
			return locs, nil
		}
	case "wiki_link", "link_destination", "atx_heading", "heading_content", "shortcut_link", "link_text", "inline_link":
		{
			parentedNode := lsp.GetParentalKind(node)

			switch parentedNode.Kind() {
			case "shortcut_link":
				doc, ok := h.Store.GetDoc(id)
				if ok {
					linkTextNode := parentedNode.NamedChild(0)
					linkText := lsp.GetNodeContent(*linkTextNode, string(doc.Content))
					footNote, ok := doc.FootNotes.GetFootNote(linkText)
					if ok {
						for _, r := range footNote.Refs {
							locs = append(locs, lsp.Location{
								URI:   params.TextDocument.URI,
								Range: r,
							})
						}
					}
				}
			case "inline_link":
				uri, ok := h.Store.Config.GetUriFromInlineNode(parentedNode, string(doc.Content), params.TextDocument.URI)
				if !ok {
					slog.Error("Failed to make uri ")
					return nil, nil
				}
				targetUri, subTarget, _ := h.Store.Config.GetMdRealUrlAndSubTarget(string(uri))
				if data.IsMdFile(string(uri)) {
					targetId := h.Store.GetIdFromURI(targetUri)
					lrefs, found := h.Store.LinkStore.GetRefs(targetId, subTarget)
					if found {
						for _, idLoc := range lrefs {
							idLocs = append(idLocs, idLoc)
						}
					}
				}

			case "atx_heading":
				subTarget, ok := data.GetSubTarget(parentedNode, string(content))
				if ok {
					refs, found := h.Store.LinkStore.GetRefs(id, subTarget)
					subrefs, subfound := doc.Headings.GetRefs(string(subTarget))
					if subfound {
						for _, sr := range subrefs {
							locs = append(locs, lsp.Location{
								Range: sr,
								URI:   uri,
							})
						}
					}
					if found {
						h.Store.FillInLocations(&locs, &refs)
					}
				}
			case "wiki_link":
				target, subTarget, isSubTarget, found := data.GetWikilinkTargets(parentedNode, content)
				if found {
					isSubheading := len(target) == 0 && isSubTarget
					if isSubheading {
						doc, ok := h.Store.GetDoc(id)
						if ok {
							ranges, ok := doc.Headings.GetRefs(string(subTarget))
							if ok {
								for _, r := range ranges {
									locs = append(locs, lsp.Location{
										URI:   params.TextDocument.URI,
										Range: r,
									})
								}
							}
						}
					} else {
						lidLocs, refFound := h.Store.GetRefsFromTarget(target, subTarget)
						if refFound {
							idLocs = append(idLocs, lidLocs...)
						}
					}
				}
			}
			locs = *h.Store.FillInLocations(&locs, &idLocs)
			if len(locs) > 0 {
				return locs, nil
			}
		}
	default:
		{
			// get files references
			lLocs, _ := h.Store.LinkStore.GetRefs(id, "")
			locs = *h.Store.FillInLocations(&locs, &lLocs)
			if len(locs) > 0 {
				return locs, nil
			}
		}
	}

	return nil, nil

}
