package data

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sylmark/lsp"
)

type Document string
type DocumentData struct {
	Trees     *lsp.Trees
	Content   Document
	Headings  *HeadingsStore
	FootNotes *FootNotesStore
}

func NewDocumentData(doc Document, trees *lsp.Trees) *DocumentData {
	return &DocumentData{
		Content: doc,
		Trees:   trees,
	}
}

// sylopti
func GetLineFromContent(content string, lineNumber int) (string, bool) {
	lines := strings.Split(content, "\n")
	if len(lines) > lineNumber {
		return lines[lineNumber], true
	}
	return "", false
}

func (doc *Document) GetLine(lineNumber int) string {
	line, _ := GetLineFromContent(string(*doc), lineNumber)
	return line
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

func GetFileURIInSameURIPath(fileName string, baseURI lsp.DocumentURI) (uri lsp.DocumentURI, err error) {
	dir, err := DirPathFromURI(baseURI)
	urlPath := filepath.Join(dir, fileName)
	return UriFromPath(urlPath)
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

// full path is relative to workspace
func (s *Store) GetPathRelRoot(uri lsp.DocumentURI) (relPath string, err error) {
	path, err := PathFromURI(uri)
	relPath, err = filepath.Rel(s.Config.RootPath, path)
	return
}

func (s *Store) GetExcerpt(id Id, rng lsp.Range) string {
	docData, ok := s.GetDoc(id)
	if !ok {
		slog.Error("Failed to get doc for GetExcerpt" + string(id))
		return ""
	}

	lines := bytes.Split([]byte(docData.Content), []byte("\n"))
	startLine := rng.Start.Line
	endLine := startLine + int(s.ExcerptLength)

	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	// syltodo make it safe
	exLines := lines[startLine:endLine]
	return fmt.Sprintf("\n`Preview`\n\n%s\n`...`\n", string(bytes.Join(exLines, []byte("\n"))))
}
