package lsp

import (
	"log/slog"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type GLink struct {
	Def  Location
	Refs []Location
}
type Store struct {
	Tags   map[Tag][]Location
	GLinks map[GTarget][]GLink
}

func newStore() Store {
	return Store{
		Tags:   map[Tag][]Location{},
		GLinks: map[GTarget][]GLink{},
	}
}

func (store *Store) loadData(uri DocumentURI, content string, rootNode *tree_sitter.Node) {
	TraverseNodeWith(rootNode, func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wikilink":
			{
				getWikilinkLink(n, content)
			}
		case "heading":
			{
				heading, ok := getHeadingTitle(n, content)
				if !ok {
					slog.Error("Could not extract heading")
					return
				}
				glink, ok := uri.getGTarget(heading)
				if !ok {
					slog.Error("Could not form glink")
					return
				}
				slog.Info("GTarget is " + string(glink))
			}
		case "tag":
			{
				store.AddTag(n, uri, &content)
			}
		case "link":
			{
				getLinkUrl(n, content)
			}
		}
	})
}

// removes data of tempStore from store
// removes on basis f URI only
func (store *Store) subtractStore(tempStore *Store) {
	s := *store
	ts := *tempStore
	if len(ts.Tags) > 0 {
		for k, v := range ts.Tags {
			sv, found := s.Tags[k]
			if found {
				newSv := []Location{}
				for _, loc := range sv {
					if loc.URI != v[0].URI {
						newSv = append(newSv, loc)
					} else {
					}
				}
				if len(newSv) == 0 {
					delete(s.Tags, k)
				} else {
					s.Tags[k] = newSv
				}
			}
		}
	}

	// syltodo TODO 🚧 wikilink headings n all

}

// appends data of tempStore into store
func (store *Store) mergeStore(tempStore *Store) {
	s := *store
	ts := *tempStore

	if len(ts.Tags) > 0 {
		for k, v := range ts.Tags {
			if sv, found := s.Tags[k]; found {
				s.Tags[k] = append(sv, v...)
			} else {
				s.Tags[k] = v
			}
		}
	}

	// syltodo TODO 🚧 wikilink headings n all

	return
}
