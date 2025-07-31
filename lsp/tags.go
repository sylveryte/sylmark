package lsp

import (
	"fmt"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (s *Store) GetTagReferences(tag Tag) []Location {
	return s.Tags[tag]
}

func (s *Store) getTagCompletions() []CompletionItem {
	completions := []CompletionItem{}
	for t, v := range s.Tags {
		completions = append(completions, CompletionItem{
			Label:         string(t),
			Kind:          ReferenceCompletion,
			Detail:        string(t),
			Documentation: fmt.Sprintf("#%d refs", len(v)),
		})
	}

	return completions
}

// returns ok
func (s *Store) AddTag(node *tree_sitter.Node, uri DocumentURI, content *string) bool {
	if s == nil {
		return false
	}

	tag := getTag(node, *content)
	location := locationFromURINode(uri, node)
	locations, found := s.Tags[tag]
	if found {
		s.Tags[tag] = append(locations, location)
	} else {
		s.Tags[tag] = []Location{location}
	}

	return true
}

func getTag(node *tree_sitter.Node, content string) Tag {

	t := getNodeContent(*node, Document(content))
	t = strings.TrimSpace(t)

	return Tag(t)
}
