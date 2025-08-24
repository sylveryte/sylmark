package data

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type DocumentStore map[lsp.DocumentURI]DocumentData

func NewDocumentStore() DocumentStore {
	return map[lsp.DocumentURI]DocumentData{}
}
func (store *Store) RemoveDoc(uri lsp.DocumentURI) (docData DocumentData, found bool) {
	if store == nil {
		slog.Error("DocumentStore is empty")
		return DocumentData{}, false
	}
	s := *store
	docData, found = s.DocStore[uri]
	if found {
		delete(s.DocStore, uri)
	}
	return docData, found
}

// returns ok
func (store *Store) AddDoc(uri lsp.DocumentURI, doc Document, tree *tree_sitter.Tree) bool {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return false
	}
	s := *store

	s.DocStore[uri] = *newDocumentData(doc, tree)
	return true
}

func (store *Store) UpdateDoc(uri lsp.DocumentURI, change lsp.TextDocumentContentChangeEvent, parse func(content string) *tree_sitter.Tree) (newDocData DocumentData, oldDocData DocumentData, ok bool) {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return DocumentData{}, DocumentData{}, false
	}
	s := *store

	if change.RangeLength == 0 {
		doc := Document(change.Text)
		tree := parse(change.Text)
		oldDocData := s.DocStore[uri]
		newDocData := *newDocumentData(doc, tree)
		s.DocStore[uri] = newDocData
		return newDocData, oldDocData, true
	} else {
		// syltodo TODO 👷 cleanup
		slog.Error("Need to handle partial change text")
		slog.Info("Contents " + change.Text)
		slog.Info(fmt.Sprintf("range length %d", change.RangeLength))
		slog.Info(fmt.Sprintf(
			"range start %d end %d",
			change.Range.Start.Line,
			change.Range.End.Line,
		))

		return DocumentData{}, DocumentData{}, false
	}

}

func (store *Store) DocDataFromURI(uri lsp.DocumentURI) (docData DocumentData, found bool) {
	if store == nil {
		return DocumentData{}, false
	}
	s := *store

	docData, found = s.DocStore[uri]
	if !found {
		path, err := PathFromURI(uri)
		if err != nil {
			slog.Error("failed to PathFromURI=" + string(uri))
			return
		}
		content := ContentFromDocPath(path)
		fdocData := newDocumentData(Document(content), nil)
		if fdocData != nil {
			s.DocStore[uri] = *fdocData
			docData = *fdocData
			return docData, true
		}
	}
	return docData, found
}
