package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type SubTargetAndRanges struct {
	subTarget SubTarget
	rng       *lsp.Range
}
type FullTargetAndLoc struct {
	FullTarget FullTarget
	rng        *lsp.Range
}

func GetTarget(uri lsp.DocumentURI) (target Target, ok bool) {

	filename := uri.GetFileName()
	splits := strings.Split(filename, ".md")
	if len(splits) < 1 {
		slog.Error("File not md " + string(uri))
		return "", false
	}
	return Target(splits[0]), true
}

// target path is relative to workspace
func (s *Store) GetVaultTarget(uri lsp.DocumentURI) (target Target, ok bool) {
	relPath, err := s.GetPathRelRoot(uri)
	if err != nil {
		slog.Error("Failed to get path from " + err.Error())
		return "", false
	}
	// getting rid of .md
	target = Target(relPath[:len(relPath)-3])
	return target, true
}

// sucess upTarget or same => ok
func GetOneUpTarget(vaultTarget Target) (target Target, ok bool) {
	splits := strings.Split(string(vaultTarget), "/")
	sl := len(splits)
	if sl > 1 {
		return Target(strings.Join([]string{splits[sl-2], splits[sl-1]}, "/")), true
	}
	return vaultTarget, false
}
func GetPlainTarget(vaultTarget Target) (target Target, ok bool) {
	splits := strings.Split(string(vaultTarget), "/")
	sl := len(splits)
	if sl > 0 {
		return Target(splits[sl-1]), true
	}
	return vaultTarget, false
}

func (s *Store) GetWikiCompletions(arg string, needEnd bool, onlyFiles bool, rng lsp.Range, id Id) []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}
	strppedArg := strings.TrimSpace(arg)
	argContainsHash := strings.ContainsRune(arg, '#')
	isWithin := (len(arg) > 0 && arg[0] == '#') || (len(strppedArg) > 0 && strppedArg[0] == '#')
	if isWithin {
		doc, ok := s.GetDoc(id)
		if ok && doc.Trees != nil {
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
		needConceal := argContainsHash && containsPipe
		if needConceal {
			arg = arg[:pipeLoc]
		}

		for target, ids := range s.TargetStore {
			match := fuzzy.MatchFold(arg, string(target))
			if match {
				// add file
				var link string
				if needEnd {
					link = fmt.Sprintf("[[%s]]", target)
				} else {
					link = fmt.Sprintf("[[%s", target)
				}
				completions = append(completions, lsp.CompletionItem{
					Label:    string(target),
					Kind:     lsp.FileCompletion,
					SortText: "b",
					TextEdit: &lsp.TextEdit{
						Range:   rng,
						NewText: link,
					},
					Detail: s.GetExcerpt(ids[0], lsp.Range{}),
				})
			}
			if onlyFiles {
				continue
			}
			// add FullTarget
			for _, id := range ids {
				subTargets := s.LinkStore.GetSubTargetsAndRanges(id)
				for _, subTargetNRange := range subTargets {
					fullTarget := FullTarget(string(target) + string(subTargetNRange.subTarget))
					match = fuzzy.MatchFold(arg, string(fullTarget))
					if match {
						var link string
						if needConceal {
							fullTarget = FullTarget(fmt.Sprintf("%s%s|%s", target, subTargetNRange.subTarget, subTargetNRange.subTarget))
						}
						if needEnd {
							link = fmt.Sprintf("[[%s]]", fullTarget)
						} else {
							link = fmt.Sprintf("[[%s", fullTarget)
						}
						var defRange lsp.Range
						if subTargetNRange.rng != nil {
							defRange = *subTargetNRange.rng
						}
						completions = append(completions, lsp.CompletionItem{
							Label:    string(fullTarget),
							Kind:     lsp.ReferenceCompletion,
							SortText: "c",
							TextEdit: &lsp.TextEdit{
								Range:   rng,
								NewText: link,
							},
							Detail: s.GetExcerpt(id, defRange),
						})
					}
				}
			}
		}
	}

	if !argContainsHash {
		dateCompletions := s.getDateCompletions(arg, needEnd, rng)
		completions = append(completions, dateCompletions...)
	}

	return completions
}
