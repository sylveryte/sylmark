package data

import (
	"bytes"
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


func SplitLines(content string) {

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
	if err != nil {
		return "", err
	}
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

	return lsp.DocumentURI(u.String()), nil
}

func ContentFromDocPath(mdDocPath string) string {
	slog.Info("reading file " + mdDocPath)
	contentByte, err := os.ReadFile(mdDocPath)
	if err != nil {
		slog.Error("Failed to read file " + mdDocPath + err.Error())
	}
	return string(contentByte)
}

func GetExcerpt(content string, rng lsp.Range, excerptLength int16) string {
	lines := bytes.Split([]byte(content), []byte("\n"))
	startLine := rng.Start.Line
	endLine := startLine + int(excerptLength)

	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	exLines := lines[startLine:endLine]
	return string(bytes.Join(exLines, []byte("\n")))
}
