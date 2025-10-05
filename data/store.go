package data

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"

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
	if changes.RangeLength == 0 {
		doc := Document(changes.Text)
		staleDoc, ok := s.GetDocMustTree(uri, parse)
		oldDocData = staleDoc
		if !ok {
			slog.Error("Failed to get old file" + string(uri))
			return
		}
		// sylopti we can use oldDocData.tree to optimze it but initial try met with some weird issues, wrong tree study and use
		tree := parse(changes.Text, nil)
		updatedDocData = *NewDocumentData(doc, tree)
		updatedDocData.Headings = s.GetLoadedDataStore(uri, parse)
		updatedDocData.FootNotes = s.GetLoadedFootNotesStore(uri, parse)
		s.AddUpdateDoc(uri, &updatedDocData)
	} else {
		// sylopti
		// TextDocumentSync is set to TDSKFull so this case won't be there but in future let's implment partial for better perf
		slog.Error("Need to handle partial change text")
		// return
		slog.Info("Contents " + changes.Text)
		slog.Info(fmt.Sprintf("range length %d", changes.RangeLength))
		slog.Info(fmt.Sprintf(
			"range start %d end %d",
			changes.Range.Start.Line,
			changes.Range.End.Line,
		))
	}

	s.UnloadData(uri, string(oldDocData.Content), oldDocData.Trees)
	s.LoadData(uri, string(updatedDocData.Content), updatedDocData.Trees)
}

func (store *Store) UnloadData(uri lsp.DocumentURI, content string, trees *lsp.Trees) {
	lsp.TraverseNodeWith(trees.GetMainTree().RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "heading":
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
