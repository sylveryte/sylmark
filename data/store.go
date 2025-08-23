package data

import (
	"net/url"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Store struct {
	tags       map[Tag][]lsp.Location
	gLinkStore GLinkStore
}

func NewStore() Store {
	return Store{
		tags:       map[Tag][]lsp.Location{},
		gLinkStore: NewGlinkStore(),
	}
}

func (store *Store) LoadData(uri lsp.DocumentURI, content string, rootNode *tree_sitter.Node) {
	unscapedUri, err := url.QueryUnescape(string(uri))
	if err == nil {
		uri = lsp.DocumentURI(unscapedUri)
	}

	store.AddFileGTarget(uri)
	lsp.TraverseNodeWith(rootNode, func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wikilink":
			{
				GetWikilinkTarget(n, content, uri)
				// if ok {
				// 	slog.Info("got wikilink " + link)
				// }
			}
		case "heading":
			{
				store.AddGTarget(n, uri, &content)
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
	if len(ts.tags) > 0 {
		for k, v := range ts.tags {
			sv, found := s.tags[k]
			if found {
				newSv := []lsp.Location{}
				for _, loc := range sv {
					if loc.URI != v[0].URI {
						newSv = append(newSv, loc)
					} else {
					}
				}
				if len(newSv) == 0 {
					delete(s.tags, k)
				} else {
					s.tags[k] = newSv
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

	if len(ts.tags) > 0 {
		for k, v := range ts.tags {
			if sv, found := s.tags[k]; found {
				s.tags[k] = append(sv, v...)
			} else {
				s.tags[k] = v
			}
		}
	}

	// syltodo TODO 🚧 wikilink headings n all
}
