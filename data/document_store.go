package data

import (
	"log/slog"
	"sylmark/lsp"
)

type DocumentStore map[lsp.DocumentURI]DocumentData

func NewDocumentStore() DocumentStore {
	return map[lsp.DocumentURI]DocumentData{}
}

// removes from DocStore and GLinkStore
func (store *Store) RemoveDoc(uri lsp.DocumentURI) (docData DocumentData, found bool) {
	if store == nil {
		slog.Error("DocumentStore is empty")
		return DocumentData{}, false
	}
	s := *store
	// remove from DocStore
	docData, found = s.DocStore[uri]
	if found {
		delete(s.DocStore, uri)
	}
	// remove from gliGLinkStore
	gtarget, ok := GetFileGTarget(uri)
	if ok {
		s.GLinkStore.RemoveDef(GTarget(gtarget), lsp.Location{
			URI:   uri,
			Range: lsp.Range{},
		})
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
func (store *Store) GetDocMustTree(uri lsp.DocumentURI, parse lsp.ParseFunction) (docData DocumentData, ok bool) {
	docData, found := store.GetDoc(uri)
	if found {
		if docData.Trees == nil {
			docData.Trees = parse(string(docData.Content), nil)
			store.AddUpdateDoc(uri, &docData)
		}
		return docData, true
	}
	return docData, false
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
