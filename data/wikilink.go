package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"

	"github.com/lithammer/fuzzysearch/fuzzy"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type GTarget string
type GTargetAndLoc struct {
	target GTarget
	loc    *lsp.Location
}

func (t GTarget) SplitHeading() (gTarget GTarget, heading string, hasHeading bool) {
	ts := string(t)
	if strings.ContainsRune(ts, '#') {
		splits := strings.Split(ts, "#")
		return GTarget(splits[0]), "# " + splits[1], true
	}
	return t, "", false
}
func (t GTarget) GetWIthinTarget() (target GTarget,  hasHeading bool) {
	ts := string(t)
	if strings.ContainsRune(ts, '#') {
		splits := strings.Split(ts, "#")
		return GTarget("#" + splits[1]), true
	}
	return "", false
}

func (t GTarget) GetFileName() (fileName string, heading string, hasHeading bool) {
	t.SplitHeading()
	fileTarget, heading, hasHeading := t.SplitHeading()
	if hasHeading {
		return fileTarget.GetFileName()
	} else {
		return string(t) + ".md", "", false
	}
}

func GetFileGTarget(uri lsp.DocumentURI) (gtarget string, ok bool) {

	filename := uri.GetFileName()
	splits := strings.Split(filename, ".md")
	if len(splits) < 1 {
		return "", false
	}
	return splits[0], true
}
func getGTarget(heading string, uri lsp.DocumentURI) (gtarget GTarget, ok bool) {

	fileGtTarget, ok := GetFileGTarget(uri)
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

	return s.GLinkStore.AddDef(gtarget, location)
}

// returns ok
// adds full target
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

	return s.GLinkStore.AddDef(gtarget, location)
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

	return s.GLinkStore.RemoveDef(gtarget, location)
}

func (s *Store) GetWikiCompletions(arg string, needEnd bool, rng lsp.Range, uri *lsp.DocumentURI) []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}
	strppedArg := strings.TrimSpace(arg)
	isWithin := (len(arg) > 0 && arg[0] == '#') || (len(strppedArg) > 0 && strppedArg[0] == '#')
	if isWithin {
		doc, ok := s.GetDoc(*uri)
		if ok && doc.Tree != nil {
			headings := GetHeadings(&doc)
			for _, target := range headings {
				var link string
				if needEnd {
					link = "[[" + target + "]]"
				} else {
					link = "[[" + target
				}
				sortText := "a"
				var kind lsp.CompletionItemKind
				kind = lsp.EnumCompletion
				completions = append(completions, lsp.CompletionItem{
					Label:    link,
					Kind:     kind,
					SortText: sortText,
					TextEdit: &lsp.TextEdit{
						Range:   rng,
						NewText: link,
					},
					Detail: "", //syltodo excerpt
				})
			}
		}
	} else {
		pipeLoc := strings.IndexRune(arg, '|')
		containsPipe := pipeLoc != -1
		needConceal := strings.ContainsRune(arg, '#') && containsPipe
		if needConceal {
			arg = arg[:pipeLoc]
		}
		for _, t := range s.GLinkStore.GetTargets() {

			target := string(t.target)
			match := fuzzy.MatchFold(arg, target)
			if match == false {
				continue
			}

			var link string
			if needEnd {
				link = "[[" + target + "]]"
			} else {
				link = "[[" + target
			}
			var excerpt string
			if t.loc != nil {
				excerpt = s.GetExcerpt(*t.loc)
			}
			isFile := !strings.ContainsRune(link, '#')
			sortText := "c"
			kind := lsp.ReferenceCompletion
			if isFile {
				sortText = "b"
				kind = lsp.FileCompletion
			}
			completions = append(completions, lsp.CompletionItem{
				Label:    link,
				Kind:     kind,
				SortText: sortText,
				TextEdit: &lsp.TextEdit{
					Range:   rng,
					NewText: link,
				},
				Detail: excerpt,
			})
			if needConceal {
				start := strings.IndexRune(link, '#')
				end := strings.IndexRune(link, ']')
				if start >= 0 {
					var concealerText string
					if end == -1 {
						concealerText = link[start+1:]
					} else {

						concealerText = link[start+1 : end]
					}
					if needEnd {
						link = "[[" + target + "|" + concealerText + "]]"
					} else {
						link = "[[" + target + "|" + concealerText
					}
					completions = append(completions, lsp.CompletionItem{
						Label:    link,
						Kind:     kind,
						SortText: "a",
						// sylopti can use InsertReplaceEdit  instead of lsp.TextEdit
						TextEdit: &lsp.TextEdit{
							Range:   rng,
							NewText: link,
						},
						Detail: excerpt,
					})
				}
			}
		}
	}

	return completions
}

func (s *Store) GetGTargetDefinition(target GTarget) []lsp.Location {
	locs, _ := s.GLinkStore.GetDefs(target)
	return locs
}

func (s *Store) GetGTargetHeadingHover(target GTarget) string {
	var totalRefs int

	refs, _ := s.GLinkStore.GetRefs(target)
	totalRefs = len(refs)
	content := fmt.Sprintf("%d references found", totalRefs)
	return content
}

func (s *Store) GetGTargetWikilinkHover(target GTarget) string {
	content := ""
	refs, found := s.GLinkStore.GetRefs(target)
	if found {
		content = fmt.Sprintf("%d references found\n", len(refs)) + content
	}

	defs, found := s.GLinkStore.GetDefs(target)
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
	if len(refs) > 0 {
		// references md
		var refmd string
		for _, loc := range refs {
			refmd = fmt.Sprintf("%s1 %s\n", refmd, loc.URI.GetFileName())
		}
		content = content + "\n---\n" + refmd
	}

	return content
}

func (s *Store) GetGTargetReferences(target GTarget) []lsp.Location {
	refs, _ := s.GLinkStore.GetRefs(target)
	return refs
}
