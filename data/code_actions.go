package data

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sylmark/lsp"
)

func (store *Store) GetCodeActions(uri lsp.DocumentURI, diagnostics []lsp.Diagnostic, rng lsp.Range, parse lsp.ParseFunction) (actions []lsp.CodeAction) {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return
	}
	s := *store

	doc, ok := s.GetDocMustTree(uri, parse)
	if !ok {
		return
	}

	node := doc.Tree.RootNode().NamedDescendantForPointRange(lsp.PointFromPosition(rng.Start), lsp.PointFromPosition(rng.Start))

	switch node.Kind() {
	case "wikilink", "wikitarget":
		{
			target, ok := GetWikilinkTarget(node, string(doc.Content), uri)
			if !ok {
				slog.Warn("Wikilink definition not found" + string(target))
				return
			}
			defs := s.GetGTargetDefinition(target)
			if len(defs) == 0 {
				dir, err := DirPathFromURI(uri)
				if err != nil {
					slog.Warn("Error getting DirPathFromURI " + err.Error())
					return
				}
				fileTarget, heading, hasHeading := target.SplitHeading()

				fileName, _, _ := fileTarget.GetFileName()
				if hasHeading {
					fileDefs := s.GetGTargetDefinition(fileTarget)
					var fileUri string
					if len(fileDefs) > 0 {
						fileUri, _ = PathFromURI(fileDefs[0].URI)
					} else {
						fileUri = filepath.Join(dir, fileName)
					}
					actions = append(actions, lsp.CodeAction{
						Title:       fmt.Sprintf("Append heading `%s` in `%s`", heading, fileName),
						Diagnostics: diagnostics,
						Command: lsp.Command{
							Title:     "Append heading",
							Command:   "append",
							Arguments: []any{fileUri, "\n" + heading + "\n"},
						},
					})
				} else {
					fileUri := filepath.Join(dir, fileName)
					actions = append(actions, lsp.CodeAction{
						Title:       fmt.Sprintf("Create file: `%s.md`", string(target)),
						Diagnostics: diagnostics,
						Command: lsp.Command{
							Title:     "Create file",
							Command:   "create",
							Arguments: []any{fileUri},
						},
					})
				}
			}
		}
	}

	return actions
}
