package data

import (
	"fmt"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (store *Store) GetDiagnostics(uri lsp.DocumentURI, parse lsp.ParseFunction) (items []lsp.Diagnostic) {
	s := *store

	id := s.GetIdFromURI(uri)
	doc, ok := s.GetDocMustTree(id, parse)
	if !ok {
		return
	}

	items = []lsp.Diagnostic{}

	content := string(doc.Content)
	lsp.TraverseNodeWith(doc.Trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "atx_heading":
			{
				subTarget, ok := GetSubTarget(n, content)
				if ok {
					refs, found := s.LinkStore.GetRefs(id, subTarget)
					subrefs, subfound := doc.Headings.GetRefs(string(subTarget))
					if found || subfound {
						if subfound && len(subrefs) > 0 {
							rng := lsp.GetRange(n)
							msg := fmt.Sprintf("%d+%d=%d ", len(refs), len(subrefs), len(refs)+len(subrefs))
							if len(refs) == 0 {
								msg = fmt.Sprintf("+%d", len(subrefs))
							}
							items = append(items, lsp.Diagnostic{
								Range:    &rng,
								Severity: lsp.DiagnosticSeverityInformation,
								Message:  msg,
							})
						} else if len(refs) > 0 {

							rng := lsp.GetRange(n)
							items = append(items, lsp.Diagnostic{
								Range:    &rng,
								Severity: lsp.DiagnosticSeverityInformation,
								Message:  fmt.Sprintf("%d", len(refs)),
							})
						}
					}
				}
			}
		}
	})
	lsp.TraverseNodeWith(doc.Trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wiki_link":
			{
				target, subTarget, _, ok := GetWikilinkTargets(n, content)
				if ok {
					isSubheading := len(target) == 0
					if isSubheading {
						starget := string(subTarget)
						_, found := doc.Headings.GetDef(starget)
						if !found {
							rng := lsp.GetRange(n)
							items = append(items, lsp.Diagnostic{
								Range:    &rng,
								Severity: lsp.DiagnosticSeverityInformation,
								Tags:     []lsp.DiagnosticTag{lsp.DiagnosticTagUnnecessary},
								Message:  "Heading Unresolved",
							})
						}
					} else {
						_, found := s.GetDefsFromTarget(target, subTarget)
						refs, rfound := s.GetRefsFromTarget(target, subTarget)
						// utils.Sprintf("%d defs, %v found | target=%s subTarget%s %drefs", len(defs), found, target, subTarget,len(refs))
						msg := "Unresolved"
						if rfound {
							if len(refs) > 1 {

								msg = fmt.Sprintf("%s | %d ", msg, len(refs))
							} else {
								msg = fmt.Sprintf("%s ", msg)
							}
						}
						if !found {
							rng := lsp.GetRange(n)
							items = append(items, lsp.Diagnostic{
								Range:    &rng,
								Severity: lsp.DiagnosticSeverityInformation,
								Tags:     []lsp.DiagnosticTag{lsp.DiagnosticTagUnnecessary},
								Message:  msg,
							})
						}
					}
				}
			}
		}
	})

	refs, rfound := s.LinkStore.GetRefs(id, "")
	if rfound && len(refs) > 0 {
		msg := fmt.Sprintf("F%d", len(refs))
		items = append(items, lsp.Diagnostic{
			Severity: lsp.DiagnosticSeverityInformation,
			Message:  msg,
			Range:    &lsp.Range{},
		})
	}
	// file warning for duplicate names
	// target, ok := GetTarget(uri)
	// if ok {
	// 	defs, dfound := s.GetDefsFromTarget(target, "")
	// 	if dfound && len(defs) > 1 {
	// 		msg := fmt.Sprintf("File name has been used %d times. Consider different name for better wikilinks.", len(defs))
	// 		items = append(items, lsp.Diagnostic{
	// 			Severity: lsp.DiagnosticSeverityWarning,
	// 			Message:  msg,
	// 			Range:    &lsp.Range{},
	// 		})
	//
	// 	}
	// }

	return items
}
