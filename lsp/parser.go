package lsp

import (
	"fmt"
	"log/slog"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type ParseFunction func(content string, oldTree *tree_sitter.Tree) *tree_sitter.Tree

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

func GetNodeContent(node tree_sitter.Node, content string) string {
	return content[node.StartByte():node.EndByte()]
}

func PointFromPosition(pos Position) tree_sitter.Point {
	targetRow := uint(pos.Line)
	targetColumn := uint(pos.Character)
	return tree_sitter.Point{Row: targetRow, Column: targetColumn}
}

func GetRange(node *tree_sitter.Node) Range {
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

func FieldText(parent *tree_sitter.Node, field string, content string) (string, bool) {
	child := parent.ChildByFieldName(field)
	if child == nil {
		return "", false
	}
	return string(content[child.StartByte():child.EndByte()]), true
}
