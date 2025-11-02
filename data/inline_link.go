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
	fileId := s.GetIdFromURI(*uri)
	sourcePath, err := DirPathFromURI(*uri)
	if err != nil {
		slog.Error("Something went wrong for path relative " + err.Error())
		return completions
	}

	if isHeadingMode {
		arg = s.Config.ProcessInlineTargetPath(arg)
		filePathUri, targetId, _, _ := s.GetInlineTargetAndSubTarget(arg, fileId)
		filePath, _ := PathFromURI(filePathUri)
		subTargets := s.LinkStore.getSubTargetsAndRanges(targetId)
		for _, subTargetNRange := range subTargets {
			if len(subTargetNRange.subTarget) == 0 {
				continue
			}
			encodedRelPath, _ := s.getInlineRelFormattedTarget(sourcePath, string(filePath))
			fullLink := FullTarget(string(encodedRelPath) + s.encodeForInlineLinkdownLinkPath(string(subTargetNRange.subTarget)))
			var link string
			link = fmt.Sprintf("[%s](%s)", text, fullLink)
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
			encodedRelPath, err := s.getInlineRelFormattedTarget(sourcePath, path)
			if err != nil {
				slog.Error("Something went wrong for path relative " + err.Error())
				continue
			}
			fn := GetFileName(path)
			if doesAliasExists == false {
				text = fn
			}
			var link string
			if s.isImage(encodedRelPath) {
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

// replaces " " with "%20" and ToLower for MdLinkWebMode
func (s *Store) encodeForInlineLinkdownLinkPath(path string) string {
	if s.Config.MdLinkWebMode {
		return strings.ToLower(strings.ReplaceAll(path, " ", "-"))
	}
	return strings.ReplaceAll(path, " ", "%20")
}

// replaces "%20" with " "
func (s *Store) DecodeForInlineLinkdownLinkPath(path string) string {
	if s.Config.MdLinkWebMode {
		return strings.ReplaceAll(path, "-", " ")
	}
	return strings.ReplaceAll(path, "%20", " ")
}
func RemoveMdExtOnly(fileName string) string {
	return strings.TrimSuffix(fileName, ".md")
}
func (s *Store) isImage(filePath string) bool {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".png", ".jpeg", ".gif", ".webp", ".avif":

		return true
	}

	return false
}

func (c *Config) GetInlineLinkTarget(node *tree_sitter.Node, content string) (path string, err error) {
	if node.Kind() != "inline_link" {
		return "", fmt.Errorf("Not inline_link")
	}
	n := node.NamedChild(1)
	if n == nil {
		return "", fmt.Errorf("No link")
	}
	path = lsp.GetNodeContent(*n, content)
	path = c.ProcessInlineTargetPath(path)
	return path, nil

}

func (c *Config) ProcessInlineTargetPath(path string) string {
	if c.MdLinkWebMode && strings.HasPrefix(path, "..") {
		cleanFilePath, found := strings.CutPrefix(path, "..")
		if found {
			path = cleanFilePath
		}
	}
	return path
}

func GetFullPathRelatedTo(fullURI lsp.DocumentURI, filePath string) (string, error) {
	dir, err := DirPathFromURI(fullURI)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filePath), nil

}

// checks if ends with .md
func IsMdFile(path string) bool {
	return strings.HasSuffix(path, ".md")
}

// adds .md at suffix if needed
func GetInlineTargetUrl(path string) string {
	if strings.HasSuffix(path, ".md") {
		return path
	}
	return path + ".md"
}

// removes .md file if config demands
func (s *Store) GetMdFormattedTargetUrl(path string) string {
	if s.Config.IncludeMdExtensionMdLink && !s.Config.MdLinkWebMode {
		return path
	}
	return RemoveMdExtOnly(path)
}

func (s *Store) getInlineRelFormattedTarget(sourcePath string, path string) (relPath string, err error) {
	relPath, err = filepath.Rel(sourcePath, path)
	if err != nil {
		return
	}
	if s.Config.MdLinkWebMode && !strings.HasPrefix(relPath, "file://") {
		relPath = s.encodeForInlineLinkdownLinkPath(relPath)
		relPath = strings.ToLower(relPath)
		relPath = filepath.Join("..", relPath)
	}
	relPath = s.encodeForInlineLinkdownLinkPath(relPath)
	relPath = s.GetMdFormattedTargetUrl(relPath)

	return relPath, nil
}

// unformats it gets real link adds .md where needed
func (s *Store) GetInlineFullTargetAndSubTarget(n *tree_sitter.Node, content string, fileId Id) (url lsp.DocumentURI, targetId Id, subTarget SubTarget, found bool) {
	fullUrl, err := s.Config.GetInlineLinkTarget(n, content)
	if err != nil {
		return
	}

	return s.GetInlineTargetAndSubTarget(fullUrl, fileId)
}

func (s *Store) GetInlineTargetAndSubTarget(fullUrl string, fileId Id) (url lsp.DocumentURI, targetId Id, subTarget SubTarget, found bool) {

	found = strings.ContainsRune(fullUrl, '#')
	var target string
	var relTarget string

	if found {
		splits := strings.SplitN(fullUrl, "#", 2)
		target = splits[0]
		subTarget = SubTarget("#" + splits[1])
	} else {
		target = fullUrl
	}
	target = GetInlineTargetUrl(target)
	target = s.DecodeForInlineLinkdownLinkPath(target)

	if !strings.HasPrefix(target, "file://") {
		fileUri, _ := s.GetUri(fileId)

		if s.Config.MdLinkWebMode {
			relTarget = target
			target, _ = GetFullPathRelatedTo(fileUri, relTarget)

			targetUri, _ := UriFromPath(target)
			id, found := s.findIdFromURIFold(targetUri)

			targetId = id
			url, found = s.GetUri(id)

			return url, targetId, subTarget, found
		}

		target, _ = GetFullPathRelatedTo(fileUri, target)
	}

	url, _ = UriFromPath(target)
	if IsMdFile(target) {
		found = true
		targetId = s.GetIdFromURI(url)
	}
	return url, targetId, subTarget, found
}
