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

	id := h.Store.GetIdFromURI(params.TextDocument.URI)
	doc, node, ok := h.DocAndNodeFromURIAndPosition(id, params.Position, h.parse)
	if !ok {
		return nil, nil
	}

	switch node.Kind() {
	case "tag":
		{
			tag := data.GetTag(node, string(doc.Content))
			locs := h.Store.GetTagReferences(tag)
			return locs, nil
		}
	case "wiki_link", "link_destination", "link_text", "shortcut_link", "inline_link":
		{
			parentedNode := lsp.GetParentalKind(node)
			switch parentedNode.Kind() {
			case "shortcut_link":
				doc, ok := h.Store.GetDoc(id)
				if ok {
					linkTextNode := parentedNode.NamedChild(0)
					linkText := lsp.GetNodeContent(*linkTextNode, string(doc.Content))
					ref, ok := doc.FootNotes.GetFootNote(linkText)
					if ok && ref.Def != nil {
						return lsp.Location{
							URI:   params.TextDocument.URI,
							Range: *ref.Def,
						}, nil
					}
				}
			case "inline_link":
				uri, ok := h.Store.Config.GetUriFromInlineNode(parentedNode, string(doc.Content), params.TextDocument.URI)
				if !ok {
					slog.Error("Failed to make uri ")
					return nil, nil
				}
				rng := lsp.Range{}
				uri, subTarget, found := h.Store.Config.GetMdRealUrlAndSubTarget(string(uri))
				if data.IsMdFile(string(uri)) {
					if found {
						id := h.Store.GetIdFromURI(uri)
						rng, _ = h.Store.LinkStore.GetDef(id, subTarget)
					}
					return lsp.Location{
						URI:   uri,
						Range: rng,
					}, nil
				}

			case "wiki_link":
				target, subTarget, isSubTarget, _ := data.GetWikilinkTargets(parentedNode, string(doc.Content))
				if ok {
					isSubheading := len(target) == 0 && isSubTarget
					if isSubheading {
						doc, ok := h.Store.GetDoc(id)
						if ok {

							rng, ok := doc.Headings.GetDef(string(subTarget))
							if ok {
								return lsp.Location{
									URI:   params.TextDocument.URI,
									Range: rng,
								}, nil
							}
						}

					} else {
						// is id even proper?? for subtarget??
						// check others refs hover
						locs := []lsp.Location{}
						defs, found := h.Store.GetDefsFromTarget(target, subTarget)

						if !found {
							// file doesn't exists create uri and open it
							fileName := target.GetFileName()
							newURI, err := data.GetFileURIInSameURIPath(fileName, params.TextDocument.URI)
							if err != nil {
								return defs, err
							}
							loc := lsp.Location{
								URI: newURI,
							}
							return loc, nil
						}

						locs = *h.Store.FillInLocations(&locs, &defs)
						return locs, nil
					}
				} else {
					slog.Warn("No Target detected = " + string(target))
				}
			}
		}
	}

	return nil, nil

}
