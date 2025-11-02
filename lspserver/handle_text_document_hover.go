package lspserver

import (
	"context"
	"encoding/json"
	"fmt"
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
	id := h.Store.GetIdFromURI(params.TextDocument.URI)

	doc, node, ok := h.DocAndNodeFromURIAndPosition(id, params.Position, h.parse)
	if !ok {
		return nil, nil
	}

	// h.Store.TargetStore.Print()
	h.Store.IdStore.Print()
	h.Store.LinkStore.Print()

	r := lsp.GetRange(node)

	var content string
	switch node.Kind() {
	case "tag":
		{
			tag := data.GetTag(node, string(doc.Content))
			content = h.Store.GetTagHover(tag)
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
					refsText := fmt.Sprintf("%d references found\n", len(footNote.Refs))
					if ok && footNote.Def != nil {
						content = footNote.Excert
					} else {
						content = "_No definition found._"
					}
					content = refsText + content
				}
			case "inline_link":
				rng := lsp.Range{}
				_, targetId, subTarget, found := h.Store.GetInlineFullTargetAndSubTarget(parentedNode, string(doc.Content), id)
				if found {
					rng, _ = h.Store.LinkStore.GetDef(targetId, subTarget)
					content += h.Store.LinkStore.GetSubTargetHover(targetId, subTarget) + "\n---"
					content += h.Store.GetExcerpt(targetId, rng)
				}

			case "atx_heading":
				subTarget, ok := data.GetSubTarget(parentedNode, string(doc.Content))
				if ok {
					refs, found := h.Store.LinkStore.GetRefs(id, subTarget)
					subrefs, subfound := doc.Headings.GetRefs(string(subTarget))
					if found {
						content = fmt.Sprintf("%d references found\n", len(refs))
					}
					if subfound && len(subrefs) > 0 {
						content += fmt.Sprintf("%d references found within file\n", len(subrefs))
					}
				}
			case "wiki_link":
				target, subTarget, _, ok := data.GetWikilinkTargets(parentedNode, string(doc.Content))
				// utils.Sprintf("Idhar tak %s %s", target, subTarget)
				if ok {
					liddefs, defFound := h.Store.GetDefsFromTarget(target, subTarget)
					// utils.Sprintf("liddefs=%d", len(liddefs))
					if defFound {
						if len(liddefs) > 1 {
							content = fmt.Sprintf("%d definitions found\n", len(liddefs))
						}
						for _, ldef := range liddefs {
							content += h.Store.LinkStore.GetSubTargetHover(ldef.Id, subTarget) + "\n---"
							content += h.Store.GetExcerpt(ldef.Id, ldef.Range)
						}
					}
				}
			}
		}
	default:
		{
			target, _ := data.GetTarget(params.TextDocument.URI)
			content += fmt.Sprintf("File Details: `%s`\n---\n", target)
			// get files references
			lLocs, _ := h.Store.LinkStore.GetRefs(id, "")
			defs, defFound := h.Store.GetDefsFromTarget(target, "")
			if len(lLocs) > 0 {
				content += fmt.Sprintf("%d references found for the file in followings\n", len(lLocs))
				dMap := map[data.Id]bool{}
				for _, idLoc := range lLocs {
					if _, f := dMap[idLoc.Id]; f {
						continue
					}
					dMap[idLoc.Id] = true
					uri, ok := h.Store.GetUri(idLoc.Id)
					if ok {
						target, ok := h.Store.GetVaultTarget(uri)
						if ok {
							content += fmt.Sprintf("\n- %s", target)
						}
					}
				}
			} else {
				content += "No references found."
			}
			if defFound && len(defs) > 1 {
				content += fmt.Sprintf("\n\n%d files found with same file name.", len(lLocs))
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
