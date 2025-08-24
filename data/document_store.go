package data

import (
	"log/slog"
	"sylmark/lsp"
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
func (store *Store) AddUpdateDoc(uri lsp.DocumentURI, docData *DocumentData) bool {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return false
	}
	s := *store

	s.DocStore[uri] = *docData
	return true
}

func (store *Store) GetDoc(uri lsp.DocumentURI) (docData DocumentData, ok bool) {
	if store == nil {
		return DocumentData{}, false
	}
	s := *store

	docData, found := s.DocStore[uri]
	if !found {
		path, err := PathFromURI(uri)
		if err != nil {
			slog.Error("failed to PathFromURI=" + string(uri))
			return
		}
		content := ContentFromDocPath(path)
		fdocData := NewDocumentData(Document(content), nil)
		if fdocData != nil {
			s.DocStore[uri] = *fdocData
			docData = *fdocData
		}
	}
	return docData, true
}
