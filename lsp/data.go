package lsp

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Target string
type Tag string

type DocumentData struct {
	Node    *tree_sitter.Node
	Content *string
}

type NodeData struct {
	URI  DocumentURI
	Node *tree_sitter.Node
}

func (h *LangHandler) loadInactiveData() {
	if h.rootPath == "" {
		slog.Error("h.rootPath is empty")
		return
	}

	filepath.WalkDir(h.rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.HasSuffix(path, ".git") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			// slog.Info(fmt.Sprintf("MDFile=%s Name=%s", path, d.Name()))
			h.loadInactiveDataOfFile(path)
		}

		return nil
	})

}

func (h *LangHandler) loadInactiveDataOfFile(mdFilePath string) {
	contentByte, err := os.ReadFile(mdFilePath)
	if err != nil {
		slog.Error("Failed to read file " + mdFilePath + err.Error())
	}
	content := string(contentByte)
	tree := h.parse(content)
	defer tree.Close()

	rootNode := tree.RootNode()

	uri, err := uriFromPath(mdFilePath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	TraverseNodeWith(rootNode, func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "wikilink":
			{
				getWikilinkLink(n, content)
			}
		case "heading":
			{
				getHeadingTitle(n, content)
			}
		case "tag":
			{
				h.inactiveStore.AddTag(n, uri, &content)
			}
		case "link":
			{
				getLinkUrl(n, content)
			}
		}
	})

	// slog.Info(fmt.Sprintf("Root node is %d", rootNode.ChildCount()))
}
