package data

import (
	"fmt"
	"strings"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (s *Store) GetTagRefs(tag Tag) int {

	clocs, found := s.tags[tag]
	if found {
		return len(clocs)
	}

	return 0
}

func (s *Store) GetTagReferences(tag Tag) []lsp.Location {
	return s.tags[tag]
}

func (s *Store) GetTagCompletions() []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}
	for t, v := range s.tags {
		completions = append(completions, lsp.CompletionItem{
			Label:         string(t),
			Kind:          lsp.ReferenceCompletion,
			Detail:        string(t),
			Documentation: fmt.Sprintf("#%d refs", len(v)),
		})
	}

	return completions
}

// returns ok
func (s *Store) AddTag(node *tree_sitter.Node, uri lsp.DocumentURI, content *string) bool {
	if s == nil {
		return false
	}

	tag := GetTag(node, *content)
	location := uri.LocationOf(node)
	locations, found := s.tags[tag]
	if found {
		s.tags[tag] = append(locations, location)
	} else {
		s.tags[tag] = []lsp.Location{location}
	}

	return true
}

func GetTag(node *tree_sitter.Node, content string) Tag {

	t := lsp.GetNodeContent(*node, content)
	t = strings.TrimSpace(t)

	return Tag(t)
}
