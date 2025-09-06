package data

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (store *Store) GetDiagnostics(uri lsp.DocumentURI, parse lsp.ParseFunction) (items []lsp.Diagnostic) {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return
	}
	s := *store

	doc, ok := s.GetDocMustTree(uri, parse)
	if !ok {
		return
	}

	items = []lsp.Diagnostic{}

	content := string(doc.Content)
	lsp.TraverseNodeWith(doc.Tree.RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wikilink":
			{
				target, ok := GetWikilinkTarget(n, content, uri)
				if ok {
					_, found := s.GLinkStore.GetDefs(target)
					refs, rfound := s.GLinkStore.GetRefs(target)
					msg := "Unresolved"
					if rfound {
						if len(refs)>1{

						msg = fmt.Sprintf("%s referrenced %d times", msg, len(refs))
					}else{

						msg = fmt.Sprintf("%s referrence ", msg)
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
		case "heading":
			{
				target, ok := GetWikilinkTarget(n, content, uri)
				if ok {
					refs, found := s.GLinkStore.GetRefs(target)
					if found {
						msg := fmt.Sprintf("Referrenced %d times", len(refs))
						rng := lsp.GetRange(n)
						items = append(items, lsp.Diagnostic{
							Range:    &rng,
							Severity: lsp.DiagnosticSeverityInformation,
							Message:  msg,
						})
					}
				}
			}
		}
	})

	fileTarget, ok := GetFileGTarget(uri)
	if ok {
		refs, rfound := s.GLinkStore.GetRefs(GTarget(fileTarget))
		if rfound {
			msg := fmt.Sprintf("File is referrenced %d times", len(refs))
			items = append(items, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverityInformation,
				Message:  msg,
				Range:    &lsp.Range{},
			})

		}
	}

	return items
}
