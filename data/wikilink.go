package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type GTarget string
type GTargetAndLoc struct {
	target GTarget
	loc    *lsp.Location
}

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

	return s.gLinkStore.AddDef(gtarget, location)
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

	return s.gLinkStore.AddDef(gtarget, location)
}

// returns ok
func (s *Store) RemoveGTarget(node *tree_sitter.Node, uri lsp.DocumentURI, content *string) bool {
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

	return s.gLinkStore.RemoveDef(gtarget, location)
}

func (s *Store) GetWikiCompletions(isWikiEnd bool, uri *lsp.DocumentURI) []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}

	for _, t := range s.gLinkStore.GetTargets() {
		var link string
		if isWikiEnd {
			link = "[[" + string(t.target) + "]]"
		} else {
			link = "[[" + string(t.target)
		}
		var excerpt string
		if t.loc != nil {
			excerpt = s.GetExcerpt(*t.loc)
		}
		completions = append(completions, lsp.CompletionItem{
			Label:  link,
			Kind:   lsp.ReferenceCompletion,
			Detail: excerpt,
		})
	}

	return completions
}

func (s *Store) GetGTargetDefinition(target GTarget) []lsp.Location {
	locs, _ := s.gLinkStore.GetDefs(target)
	return locs
}

func (s *Store) GetGTargetHeadingHover(target GTarget) string {
	var totalRefs int

	refs, _ := s.gLinkStore.GetRefs(target)
	totalRefs = len(refs)
	content := fmt.Sprintf("%d references found", totalRefs)
	return content
}

func (s *Store) GetGTargetWikilinkHover(target GTarget) string {
	content := ""
	refs, found := s.gLinkStore.GetRefs(target)
	if found {
		content = fmt.Sprintf("%d references found\n", len(refs)) + content
	}

	defs, found := s.gLinkStore.GetDefs(target)
	if !found {
		content = "No definition found."
	} else {
		if len(defs) == 1 {
			loc := defs[0]
			excerpt := s.GetExcerpt(loc)
			content = fmt.Sprintf("%s\n---\n%s", content, excerpt)
		} else if len(defs) > 1 {
			content = fmt.Sprintf("%d definitions found\n---\n", len(defs))
			for _, loc := range defs {
				excerpt := s.GetExcerpt(loc)
				content = content + fmt.Sprintf("\n%s\n---", excerpt)
			}
		}

	}
	slog.Info(fmt.Sprintf("Refmd con %d", len(refs)))
	if len(refs) > 0 {
		// references md
		var refmd string
		for _, loc := range refs {
			refmd = fmt.Sprintf("%s1 %s\n", refmd, loc.URI.GetFileName())
		}
		slog.Info("Refmd is " + refmd)
		content = content +"\n---\n" + refmd
	}

	return content
}

func (s *Store) GetGTargetReferences(target GTarget) []lsp.Location {
	refs, _ := s.gLinkStore.GetRefs(target)
	return refs
}
