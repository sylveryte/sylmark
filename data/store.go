package data

import (
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type GLink struct {
	Def  lsp.Location
	Refs []lsp.Location
}
type Store struct {
	Tags   map[Tag][]lsp.Location
	GLinks map[GTarget][]GLink
}

func NewStore() Store {
	return Store{
		Tags:   map[Tag][]lsp.Location{},
		GLinks: map[GTarget][]GLink{},
	}
}

func (store *Store) LoadData(uri lsp.DocumentURI, content string, rootNode *tree_sitter.Node) {
	lsp.TraverseNodeWith(rootNode, func(n *tree_sitter.Node) {
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
				glink, ok := GetGTarget(heading, uri)
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
func (store *Store) SubtractStore(tempStore *Store) {
	s := *store
	ts := *tempStore
	if len(ts.Tags) > 0 {
		for k, v := range ts.Tags {
			sv, found := s.Tags[k]
			if found {
				newSv := []lsp.Location{}
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
func (store *Store) MergeStore(tempStore *Store) {
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
}
