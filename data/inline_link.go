package data

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sylmark/lsp"

	"github.com/lithammer/fuzzysearch/fuzzy"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (s *Store) GetInlineLinkCompletions(arg string, text string, rng lsp.Range, uri *lsp.DocumentURI) []lsp.CompletionItem {
	isText := len(text) != 0
	completions := []lsp.CompletionItem{}
	strppedArg := strings.TrimSpace(arg)

	// md files
	allFiles := s.OtherFiles
	allFiles = append(allFiles, s.MdFiles...)
	for _, path := range s.OtherFiles {
		if match := fuzzy.MatchFold(strppedArg, path); match == false {
			continue
		}
		sourcePath, err := DirPathFromURI(*uri)
		if err != nil {
			slog.Error("Something went wrong for path relative " + err.Error())
			continue
		}
		relPath, err := filepath.Rel(sourcePath, path)
		encodedRelPath := encodeForInlineLinkdownLinkPath(relPath)
		if err != nil {
			slog.Error("Something went wrong for path relative " + err.Error())
			continue
		}
		fn := GetFileName(path)
		if isText == false {
			text = fn
		}
		var link string
		if s.isImage(relPath) {
			link = fmt.Sprintf("![%s](%s)", text, encodedRelPath)
		} else {
			link = fmt.Sprintf("[%s](%s)", text, encodedRelPath)
		}
		completions = append(completions, lsp.CompletionItem{
			Label:    link,
			Kind:     lsp.FileCompletion,
			SortText: "a",
			TextEdit: &lsp.TextEdit{
				Range:   rng,
				NewText: link,
			},
			Detail: "",
		})

	}
	return completions
}

func GetFileName(path string) string {
	return RemoveMdExtOnly(filepath.Base(path))
}

// replaces " " with "%20"
func encodeForInlineLinkdownLinkPath(path string) string {
	return strings.ReplaceAll(path, " ", "%20")
}

// replaces "%20" with " "
func DecodeForInlineLinkdownLinkPath(path string) string {
	return strings.ReplaceAll(path, "%20", " ")
}
func RemoveMdExtOnly(fileName string) string {
	return strings.ReplaceAll(fileName, ".md", "")
}
func (s *Store) isImage(filePath string) bool {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".png", ".jpeg", ".gif", ".webp", ".avif":

		return true
	}

	return false
}

func GetInlineLinkTarget(node *tree_sitter.Node, content string, uri lsp.DocumentURI) (path string, err error) {
	if node.Kind() != "inline_link" {
		return "", fmt.Errorf("Not inline_link")
	}
	n := node.NamedChild(1)
	path = lsp.GetNodeContent(*n, content)
	return path, nil

}

func GetFullPathRelatedTo(fullURI lsp.DocumentURI, filePath string) (string, error) {
	dir, err := DirPathFromURI(fullURI)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filePath), nil

}
