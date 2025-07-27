package lsp

import (
	"fmt"
	"log/slog"
)

type Document string

type DocumentStore map[DocumentURI]Document

func newDocumentStore() DocumentStore {
	return map[DocumentURI]Document{}
}

// returns ok
func (s DocumentStore) AddDoc(uri DocumentURI, doc Document) bool {
	if s == nil {
		return false
	}

	existingDoc, ok := s[uri]
	if ok {
		slog.Info("There exists docs for " + string(uri) + string(existingDoc))
	}
	s[uri] = doc

	return true
}

// returns ok
func (s DocumentStore) UpdateDoc(uri DocumentURI, change TextDocumentContentChangeEvent) bool {

	existingDoc, ok := s[uri]
	if ok {
		slog.Info("There exists docs for " + string(uri) + string(existingDoc))
	}

	slog.Info("Contents " + change.Text)
	slog.Info(fmt.Sprintf("range length %d", change.RangeLength))
	slog.Info(fmt.Sprintf("range start %d end %d", change.Range.Start.Line, change.Range.End.Line))

	return true
}

func (s DocumentStore) GetDoc(uri DocumentURI) (doc Document, found bool) {
	if s == nil {
		return "", false
	}
	doc, found = s[uri]
	return
}
