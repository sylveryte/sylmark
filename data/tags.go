package data

import (
	"fmt"
	"strings"
	"sylmark-server/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Tag string

func (s *Store) GetTagRefs(tag Tag) int {

	clocs, found := s.Tags[tag]
	if found {
		return len(clocs)
	}

	return 0
}

func (s *Store) GetTagReferences(tag Tag) []lsp.Location {
	return s.Tags[tag]
}

func (s *Store) GetTagCompletions() []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}
	for t, v := range s.Tags {
		completions = append(completions, lsp.CompletionItem{
			Label:         string(t),
			Kind:          lsp.ClassCompletion,
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
	locations, found := s.Tags[tag]
	if found {
		s.Tags[tag] = append(locations, location)
	} else {
		s.Tags[tag] = []lsp.Location{location}
	}

	return true
}

// returns ok
func (s *Store) RemoveTag(node *tree_sitter.Node, uri lsp.DocumentURI, content *string) bool {
	if s == nil {
		return false
	}

	tag := GetTag(node, *content)
	loc := uri.LocationOf(node)
	tagLocs, found := s.Tags[tag]
	if found {
		var newLocations []lsp.Location

		for _, tagLoc := range tagLocs {
			if tagLoc.URI == loc.URI && tagLoc.Range.Start == loc.Range.Start {
				continue
			}
			newLocations = append(newLocations, tagLoc)
		}

		if len(newLocations) == 0 {
			delete(s.Tags, tag)
		} else {
			s.Tags[tag] = newLocations
		}
	}

	return true
}

func GetTag(node *tree_sitter.Node, content string) Tag {

	t := lsp.GetNodeContent(*node, content)
	t = strings.TrimSpace(t)

	return Tag(t)
}

func (s *Store) GetTagHover(tag Tag) string {
	if s == nil {
		return ""
	}
	totalRefs := s.GetTagRefs(tag)
	return fmt.Sprintf("%d references of %s", totalRefs, tag)
}
