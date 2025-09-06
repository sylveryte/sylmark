package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"
	"sylmark/utils"
)

func (store *Store) GetCompletions(params lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	completions := []lsp.CompletionItem{}

	doc, found := store.GetDoc(params.TextDocument.URI)
	if !found {
		slog.Error("not found" + string(params.TextDocument.URI))
		return completions, fmt.Errorf("Not found")
	}

	line := doc.Content.GetLine(params.Position.Line)

	isTag := false
	isWiki := false
	isWikiEnd := true

	// TODO fix sylbug can we improve if say [[SOme space link starting from here#]]
	// this case where space is not considrered by below code
	before, after, found := utils.FindWord(params.Position.Character, line)
	before = strings.TrimSpace(before)
	if len(before) > 0 && before[0] == '#' {
		isTag = true
	} else if len(before) > 1 {
		if before[0] == '[' && before[1] == '[' {
			isWiki = true
		}
		if len(after) > 1 {
			if after[0] == ']' && after[1] == ']' {
				isWikiEnd = false
			}
		}
	}

	// Tags
	if isTag {
		tagCompletions := store.GetTagCompletions()
		completions = append(completions, tagCompletions...)
	} else if isWiki {
		// wiklink
		wikiCompletions := store.GetWikiCompletions(isWikiEnd, nil)
		completions = append(completions, wikiCompletions...)
	}

	return completions, nil
}
