package lsp

import (
	"log/slog"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)


func (h *LangHandler) DocAndNodeFromURIAndPosition(uri DocumentURI, position Position) (doc Document, node *tree_sitter.Node, ok bool) {
	docData, ok := h.openedDocs.docDataFromURI(uri)
	if !ok {
		slog.Error("Document missing" + string(uri))
		return "", nil, false
	}
	point := pointFromPosition(position)

	doc = docData.Content
	node = docData.Tree.RootNode().NamedDescendantForPointRange(point, point)

	ok = true
	return
}
