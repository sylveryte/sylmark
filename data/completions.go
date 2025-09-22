package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"
)

func (store *Store) GetCompletions(params lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	completions := []lsp.CompletionItem{}

	doc, found := store.GetDoc(params.TextDocument.URI)
	if !found {
		slog.Error("not found" + string(params.TextDocument.URI))
		return completions, fmt.Errorf("Not found")
	}

	line := doc.Content.GetLine(params.Position.Line)

	// this case where space is not considrered by below code
	kind, arg, cstart, cend := analyzeTriggerKind(params.Position.Character, line)
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
	arg = strings.TrimSpace(arg)

	// Tags
	switch kind {
	case CompletionTag:
		tagCompletions := store.GetTagCompletions(arg, rng)
		completions = append(completions, tagCompletions...)
	case CompletionWiki:
		// wiklink
		wikiCompletions := store.GetWikiCompletions(arg, true, rng, &params.TextDocument.URI)
		completions = append(completions, wikiCompletions...)
	case CompletionWikiWithEnd:
		// wiklink
		wikiCompletions := store.GetWikiCompletions(arg, false, rng, &params.TextDocument.URI)
		completions = append(completions, wikiCompletions...)
	}

	return completions, nil
}

type CompletionTriggerKind int8

const (
	CompletionNone        CompletionTriggerKind = 0
	CompletionTag         CompletionTriggerKind = 1
	CompletionWiki        CompletionTriggerKind = 2
	CompletionWikiWithEnd CompletionTriggerKind = 3
)

func analyzeTriggerKind(char int, line string) (kind CompletionTriggerKind, arg string, cstart, cend int) {

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
			if wikistart+1 < char {
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
		}

	}

	return
}
