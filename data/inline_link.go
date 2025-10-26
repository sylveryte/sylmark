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
	doesAliasExists := len(text) != 0
	completions := []lsp.CompletionItem{}
	mdFilesOnly := len(arg) > 1 && arg[0] == ' '
	otherFilesOnly := mdFilesOnly && len(arg) > 2 && arg[1] == ' '
	strppedArg := strings.TrimSpace(arg)
	isHeadingMode := strings.ContainsRune(arg, '#')

	if isHeadingMode {
		splits := strings.SplitN(arg, "#", 2)
		filePath := splits[0]
		subTarget := splits[1]
		fullFilePath, err := GetFullPathRelatedTo(*uri, filePath)
		if err != nil {
			slog.Warn(fmt.Sprintf("Link file issue %s", err.Error()))
		}
		targetUri, err := UriFromPath(fullFilePath)
		if err != nil {
			slog.Warn(fmt.Sprintf("URI failed file issue %s", err.Error()))
		}
		targetId := s.GetIdFromURI(targetUri)
		subTargets := s.LinkStore.GetSubTargetsAndRanges(targetId)
		for _, subTargetNRange := range subTargets {
			if len(subTargetNRange.subTarget) == 0 {
				continue
			}
			fullLink := FullTarget(string(filePath) + string(subTargetNRange.subTarget))
			match := fuzzy.MatchFold(arg+subTarget, string(fullLink))
			var link string
			if match {
				link = fmt.Sprintf("[%s](%s)", text, fullLink)
			}
			completions = append(completions, lsp.CompletionItem{
				Label:    link,
				Kind:     lsp.ReferenceCompletion,
				SortText: "b",
				TextEdit: &lsp.TextEdit{
					Range:   rng,
					NewText: link,
				},
				Detail: "",
			})
		}
	} else {

		includeMdFiles := !otherFilesOnly
		includeOtherFiles := !mdFilesOnly

		// md files
		var files []string
		if includeMdFiles {
			for u := range s.IdStore.uri {
				f, er := PathFromURI(u)
				if er != nil {
					continue
				}
				files = append(files, f)
			}
		}
		if includeOtherFiles {
			files = append(files, s.OtherFiles...)
		}
		for _, path := range files {
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
			if doesAliasExists == false {
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
	if n == nil {
		return "", fmt.Errorf("No link")
	}
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

func GetUriFromInlineNode(inlineNode *tree_sitter.Node, content string, relUri lsp.DocumentURI) (lsp.DocumentURI, bool) {

	filePath, err := GetInlineLinkTarget(inlineNode, content, relUri)
	if err != nil {
		slog.Error("File doesnt exist")
		return "", false
	}
	fullFilePath, err := GetFullPathRelatedTo(relUri, filePath)
	if err != nil {
		slog.Error("Fialed to get full path" + err.Error())
		return "", false
	}
	uri, err := UriFromPath(fullFilePath)
	if err != nil {
		slog.Error("Failed to make uri " + err.Error())
		return "", false
	}
	return uri, true
}

func IsMdFile(path string) bool {
	return strings.HasSuffix(path, ".md")
}
func GetUrlAndSubTarget(fullUrl string) (url lsp.DocumentURI, subTarget SubTarget, found bool) {
	found = strings.ContainsRune(fullUrl, '#')
	if found {
		splits := strings.SplitN(fullUrl, "#", 2)
		url = lsp.DocumentURI(splits[0])
		subTarget = SubTarget("#" + splits[1])
		return url, subTarget, found
	} else {
		return lsp.DocumentURI(fullUrl), subTarget, found
	}
}
