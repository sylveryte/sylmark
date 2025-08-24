package data

import (
	"log/slog"
	"net/url"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Store struct {
	tags       map[Tag][]lsp.Location
	gLinkStore GLinkStore
	DocStore   DocumentStore
	ExcerptLength int16
}

func NewStore() Store {
	return Store{
		tags:       map[Tag][]lsp.Location{},
		gLinkStore: NewGlinkStore(),
		DocStore:   NewDocumentStore(),
		ExcerptLength: 10,
	}
}

func (s *Store) SyncChangedDocument(uri lsp.DocumentURI, changes lsp.TextDocumentContentChangeEvent, parse lsp.ParseFunction) {

	unscapedUri, err := url.QueryUnescape(string(uri))
	if err == nil {
		uri = lsp.DocumentURI(unscapedUri)
	} else {
		slog.Error("Failed to unscapedUri")
		return
	}

	var updatedDocData, oldDocData DocumentData
	// update data into openedDocs
	if changes.RangeLength == 0 {
		doc := Document(changes.Text)
		staleDoc, ok := s.GetDoc(uri)
		oldDocData = staleDoc
		if !ok {
			slog.Error("Failed to get old file" + string(uri))
			return
		}
		// sylopti we can use oldDocData.tree to optimze it but initial try met with some weird issues, wrong tree study and use
		tree := parse(changes.Text, nil)
		updatedDocData = *NewDocumentData(doc, tree)
		s.AddUpdateDoc(uri, &updatedDocData)
	} else {
		// sylopti
		// TextDocumentSync is set to TDSKFull so this case won't be there but in future let's implment partial for better perf
		slog.Error("Need to handle partial change text")
		return
		// slog.Info("Contents " + change.Text)
		// slog.Info(fmt.Sprintf("range length %d", change.RangeLength))
		// slog.Info(fmt.Sprintf(
		// 	"range start %d end %d",
		// 	change.Range.Start.Line,
		// 	change.Range.End.Line,
		// ))
	}

	// UnloadData
	s.UnloadData(uri, string(oldDocData.Content), oldDocData.Tree.RootNode())

	// LoadData
	s.LoadData(uri, string(updatedDocData.Content), updatedDocData.Tree.RootNode())
}

func (store *Store) UnloadData(uri lsp.DocumentURI, content string, rootNode *tree_sitter.Node) {
	lsp.TraverseNodeWith(rootNode, func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wikilink":
			{
				target, ok := GetWikilinkTarget(n, content, uri)
				if ok {
					loc := uri.LocationOf(n)
					store.gLinkStore.RemoveRef(target, loc)
				}
			}
		case "heading":
			{
				store.RemoveGTarget(n, uri, &content)
			}
		case "tag":
			{
				store.RemoveTag(n, uri, &content)
			}
		}
	})
}

func (store *Store) LoadData(uri lsp.DocumentURI, content string, rootNode *tree_sitter.Node) {

	store.AddFileGTarget(uri)
	lsp.TraverseNodeWith(rootNode, func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wikilink":
			{
				target, ok := GetWikilinkTarget(n, content, uri)
				if ok {
					loc := uri.LocationOf(n)
					store.gLinkStore.AddRef(target, loc)
				}
			}
		case "heading":
			{
				store.AddGTarget(n, uri, &content)
			}
		case "tag":
			{
				store.AddTag(n, uri, &content)
			}
		}
	})
}
