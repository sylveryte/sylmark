package data

import (
	"log/slog"
	"sylmark/lsp"
)

type DocumentStore map[Id]DocumentData

func NewDocumentStore() DocumentStore {
	return map[Id]DocumentData{}
}

// removes from DocStore and GLinkStore
func (store *Store) RemoveDoc(id Id) (docData DocumentData, found bool) {
	if store == nil {
		slog.Error("Store is empty")
		return DocumentData{}, false
	}
	s := *store
	// remove from DocStore
	docData, found = s.DocStore[id]
	if found {
		delete(s.DocStore, id)
	}
	// remove from gliGLinkStore
	s.LinkStore.RemoveDef(id, "", lsp.Range{})
	return docData, found
}

// returns ok
func (store *Store) AddUpdateDoc(id Id, docData *DocumentData) bool {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return false
	}
	s := *store

	s.DocStore[id] = *docData
	return true
}
func (s *Store) GetDocMustTree(id Id, parse lsp.ParseFunction) (docData DocumentData, ok bool) {
	docData, found := s.GetDoc(id)
	if found {
		if docData.Trees == nil {
			docData.Trees = parse(string(docData.Content), nil)
			s.AddUpdateDoc(id, &docData)
		}
		return docData, true
	}
	return docData, false
}

func (s *Store) AddDoc(id Id) *DocumentData {
	uri, _ := s.GetUri(id)
	path, err := PathFromURI(uri)
	if err != nil {
		slog.Error("failed to PathFromURI=" + string(uri))
		return nil
	}
	content := ContentFromDocPath(path)
	fdocData := NewDocumentData(Document(content), nil)
	s.DocStore[id] = *fdocData
	return fdocData
}

func (s *Store) GetDoc(id Id) (docData DocumentData, ok bool) {
	docData, found := s.DocStore[id]
	if !found {
		doc := s.AddDoc(id)
		if doc == nil {
			return
		}
		docData = *doc
	}
	return docData, true
}
