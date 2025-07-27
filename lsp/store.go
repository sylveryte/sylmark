package lsp

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Store struct {
	Tags map[Tag][]Location
}

func newStore() Store {
	return Store{
		Tags: map[Tag][]Location{},
	}
}

// returns ok
func (s *Store) AddTag(node *tree_sitter.Node, uri DocumentURI, content *string) bool {
	if s == nil {
		return false
	}

	tag := getTag(node, *content)
	location := getLocation(uri, node)
	locations, found := s.Tags[tag]
	if found {
		s.Tags[tag] = append(locations, location)
	} else {
		s.Tags[tag] = []Location{location}
	}

	return true
}
