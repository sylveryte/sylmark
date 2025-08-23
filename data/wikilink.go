package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Target string  // is like # Some heading
type GTarget string // is full GLink target
type Tag string

func getFileGTarget(uri lsp.DocumentURI) (gtarget string, ok bool) {

	filename := uri.GetFileName()
	splits := strings.Split(filename, ".md")
	if len(splits) < 1 {
		return "", false
	}
	return splits[0], true
}
func getGTarget(heading string, uri lsp.DocumentURI) (gtarget GTarget, ok bool) {

	fileGtTarget, ok := getFileGTarget(uri)
	if ok {
		if len(heading) > 0 {
			return GTarget(fileGtTarget + "#" + heading), true
		} else {
			return GTarget(fileGtTarget), true
		}
	} else {
		return "", false
	}
}

func (s *Store) AddFileGTarget(uri lsp.DocumentURI) bool {

	gtarget, ok := getGTarget("", uri)
	if !ok {
		slog.Error("Could not form gtarget")
		return false
	}

	location := uri.LocationOfFile()

	return s.gLinkStore.AddDef(gtarget, location.URI, location.Range)
}

// returns ok
func (s *Store) AddGTarget(node *tree_sitter.Node, uri lsp.DocumentURI, content *string) bool {
	if s == nil {
		return false
	}

	heading, ok := getHeadingTitle(node, *content)
	if !ok {
		slog.Error("Could not extract heading")
		return false
	}
	gtarget, ok := getGTarget(heading, uri)
	if !ok {
		slog.Error("Could not form gtarget")
		return false
	}

	location := uri.LocationOf(node)

	return s.gLinkStore.AddDef(gtarget, location.URI, location.Range)
}

func (s *Store) GetWikiCompletions(isWikiEnd bool, uri *lsp.DocumentURI) []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}

	for _, t := range s.gLinkStore.GetTargets() {
		var link string
		if isWikiEnd {
			link = "[[" + string(t) + "]]"
		} else {
			link = "[[" + string(t)
		}
		completions = append(completions, lsp.CompletionItem{
			Label:  link,
			Kind:   lsp.FileCompletion,
			Detail: string(t),
		})
	}

	return completions
}

func (s *Store) GetGTargetDefinition(target GTarget) []lsp.Location {
	locs, _ := s.gLinkStore.GetDefs(target)
	slog.Info(fmt.Sprintf("returing %d ", len(locs)))
	return locs
}

func getWikiHoverExcerpt(loc lsp.Location, excerptLength int16) string {
	filePath, err := PathFromURI(loc.URI)
	if err != nil {
		slog.Error("failed to read file URI = " + string(loc.URI))
	}
	slog.Error("filePath path is = " + string(filePath))
	cont := ContentFromDocPath(filePath)
	return GetExcerpt(cont, loc.Range, excerptLength)
}

func (s *Store) GetWikiHoverReferences(target GTarget) []lsp.Location {
	// glinks, _, gfound := s.gLinkStore.GetGLinks(target)
	// if !gfound {
	// 	return "No references found."
	// }
	// return ""
}

func (s *Store) GetGTargetHeadingHover(target GTarget) string {
	// var totalRefs int
	//
	// glinks, _, gfound := s.gLinkStore.GetGLinks(target)
	// glinks(func(g glink) bool {
	// 	totalRefs = totalRefs + len(g.refs)
	// 	return true
	// })
	//
	// content := fmt.Sprintf("%d references found", totalRefs)
	// return content
}

func (s *Store) GetGTargetWikilinkHover(target GTarget, excerptLength int16) string {
	glinks, gcount, gfound := s.gLinkStore.GetGLinks(target)
	if !gfound {
		return "No definition found."
	}
	content := ""

	if gcount > 1 {
		content = fmt.Sprintf("%d definitions found", gcount)
	}

	glinks(func(g glink) bool {
		loc := g.def
		excerpt := getWikiHoverExcerpt(loc, excerptLength)
		content = fmt.Sprintf("%s\n---\n%s", content, excerpt)
		return true
	})

	return content
}
