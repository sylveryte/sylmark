package data

import (
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// returns ok
func (s *Store) AddGTarget(node *tree_sitter.Node, uri lsp.DocumentURI, content *string) bool {
	if s == nil {
		return false
	}

	heading, ok := getHeadingTitle(node, *content)
	if !ok {
		slog.Error("Could not extract heading")
		return false
	}
	gtarget, ok := GetGTarget(heading, uri)
	if !ok {
		slog.Error("Could not form gtarget")
		return false
	}

	location := uri.LocationOf(node)

	s.gLinkStore.AddDef(gtarget, location.URI, location.Range)

	return true
}

func (s *Store) GetWikiCompletions(isWikiEnd bool, uri *lsp.DocumentURI) []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}

	for _, t := range s.gLinkStore.GetTargets() {
		var link string
		if isWikiEnd {
			link = "[[" + string(t) + "]]"
		} else {
			link = "[[" + string(t)
		}
		completions = append(completions, lsp.CompletionItem{
			Label:  link,
			Kind:   lsp.FileCompletion,
			Detail: string(t),
		})
	}

	return completions
}
