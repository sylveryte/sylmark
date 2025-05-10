package lsp

import (
	"fmt"
	"log/slog"

	tree_sitter_sylmark "github.com/sylveryte/tree-sitter-sylmark/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (h *LangHandler) SetupGrammars() {
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_sylmark.Language())
	parser.SetLanguage(language)

	h.Parser = parser

	slog.Info("Grammars are set")
}

func (h *LangHandler) parseTreesitter(content string) {

	tree := h.Parser.Parse([]byte(content), nil)
	rootNode := tree.RootNode()
	slog.Info(fmt.Sprintf("n child %d ", rootNode.ChildCount()))
	h.printTsTree(*rootNode, 0, content)

}

func (h *LangHandler) printTsTree(node tree_sitter.Node, depth int, cont string) {
	indent := ""
	for range depth {
		indent += "  "
	}

	nodeConte := cont[node.StartByte():node.EndByte()]
	slog.Info(fmt.Sprintf("%sNode: %s (%s), Range: (%d,%d)-(%d,%d) Text: (%s)", indent, node.Kind(), node.GrammarName(), node.StartByte(), node.StartPosition(), node.EndByte(), node.EndPosition(), nodeConte))

	for i := range int(node.ChildCount()) {
		child := node.Child(uint(i))
		h.printTsTree(*child, depth+1, cont)
	}
}
