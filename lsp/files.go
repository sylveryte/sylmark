package lsp

import (
	"net/url"
	"os"
	"path/filepath"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (h *LangHandler) addRootPathAndLoad(dir string) {
	h.rootPath = dir
	h.loadInactiveData()
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

func getLocation(uri DocumentURI, node *tree_sitter.Node) Location {

	return Location{
		URI:   uri,
		Range: getRange(node),
	}

}
