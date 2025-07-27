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

func (h *LangHandler) parse(content string) *tree_sitter.Tree {
	return h.Parser.Parse([]byte(content), nil)
}

func printTsTree(node tree_sitter.Node, depth int, cont string) {
	indent := ""
	for range depth {
		indent += "  "
	}

	nodeConte := cont[node.StartByte():node.EndByte()]
	slog.Info(fmt.Sprintf("%sNode: %s (%s), Range: (%d,%d)-(%d,%d) Text: (%s)", indent, node.Kind(), node.GrammarName(), node.StartByte(), node.StartPosition(), node.EndByte(), node.EndPosition(), nodeConte))

	for i := range int(node.ChildCount()) {
		child := node.Child(uint(i))
		printTsTree(*child, depth+1, cont)
	}
}

func getNodeContent(node tree_sitter.Node, content Document) string {
	return string(content)[node.StartByte():node.EndByte()]
}

func pointFromPosition(pos Position) tree_sitter.Point {
	targetRow := uint(pos.Line)
	targetColumn := uint(pos.Character)
	return tree_sitter.Point{Row: targetRow, Column: targetColumn}
}

func getRange(node *tree_sitter.Node) Range {
	r := Range{}
	r.Start = Position{
		Line:      int(node.StartPosition().Row),
		Character: int(node.StartPosition().Column),
	}
	r.End = Position{
		Line:      int(node.EndPosition().Row),
		Character: int(node.EndPosition().Column),
	}

	return r
}

func TraverseNodeWith(node *tree_sitter.Node, action func(*tree_sitter.Node)) {

	for i := range int(node.ChildCount()) {
		child := node.Child(uint(i))
		action(child)
		TraverseNodeWith(child, action)
	}

}

func fieldText(parent *tree_sitter.Node, field string, content string) (string, bool) {
	child := parent.ChildByFieldName(field)
	if child == nil {
		return "", false
	}
	return string(content[child.StartByte():child.EndByte()]), true
}
