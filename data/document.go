package data

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"sylmark-server/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Document string
type DocumentData struct {
	Tree    *tree_sitter.Tree
	Content Document
}

func NewDocumentData(doc Document, tree *tree_sitter.Tree) *DocumentData {
	return &DocumentData{
		Content: doc,
		Tree:    tree,
	}
}

// sylopti
func (doc *Document) GetLine(lineNumber int) string {
	for i, v := range bytes.Split([]byte(*doc), []byte("\n")) {
		if i == lineNumber {
			return string(v)
		}
	}
	return ""
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


func GetDirPathFromURI(uri lsp.DocumentURI) (string, error) {
	path, err := PathFromURI(uri)
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

func PathFromURI(uri lsp.DocumentURI) (string, error) {
	parsedUrl, err := url.Parse(string(uri))
	if err != nil {
		return "", err
	}

	return parsedUrl.Path, nil
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
	// remove %20
	cleanedURI, err := CleanUpURI(u.String())

	return lsp.DocumentURI(cleanedURI), err
}

func CleanUpURI(uri string) (lsp.DocumentURI, error) {
	// remove %20
	unscapedUri, err := url.QueryUnescape(uri)

	return lsp.DocumentURI(unscapedUri), err
}

func ContentFromDocPath(mdDocPath string) string {
	contentByte, err := os.ReadFile(mdDocPath)
	if err != nil {
		slog.Error("Failed to read file " + mdDocPath + err.Error())
	}
	return string(contentByte)
}

func (store *Store) GetExcerpt(loc lsp.Location) string {
	docData, ok := store.GetDoc(loc.URI)
	if !ok {
		slog.Error("Failed to get doc for GetExcerpt" + string(loc.URI))
		return ""
	}

	lines := bytes.Split([]byte(docData.Content), []byte("\n"))
	startLine := loc.Range.Start.Line
	endLine := startLine + int(store.ExcerptLength)

	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	exLines := lines[startLine:endLine]
	return fmt.Sprintf("\n`Preview`\n\n%s\n`...`\n", string(bytes.Join(exLines, []byte("\n"))))
}
