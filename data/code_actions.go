package data

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sylmark/lsp"
)

func (store *Store) GetCodeActions(id Id, diagnostics []lsp.Diagnostic, rng lsp.Range, parse lsp.ParseFunction) (actions []lsp.CodeAction) {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return
	}
	s := *store

	doc, ok := s.GetDocMustTree(id, parse)
	if !ok {
		return
	}

	node := doc.Trees.GetInlineTree().RootNode().NamedDescendantForPointRange(lsp.PointFromPosition(rng.Start), lsp.PointFromPosition(rng.Start))

	node = lsp.GetParentalKind(node)
	if node.Kind() == "wiki_link" {
		target, subTarget, isSubTarget, ok := GetWikilinkTargets(node, string(doc.Content))
		if !ok {
			slog.Warn("Wikilink definition not found" + string(target))
			return
		}

		defIdLocs, ok := s.GetDefsFromTarget(target, subTarget)

		if len(defIdLocs) != 0 {
			return
		}
		uri, _ := s.GetUri(id)
		dir, err := DirPathFromURI(uri)
		if err != nil {
			slog.Warn("Error getting DirPathFromURI " + err.Error())
			return
		}
		fileName := target.GetFileName()
		fileUri := filepath.Join(dir, fileName)
		if isSubTarget {
			ids := s.getIds(target)
			title := "Append heading"
			Title := fmt.Sprintf("Append heading `%s` in `%s`", subTarget, fileName)
			var fileUris []string
			if ids[0] == 0 {
				// append heading only need
				title = "Create file and append heading"
				Title = fmt.Sprintf("Create file and Append heading `%s` in `%s`", subTarget, fileName)

				fileUris = append(fileUris, fileUri)
			} else {
				for _, id := range ids {
					uri, ok = s.GetUri(id)
					if ok {
						fileUris = append(fileUris, string(uri))
					}
				}
			}
			actions = append(actions, lsp.CodeAction{
				Title:       Title,
				Diagnostics: diagnostics,
				Command: lsp.Command{
					Title:     title,
					Command:   "append",
					Arguments: []any{fileUri, "\n" + subTarget + "\n"},
				},
			})
		} else {
			actions = append(actions, lsp.CodeAction{
				Title:       fmt.Sprintf("Create file: `%s`", fileName),
				Diagnostics: diagnostics,
				Command: lsp.Command{
					Title:     "Create file",
					Command:   "create",
					Arguments: []any{fileUri},
				},
			})
		}
	}

	return actions
}
