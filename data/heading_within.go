package data

import (
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Subheading struct {
	Def  lsp.Range
	Refs []lsp.Range
}

func newSubheading() Subheading {
	return Subheading{
		Refs: []lsp.Range{},
	}
}

type HeadingsStore map[string]Subheading

func NewHeadingStore() HeadingsStore {
	return map[string]Subheading{}
}

func (store *HeadingsStore) GetDef(target string) (rng lsp.Range, found bool) {
	s := *store
	subHeading, found := s[target]
	if found {
		return subHeading.Def, found
	}
	return
}

// returns ok
func (store *HeadingsStore) SetDef(target string, rng lsp.Range) bool {
	if store == nil {
		return false
	}
	s := *store

	subh, found := s[target]
	if !found {
		subh = newSubheading()
	}
	subh.Def = rng
	s[target] = subh
	return true
}
func (store *HeadingsStore) AddRef(target string, rng lsp.Range) bool {
	if store == nil {
		return false
	}
	s := *store

	subh, found := s[target]
	if !found {
		subh = newSubheading()
	}
	subh.Refs = append(subh.Refs, rng)
	s[target] = subh
	return true
}
func (store *HeadingsStore) GetRefs(target string) (rng []lsp.Range, found bool) {
	s := *store
	subHeading, found := s[target]
	if found {
		return subHeading.Refs, found
	}
	return
}

func (s *Store) GetLoadedDataStore(id Id, parse lsp.ParseFunction) *HeadingsStore {

	store := NewHeadingStore()
	docData, ok := s.GetDocMustTree(id, parse)
	if ok {
		lsp.TraverseNodeWith(docData.Trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
			switch n.Kind() {
			case "atx_heading":
				{
					link, _ := GetSubTarget(n, string(docData.Content))
					store.SetDef(string(link), lsp.GetRange(n))
				}
			}
		})
		lsp.TraverseNodeWith(docData.Trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
			switch n.Kind() {
			case "wiki_link":
				{
					target, subTarget, isSubTarget, ok := GetWikilinkTargets(n, string(docData.Content))
					if ok {
						isSubheading := len(target) == 0 && isSubTarget
						if isSubheading {
							store.AddRef(string(subTarget), lsp.GetRange(n))
						}
					}
				}
			}
		})
	}
	return &store
}

func GetHeadings(docData *DocumentData) []string {
	if docData.Headings == nil {
		slog.Error("GetHeadings docData should not be nil")
		return []string{}
	}
	headings := []string{}
	if docData.Headings != nil {
		for k := range *docData.Headings {
			headings = append(headings, k)
		}
	}
	return headings
}
