package data

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"
	"time"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Store struct {
	Tags          map[Tag][]lsp.Location
	GLinkStore    GLinkStore
	DocStore      DocumentStore
	DateStore     DateStore
	LastOpenedDoc lsp.DocumentURI
	ExcerptLength int16
	Config        Config
	MdFiles       []string
	OtherFiles    []string
}

func NewStore() Store {
	return Store{
		Tags:          map[Tag][]lsp.Location{},
		GLinkStore:    NewGlinkStore(),
		DocStore:      NewDocumentStore(),
		DateStore:     NewDateStore(),
		MdFiles:       []string{},
		OtherFiles:    []string{},
		Config:        NewConfig(),
		ExcerptLength: 10,
	}
}

func (s *Store) SyncChangedDocument(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent, parse lsp.ParseFunction) {

	var updatedDocData, oldDocData DocumentData
	// update data into openedDocs
	doc := Document(changes.Text)
	staleDoc, ok := s.GetDocMustTree(uri, parse)
	oldDocData = staleDoc
	if !ok {
		slog.Error("Failed to get old file" + string(uri))
		return
	}
	updatedDocData = *s.UpdateAndReloadDoc(uri, string(doc), parse)
	s.UnloadData(uri, string(oldDocData.Content), oldDocData.Trees)
	s.LoadData(uri, string(updatedDocData.Content), updatedDocData.Trees)

	// sylopti this is for lsp.TDSKIncremental
	// TextDocumentSync is set to TDSKFull so this case won't be there but in future let's implment partial for better perf
	// slog.Info(fmt.Sprintf("RangeLength changed %d", changes.RangeLength))
	// slog.Error("Need to handle partial change text")
	// // return
	// slog.Info("Contents " + changes.Text)
	// slog.Info(fmt.Sprintf("range length %d", changes.RangeLength))
	// slog.Info(fmt.Sprintf(
	// 	"range line start %d end %d",
	// 	changes.Range.Start.Line,
	// 	changes.Range.End.Line,
	// ))
	// slog.Info(fmt.Sprintf(
	// 	"range char start %d end %d",
	// 	changes.Range.Start.Character,
	// 	changes.Range.End.Character,
	// ))
}

func (store *Store) UnloadData(uri lsp.DocumentURI, content string, trees *lsp.Trees) {
	lsp.TraverseNodeWith(trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "atx_heading":
			{
				store.RemoveGTarget(n, uri, &content)
			}
		}
	})
	lsp.TraverseNodeWith(trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wiki_link":
			{
				target, ok := GetWikilinkTarget(n, content, uri)
				if ok {
					loc := uri.LocationOf(n)
					store.GLinkStore.RemoveRef(target, loc)
				}
			}
		case "tag":
			{
				store.RemoveTag(n, uri, &content)
			}
		}
	})
}

func (store *Store) LoadData(uri lsp.DocumentURI, content string, trees *lsp.Trees) {

	store.AddFileGTarget(uri)
	lsp.TraverseNodeWith(trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "atx_heading":
			{
				store.AddGTarget(n, uri, &content)
			}
		}
	})

	lsp.TraverseNodeWith(trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wiki_link":
			{
				target, ok := GetWikilinkTarget(n, content, uri)
				if ok {
					isSubheading := len(target) > 0 && target[0] == '#'
					if !isSubheading {
						loc := uri.LocationOf(n)
						store.GLinkStore.AddRef(target, loc)
					}
				}
			}
		case "tag":
			{
				store.AddTag(n, uri, &content)
			}
		}
	})
}

func (store *Store) UpdateAndReloadDoc(uri lsp.DocumentURI, content string, parse lsp.ParseFunction) *DocumentData {
	t := time.Now()
	trees := parse(content, nil)
	doc := Document(content)
	slog.Info(fmt.Sprintf("%dms<==parsing time", time.Since(t).Milliseconds()))

	// slog.Info("First main---------------")
	// lsp.PrintTsTree(*trees.GetMainTree().RootNode(), 0, content)
	// slog.Info("Now inline-------------")
	// lsp.PrintTsTree(*trees.GetInlineTree().RootNode(), 0, content)

	docData := NewDocumentData(doc, trees)

	// important to update doc first since GetLoadedDataStore fetches it
	store.AddUpdateDoc(uri, docData)
	docData.Headings = store.GetLoadedDataStore(uri, parse)
	docData.FootNotes = store.GetLoadedFootNotesStore(uri, parse)
	// finally update with stores
	store.AddUpdateDoc(uri, docData)

	return docData
}
