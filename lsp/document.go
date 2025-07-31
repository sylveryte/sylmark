package lsp

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Document string
type DocumentData struct {
	Tree    *tree_sitter.Tree
	Content Document
}

func newDocumentData(doc Document, tree *tree_sitter.Tree) *DocumentData {
	return &DocumentData{
		Content: doc,
		Tree:    tree,
	}
}

type DocumentStore map[DocumentURI]DocumentData

func newDocumentStore() DocumentStore {
	return map[DocumentURI]DocumentData{}
}
func (store *DocumentStore) removeDoc(uri DocumentURI) (docData DocumentData, found bool) {
	if store == nil {
		slog.Error("DocumentStore is empty")
		return DocumentData{}, false
	}
	s := *store
	docData, found = s[uri]
	if found {
		delete(s, uri)
	}
	return docData, found
}

// returns ok
func (store *DocumentStore) addDoc(uri DocumentURI, doc Document, tree *tree_sitter.Tree) bool {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return false
	}
	s := *store

	s[uri] = *newDocumentData(doc, tree)
	return true
}

func (store *DocumentStore) updateDoc(uri DocumentURI, change TextDocumentContentChangeEvent, h *LangHandler) (newDocData DocumentData, oldDocData DocumentData, ok bool) {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return DocumentData{}, DocumentData{}, false
	}
	s := *store

	if change.RangeLength == 0 {
		doc := Document(change.Text)
		tree := h.parse(change.Text)
		oldDocData := s[uri]
		newDocData := *newDocumentData(doc, tree)
		s[uri] = newDocData
		return newDocData, oldDocData, true
	} else {
		// syltodo TODO 👷
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

func (store *DocumentStore) docDataFromURI(uri DocumentURI) (docData DocumentData, found bool) {
	if store == nil {
		return DocumentData{}, false
	}
	s := *store

	docData, found = s[uri]
	return
}

func dirPathFromURI(uri DocumentURI) (path string, er error) {
	parsedUrl, err := url.Parse(string(uri))
	if err != nil {
		return "", err
	}

	filePath := parsedUrl.Path
	dir := filePath
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		dir = filepath.Dir(filePath)
	}

	dir = filepath.Clean(dir)
	return dir, nil
}

func uriFromPath(path string) (DocumentURI, error) {

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	uriPath := filepath.ToSlash(absPath)
	u := url.URL{
		Scheme: "file",
		Path:   uriPath,
	}

	return DocumentURI(u.String()), nil
}

func locationFromURINode(uri DocumentURI, node *tree_sitter.Node) Location {

	return Location{
		URI:   uri,
		Range: getRange(node),
	}

}
func contentFromDocPath(mdDocPath string) string {
	contentByte, err := os.ReadFile(mdDocPath)
	if err != nil {
		slog.Error("Failed to read file " + mdDocPath + err.Error())
	}
	return string(contentByte)
}
