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

// simply gives wiki links [[gtarget]] by \n
// func (s *Store) GetMdReferences(refs []lsp.Location) string {
// 	var content string
//
// 	for _, loc := range refs {
// 		lsp.GetNodeContent(loc.) // dont' have startByte endByte
// 	}
// }

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
	return ""
}

func (s *Store) GetGTargetWikilinkHover(target GTarget, excerptLength int16) string {
	defs, found := s.gLinkStore.GetDefs(target)
	content := ""
	if !found {
		content = "No definition found."
	} else {
		if len(defs) == 1 {
			loc := defs[0]
			excerpt := getWikiHoverExcerpt(loc, excerptLength)
			content = fmt.Sprintf("%s\n---\n%s", content, excerpt)
			return content
		}

		if len(defs) > 1 {
			content = fmt.Sprintf("%d definitions found\n---", len(defs))
		}

		for _, loc := range defs {
			excerpt := getWikiHoverExcerpt(loc, excerptLength)
			content = content + fmt.Sprintf("\n%s\n---", excerpt)
		}
	}

	refs, found := s.gLinkStore.GetRefs(target)
	if found {
		content = content + fmt.Sprintf("%d references found", len(refs))
	}

	return content
}

func (s *Store) GetGTargetReferences(target GTarget) []lsp.Location {
	slog.Info("getting target refs for "+string(target))
	refs, _ := s.gLinkStore.GetRefs(target)
	return refs
}
