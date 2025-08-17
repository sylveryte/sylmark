package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"
	"sylmark/utils"
)

func (store *Store) GetCompletions(params lsp.CompletionParams, openedDocs DocumentStore) ([]lsp.CompletionItem, error) {
	completions := []lsp.CompletionItem{}

	doc, found := openedDocs.DocDataFromURI(params.TextDocument.URI)
	if !found {
		slog.Info("not found" + string(params.TextDocument.URI))
		return completions, fmt.Errorf("Not found")
	}

	line := doc.Content.GetLine(params.Position.Line)
	slog.Info(fmt.Sprintf("%v==>", params.Position.Character) + "mila re==[" + line + "]")
	// slog.Info("mila re==[" + line+"] ["+line[params.Position.Character:]+"]")

	isTag := false
	isWiki := false
	isWikiEnd := true

	before, after, found := utils.FindWord(params.Position.Character, line)
	before = strings.TrimSpace(before)
	// slog.Info(fmt.Sprintf("B-[%s] A-[%s] f-[%t]", before, after, found))
	if before[0] == '#' {
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
	// slog.Info(fmt.Sprintf("isTag=%t isWiki=%t isWikiEnd=%t", isTag, isWiki, isWikiEnd))

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
