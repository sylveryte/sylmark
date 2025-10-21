package data

import (
	"fmt"
	"log/slog"
	"strings"
	"sylmark/lsp"

	"github.com/lithammer/fuzzysearch/fuzzy"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type FootNoteRef struct {
	Def    *lsp.Range
	Refs   []lsp.Range
	Excert string
}

type FootNotesStore map[string]FootNoteRef

func GetNewFootNotesStore() *FootNotesStore {
	return &FootNotesStore{}
}
func (fs *FootNotesStore) AddDefRef(target string, isDefinition bool, rng lsp.Range, excert string) {
	if fs == nil {
		slog.Error("FootNotesStore is nil")
		return
	}
	s := *fs
	fn, ok := s[target]
	if !ok {
		fn = FootNoteRef{}
	}
	if isDefinition {
		fn.Def = &rng
		fn.Excert = excert
	} else {
		// is a reference
		fn.Refs = append(fn.Refs, rng)
	}
	s[target] = fn
}

func (fs *FootNotesStore) GetFootNote(target string) (FootNoteRef, bool) {
	if fs == nil {
		slog.Error("FootNotesStore is nil")
	}
	s := *fs
	r, ok := s[target]
	return r, ok
}

func (s *Store) GetLoadedFootNotesStore(id Id, parse lsp.ParseFunction) *FootNotesStore {
	fs := GetNewFootNotesStore()
	docData, ok := s.GetDocMustTree(id, parse)
	if ok {
		lsp.TraverseNodeWith(docData.Trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
			switch n.Kind() {
			case "shortcut_link":
				{
					rng := lsp.GetRange(n)
					linkTextNode := n.NamedChild(0)
					linkText := lsp.GetNodeContent(*linkTextNode, string(docData.Content))
					line := docData.Content.GetLine(rng.Start.Line)
					excert, ok := getExecertOfShortcutLink(rng.End.Character-1, line)
					fs.AddDefRef(linkText, ok, rng, excert)
				}
			}
		})
	}

	return fs
}

func getExecertOfShortcutLink(endChar int, line string) (excert string, ok bool) {
	start := endChar + 2
	if start < len(line) && line[start-1] == ':' {
		return line[start : len(line)-1], true
	}
	return "", false
}

func (s *Store) GetFootNoteCompletions(arg string, rng lsp.Range, id Id) []lsp.CompletionItem {
	completions := []lsp.CompletionItem{}
	strppedArg := strings.TrimSpace(arg)
	doc, ok := s.GetDoc(id)
	if !ok {
		slog.Error("Failed to get doc for GetFootNoteCompletions")
		return completions
	}

	for target, footNote := range *doc.FootNotes {

		match := fuzzy.MatchFold(strppedArg, target)
		if match == false {
			continue
		}

		link := fmt.Sprintf("[%s]", target)

		var excert string
		if footNote.Def != nil {
			excert = footNote.Excert
		} else {
			excert = "_No definition found._"
		}

		completions = append(completions, lsp.CompletionItem{
			Label:    link,
			Kind:     lsp.VariableCompletion,
			SortText: "a",
			TextEdit: &lsp.TextEdit{
				Range:   rng,
				NewText: link,
			},
			Detail: excert,
		})
	}

	return completions
}
