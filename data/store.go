package data

import (
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Store struct {
	Tags          map[Tag][]lsp.Location
	LinkStore     LinkStore
	DocStore      DocumentStore
	TargetStore   TargetStore
	IdStore       IdStore
	DateStore     DateStore
	LastOpenedDoc lsp.DocumentURI
	ExcerptLength int16
	Config        Config
	OtherFiles    []string
}

func NewStore() Store {
	return Store{
		Tags:          map[Tag][]lsp.Location{},
		LinkStore:     NewlinkStore(),
		TargetStore:   NewTargetStore(),
		DocStore:      NewDocumentStore(),
		DateStore:     NewDateStore(),
		IdStore:       NewIdStore(),
		OtherFiles:    []string{},
		Config:        NewConfig(),
		ExcerptLength: 10,
	}
}

func (s *Store) SyncChangedDocument(id Id, changes lsp.TextDocumentContentChangeEvent, parse lsp.ParseFunction) {

	var updatedDocData, oldDocData DocumentData
	// update data into openedDocs
	doc := Document(changes.Text)
	staleDoc, ok := s.GetDocMustTree(id, parse)
	oldDocData = staleDoc
	if !ok {
		slog.Error("Failed to get old file")
		return
	}
	updatedDocData = *s.UpdateAndReloadDoc(id, string(doc), parse)
	s.UnloadData(id, string(oldDocData.Content), oldDocData.Trees)
	s.LoadData(id, string(updatedDocData.Content), updatedDocData.Trees)

	// sylopti this is for lsp.TDSKIncremental
	// TextDocumentSync is set to TDSKFull so this case won't be there but in future let's implment partial for better perf
	// utils.Sprintf("RangeLength changed %d", changes.RangeLength)
	// slog.Error("Need to handle partial change text")
	// // return
	// utils.Sprintf(fmt.Sprintf("range length %d", changes.RangeLength)
	// utils.Sprintf
	// 	"range line start %d end %d",
	// 	changes.Range.Start.Line,
	// 	changes.Range.End.Line,
	// ))
	// utils.Sprintf
	// 	"range char start %d end %d",
	// 	changes.Range.Start.Character,
	// 	changes.Range.End.Character,
	// ))
}

func (s *Store) UnloadData(id Id, content string, trees *lsp.Trees) {
	// utils.Sprintf("UnloadData id=%d", id)
	uri, _ := s.GetUri(id)
	lsp.TraverseNodeWith(trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "atx_heading":
			{
				subTarget, ok := GetSubTarget(n, content)
				if ok {
					s.LinkStore.RemoveDef(id, subTarget, lsp.GetRange(n))
				}
			}
		}
	})
	lsp.TraverseNodeWith(trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wiki_link":
			{
				target, subTarget, _, ok := GetWikilinkTargets(n, content)
				if ok {
					ids := s.getIds(target)
					for _, defId := range ids {
						loc := id.LocationOf(n)
						s.LinkStore.RemoveRef(defId, subTarget, loc)
					}
				}
			}
		case "tag":
			{
				s.RemoveTag(n, uri, &content)
			}
		case "inline_link":
			linkUri, ok := GetUriFromInlineNode(n, content, uri)
			subTarget := SubTarget("")
			linkUri, subTarget, _ = GetUrlAndSubTarget(string(linkUri))
			if ok && IsMdFile(string(linkUri)) {
				linkId := s.GetIdFromURI(linkUri)
				loc := id.LocationOf(n)
				s.LinkStore.RemoveRef(linkId, subTarget, loc)
			}
		}
	})
}

func (s *Store) LoadData(id Id, content string, trees *lsp.Trees) {
	// utils.Sprintf("LoadData id=%d", id)

	uri, _ := s.GetUri(id)
	s.LinkStore.AddFileGTarget(id)
	lsp.TraverseNodeWith(trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "atx_heading":
			{
				subTarget, ok := GetSubTarget(n, content)
				if ok {
					s.LinkStore.AddDef(id, subTarget, lsp.GetRange(n))
				}
			}
		}
	})

	lsp.TraverseNodeWith(trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wiki_link":
			{
				target, subTarget, _, ok := GetWikilinkTargets(n, content)

				if ok {
					isSubheading := len(target) == 0
					if !isSubheading {
						defIds := s.getIds(target)
						for _, defId := range defIds {
							loc := id.LocationOf(n)
							s.LinkStore.AddRef(defId, subTarget, loc)
						}
					}
				}
			}
		case "tag":
			{
				s.AddTag(n, uri, &content)
			}
		case "inline_link":
			linkUri, ok := GetUriFromInlineNode(n, content, uri)
			subTarget := SubTarget("")
			linkUri, subTarget, _ = GetUrlAndSubTarget(string(linkUri))
			if ok && IsMdFile(string(linkUri)) {
				linkId := s.GetIdFromURI(linkUri)
				loc := id.LocationOf(n)
				s.LinkStore.AddRef(linkId, subTarget, loc)
			}
		}
	})
}

func (s *Store) UpdateAndReloadDoc(id Id, content string, parse lsp.ParseFunction) *DocumentData {
	// t := time.Now()
	trees := parse(content, nil)
	doc := Document(content)
	// utils.Sprintf("%dms<==parsing time", time.Since(t).Milliseconds())

	// utils.Sprintf(*trees.GetMainTree().RootNode(), 0, content
	// utils.Sprintf(*trees.GetInlineTree().RootNode(), 0, content
	docData := NewDocumentData(doc, trees)

	// important to update doc first since GetLoadedDataStore fetches it
	s.AddUpdateDoc(id, docData)
	docData.Headings = s.GetLoadedDataStore(id, parse)
	docData.FootNotes = s.GetLoadedFootNotesStore(id, parse)
	// finally update with stores
	s.AddUpdateDoc(id, docData)

	return docData
}
