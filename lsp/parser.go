package lsp

import (
	"fmt"
	"log/slog"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Trees [2]*tree_sitter.Tree

func (t *Trees) GetMainTree() *tree_sitter.Tree {
	if t != nil {
		return t[0]
	}
	slog.Error("No maintree")
	return nil
}
func (t *Trees) GetInlineTree() *tree_sitter.Tree {
	if t != nil {
		return t[1]
	}
	slog.Error("No inlinetree")
	return nil
}

type ParseFunction func(content string, oldTrees *Trees) *Trees

func PrintTsTree(node tree_sitter.Node, depth int, cont string) {
	indent := ""
	for range depth {
		indent += "  "
	}

	nodeConte := cont[node.StartByte():node.EndByte()]
	slog.Info(fmt.Sprintf("%sNode: %s (%s), Range: (%d,%d)-(%d,%d) Text: (%s)", indent, node.Kind(), node.GrammarName(), node.StartByte(), node.StartPosition(), node.EndByte(), node.EndPosition(), nodeConte))

	for i := range int(node.NamedChildCount()) {
		child := node.NamedChild(uint(i))
		PrintTsTree(*child, depth+1, cont)
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

	for i := range int(node.NamedChildCount()) {
		child := node.NamedChild(uint(i))
		action(child)
		TraverseNodeWith(child, action)
	}

}

func IsInlineParseNeeded(node *tree_sitter.Node) bool {
	return node.Kind() == "inline" || node.Kind() == "paragraph"
}

func FieldText(parent *tree_sitter.Node, field string, content string) (string, bool) {
	child := parent.ChildByFieldName(field)
	if child == nil {
		return "", false
	}
	return string(content[child.StartByte():child.EndByte()]), true
}

func GetParentalKind(node *tree_sitter.Node) *tree_sitter.Node {
	k := node.Kind()
	switch k {
	case "link_destination", "link_text", "heading_content":
		return node.Parent()
	}
	return node
}
