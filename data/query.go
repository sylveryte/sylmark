package data

import (
	"log/slog"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

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

func (s *Store) GetSemanticTokens(id Id, parse lsp.ParseFunction) lsp.SemantiTokens {
	intTokens := []uint{}

	// get tokens and convert them to intTokens
	docData, found := s.GetDocMustTree(id, parse)
	if !found {
		slog.Error("Shocking doc not found for SemantiTokens")
		return lsp.SemantiTokens{
			Data: []uint{},
		}
	}

	// these token poistions are relative to last one hence lastLastLine
	var lastLine uint
	var lastStart uint
	lsp.TraverseNodeWith(docData.Trees.GetInlineTree().RootNode(), func(n *tree_sitter.Node) {
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

	return lsp.SemantiTokens{
		Data: intTokens,
	}
}
