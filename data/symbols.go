package data

import (
	"strings"
	"sylmark/lsp"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

func (s *Store) GetAllSymbols(query string) (symbols []lsp.WorkspaceSymbol) {

	if len(query) == 0 {
		return
	}
	isFileOnly := query[0] == ' '

	for t, gl := range s.GLinkStore {
		target := string(t)

		if len(gl.Defs) > 0 {
		outer:
			for _, loc := range gl.Defs {
				match := fuzzy.MatchFold(query, loc.URI.GetFileName()+" "+target)
				if match == false {
					continue outer
				}
				isFile := !strings.ContainsRune(target, '#')
				if isFileOnly && !isFile {
					continue
				}
				kind := lsp.SymbolKindKey
				if isFile {
					kind = lsp.SymbolKindFile
				}

				symbols = append(symbols, lsp.WorkspaceSymbol{
					Name:     target,
					Kind:     kind,
					Location: loc,
				})
			}
		}
	}

	for t, refs := range s.Tags {
		target := string(t)
		match := fuzzy.MatchFold(query, target)
		if match == false {
			continue
		}

		for _, loc := range refs {
			symbols = append(symbols, lsp.WorkspaceSymbol{
				Name:     target,
				Kind:     lsp.SymbolKindEnum,
				Location: loc,
			})
		}
	}

	return
}
