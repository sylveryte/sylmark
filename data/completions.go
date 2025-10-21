package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"
)

func (s *Store) GetCompletions(params lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	completions := []lsp.CompletionItem{}

	id := s.GetIdFromURI(params.TextDocument.URI)
	doc, found := s.GetDoc(id)
	if !found {
		slog.Error("not found" + string(params.TextDocument.URI))
		return completions, fmt.Errorf("Not found")
	}

	line := doc.Content.GetLine(params.Position.Line)

	// utils.Sprintf("%d=%d Line=%s", params.Position.Character, params.CompletionContext.TriggerKind, line)

	// this case where space is not considrered by below code
	kind, arg, arg2, cstart, cend := analyzeTriggerKind(params.Position.Character, line)
	// utils.Sprintf("Trigge kind is %d", kind)
	rng := lsp.Range{
		Start: lsp.Position{
			Line:      params.Position.Line,
			Character: cstart,
		},
		End: lsp.Position{
			Line:      params.Position.Line,
			Character: cend,
		},
	}
	oneSpace := len(arg) > 0 && arg[0] == ' '
	// twoSpace := false
	// if oneSpace {
	// 	twoSpace = len(arg) > 1 && arg[1] == ' '
	// }
	arg = strings.TrimSpace(arg)

	// Tags
	switch kind {
	case CompletionTag:
		tagCompletions := s.GetTagCompletions(arg, rng)
		completions = append(completions, tagCompletions...)
	case CompletionWiki:
		// wiklink
		wikiCompletions := s.GetWikiCompletions(arg, true, oneSpace, rng, id)
		completions = append(completions, wikiCompletions...)
	case CompletionWikiWithEnd:
		// wiklink
		wikiCompletions := s.GetWikiCompletions(arg, false, oneSpace, rng, id)
		completions = append(completions, wikiCompletions...)
	case CompletionInlineLink, CompletionInlineLinkEnd:
		// inlinedownslinks syltodo cleanup args for below
		inlinedownLinkCompletions := s.GetInlineLinkCompletions(arg, "", rng, &params.TextDocument.URI)
		completions = append(completions, inlinedownLinkCompletions...)
	case CompletionInlineLinkHasText, CompletionInlineLinkEndHasText:
		// inlinedownslinks syltodo cleanup args for below
		inlinedownLinkCompletions := s.GetInlineLinkCompletions(arg, arg2, rng, &params.TextDocument.URI)
		completions = append(completions, inlinedownLinkCompletions...)
	case CompletionFootNote, CompletionFootNoteEnd:
		// inlinedownslinks syltodo cleanup args for below
		footNoteCompletions := s.GetFootNoteCompletions(arg, rng, id)
		completions = append(completions, footNoteCompletions...)
	}

	return completions, nil
}

type CompletionTriggerKind int8

const (
	CompletionNone                 CompletionTriggerKind = 0
	CompletionTag                  CompletionTriggerKind = 1
	CompletionWiki                 CompletionTriggerKind = 2
	CompletionWikiWithEnd          CompletionTriggerKind = 3
	CompletionInlineLink           CompletionTriggerKind = 4
	CompletionInlineLinkEnd        CompletionTriggerKind = 5
	CompletionInlineLinkHasText    CompletionTriggerKind = 6
	CompletionInlineLinkEndHasText CompletionTriggerKind = 7
	CompletionFootNote             CompletionTriggerKind = 8
	CompletionFootNoteEnd          CompletionTriggerKind = 9
)

func analyzeTriggerKind(char int, line string) (kind CompletionTriggerKind, arg string, arg2 string, cstart, cend int) {

	if len(line) > 0 && char > 0 && char <= len(line) {
		// note char is 1 indexed

		// look for # tag
		tag := -1
		for i := char - 1; ; {
			if i >= 0 && line[i] != ' ' {
				if line[i] == '#' {
					tag = i
					break
				}
			} else {
				break
			}
			i--
		}

		// loof for [
		footNote := -1
		for i := char - 1; i >= 0; i-- {
			ch := line[i]
			if ch == '[' {
				if !(i > 0 && line[i-1] == '[') {
					footNote = i
				}
				break
			}
		}

		// look for (
		inlinelink := -1
		inlinelinkurlstart := -1
		for i := char - 1; i > 0; i-- {
			ch := line[i]
			if ch == '(' {
				inlinelinkurlstart = i
				if line[i-1] == ']' {
					// find start
					for j := i - 2; j >= 0; j-- {
						if line[j] == '[' {
							inlinelink = j
							break
						} else if line[j] == ']' {
							break
						}
					}
				}
				break
			}
		}

		// look for [[, stop if ]] is seen
		wikistart := -1
		for i := char - 1; i > 0; i-- {
			ch := line[i]
			if ch == '[' {
				if line[i-1] == '[' {
					wikistart = i
					break
				}
			} else if line[i] == ']' {
				wikistart = -1
				break
			}
		}

		if tag > -1 && wikistart == -1 {
			// sylopti can make it so it takes right side into consideration for eg if cursor is at 2 for #sup can give arg=sup instead of arg=su
			// check right for better arg and better cend
			cstart = tag
			arg = line[tag+1 : char]
			cend = char
			kind = CompletionTag
		} else if wikistart > -1 {
			if wikistart+1 <= char {
				wikiend := char
				wikiendWord := char
				for wikiend < len(line) {
					if line[wikiend] == ' ' && wikiendWord == char {
						wikiendWord = wikiend
					}
					if line[wikiend] == ']' {
						kind = CompletionWikiWithEnd
						break
					}
					wikiend++
				}
				if kind == CompletionWikiWithEnd {
					arg = line[wikistart+1 : wikiend]
					cend = wikiend
				} else if wikiendWord != char {
					arg = line[wikistart+1 : wikiendWord]
				} else {
					arg = line[wikistart+1 : char]
				}
			}
			cstart = wikistart - 1
			if kind != CompletionWikiWithEnd {
				cend = char
				kind = CompletionWiki
			}
		} else if inlinelink > -1 {
			// find end
			cstart = inlinelink
			cstartArgStart := 1
			// check if is image then consider ! as well
			if cstart > 0 && line[cstart-1] == '!' {
				cstart -= 1
				cstartArgStart = 2
			}
			cend = char
			kind = CompletionInlineLink

			// look for closing )
			for i := char; i < len(line); i++ {
				if line[i] == ')' {
					cend = i + 1
					kind = CompletionInlineLinkEnd
				} else if line[i] == '[' || line[i] == ']' {
					break
				}
			}
			arg2 = line[cstart+cstartArgStart : inlinelinkurlstart-1]
			if len(strings.TrimSpace(arg2)) != 0 {
				// it's not empty
				switch kind {
				case CompletionInlineLink:
					kind = CompletionInlineLinkHasText
				case CompletionInlineLinkEnd:
					kind = CompletionInlineLinkEndHasText
				}
			}
			arg = line[inlinelinkurlstart+1 : char]
		} else if footNote > -1 {
			cstart = footNote
			cend = char
			kind = CompletionFootNote
			arg = line[footNote+1 : char]
			for i := char; i < len(line); i++ {
				if line[i] == ']' {
					cend = i + 1
					kind = CompletionFootNoteEnd
				} else if line[i] == '[' || line[i] == '(' || line[i] == ')' {
					break
				}
			}

		}

	}

	return
}
