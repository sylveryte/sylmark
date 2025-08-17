package data

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"sylmark/lsp"

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

type DocumentStore map[lsp.DocumentURI]DocumentData

func NewDocumentStore() DocumentStore {
	return map[lsp.DocumentURI]DocumentData{}
}


// syltodo TODO optimize it
func (doc *Document) GetLine(lineNumber int) string {
	for i, v := range bytes.Split([]byte(*doc), []byte("\n")) {
		if i == lineNumber {
			return string(v)
		}
	}
	return ""
}

func (store *DocumentStore) RemoveDoc(uri lsp.DocumentURI) (docData DocumentData, found bool) {
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
func (store *DocumentStore) AddDoc(uri lsp.DocumentURI, doc Document, tree *tree_sitter.Tree) bool {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return false
	}
	s := *store

	s[uri] = *newDocumentData(doc, tree)
	return true
}

func (store *DocumentStore) UpdateDoc(uri lsp.DocumentURI, change lsp.TextDocumentContentChangeEvent, parse func(content string) *tree_sitter.Tree) (newDocData DocumentData, oldDocData DocumentData, ok bool) {
	if store == nil {
		slog.Error("DocumentStore not defined")
		return DocumentData{}, DocumentData{}, false
	}
	s := *store

	if change.RangeLength == 0 {
		doc := Document(change.Text)
		tree := parse(change.Text)
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

func (store *DocumentStore) DocDataFromURI(uri lsp.DocumentURI) (docData DocumentData, found bool) {
	if store == nil {
		return DocumentData{}, false
	}
	s := *store

	docData, found = s[uri]
	return
}

func DirPathFromURI(uri lsp.DocumentURI) (path string, er error) {
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

func UriFromPath(path string) (lsp.DocumentURI, error) {

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	uriPath := filepath.ToSlash(absPath)
	u := url.URL{
		Scheme: "file",
		Path:   uriPath,
	}

	return lsp.DocumentURI(u.String()), nil
}

func ContentFromDocPath(mdDocPath string) string {
	contentByte, err := os.ReadFile(mdDocPath)
	if err != nil {
		slog.Error("Failed to read file " + mdDocPath + err.Error())
	}
	return string(contentByte)
}
