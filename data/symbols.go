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
	query = strings.TrimSpace(query)

	for uri, id := range s.IdStore.uri {
		rootPath, err := s.GetPathRelRoot(uri)
		if err != nil {
			continue
		}
		match := fuzzy.MatchFold(query, rootPath)
		if match {
			var name string
			target, ok := GetTarget(uri)
			if ok {
				name = string(target)
			} else {
				name = rootPath
			}
			symbols = append(symbols, lsp.WorkspaceSymbol{
				Name: name,
				Kind: lsp.SymbolKindFile,
				Location: lsp.Location{
					URI:   uri,
					Range: lsp.Range{},
				},
			})
		}
		if !isFileOnly {
			// add subtargets as well
			for _, st := range s.LinkStore.GetSubTargetsAndRanges(id) {
				nr := rootPath + "#" + string(st.subTarget)
				match := fuzzy.MatchFold(query, nr)
				if match {
					var name string
					target, ok := GetTarget(uri)
					if ok {
						name = string(target) + "#" + string(st.subTarget)
					} else {
						name = nr
					}
					if st.rng == nil {
						continue
					}
					symbols = append(symbols, lsp.WorkspaceSymbol{
						Name: name,
						Kind: lsp.SymbolKindKey,
						Location: lsp.Location{
							URI:   uri,
							Range: *st.rng,
						},
					})
				}
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
