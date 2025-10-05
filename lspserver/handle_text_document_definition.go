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
	case "wiki_link", "link_destination", "link_text", "shortcut_link", "inline_link":
		{
			parentedNode := lsp.GetParentalKind(node)
			switch parentedNode.Kind() {
			case "shortcut_link":
				doc, ok := h.Store.GetDoc(params.TextDocument.URI)
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
				filePath, err := data.GetInlineLinkTarget(parentedNode, string(doc), params.TextDocument.URI)
				if err != nil {
					slog.Error("File doesnt exist")
					return nil, nil
				}
				fullFilePath, err := data.GetFullPathRelatedTo(params.TextDocument.URI, filePath)
				if err != nil {
					slog.Error("Fialed to get full path" + err.Error())
					return nil, nil
				}
				uri, err := data.UriFromPath(fullFilePath)
				if err != nil {
					slog.Error("Failed to make uri " + err.Error())
					return nil, nil
				}
				return lsp.Location{
					URI:   uri,
					Range: lsp.Range{},
				}, nil

			case "wiki_link":
				target, ok := data.GetWikilinkTarget(parentedNode, string(doc), params.TextDocument.URI)
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
	}

	return nil, nil

}
