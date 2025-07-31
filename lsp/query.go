package lsp

import (
	"log/slog"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (h *LangHandler) getTagRefs(tag Tag) int {

	clocs, found := h.store.Tags[tag]
	if found {
		return len(clocs)
	}

	return 0
}

func (h *LangHandler) getSemanticTokens(uri DocumentURI) SemantiTokens {
	intTokens := []uint{}

	// get tokens and convert them to intTokens
	docData, found := h.openedDocs[uri]
	if !found {
		slog.Info("Shocking doc not found for SemantiTokens" + string(uri))
		return SemantiTokens{
			Data: []uint{},
		}
	}

	// these token poistions are relative to last one hence lastLastLine
	var lastLine uint
	var lastStart uint
	TraverseNodeWith(docData.Tree.RootNode(), func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "tag":
			{
				token := getSemanticToken(n, 0, 0)
				lastLastLine := lastLine
				lastLastStart := lastStart // column in case two are on same row columns relative needed
				lastLine = token[0]
				lastStart = token[1]
				token[0] = token[0] - lastLastLine
				if lastLastLine == lastLine {
					token[1] = token[1] - lastLastStart
				}
				intTokens = append(intTokens, token...)
			}
		}
	})

	return SemantiTokens{
		Data: intTokens,
	}
}

func getSemanticToken(node *tree_sitter.Node, tokenTypeIndex uint, tokenModifierIndex uint) []uint {

	return []uint{
		node.StartPosition().Row,    //line
		node.StartPosition().Column, //char
		// node.EndPosition().Column - node.StartPosition().Column, //length
		node.EndByte() - node.StartByte(), //length
		tokenTypeIndex,
		tokenModifierIndex,
	}
}
