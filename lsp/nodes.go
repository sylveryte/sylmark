package lsp

import (
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func getTag(node *tree_sitter.Node, content string) Tag {

	t := getNodeContent(*node, Document(content))
	t = strings.TrimSpace(t)

	return Tag(t)
}

func getLinkUrl(node *tree_sitter.Node, content string) (url string, ok bool) {

	url, ok = fieldText(node, "url", content)

	if !ok {
		return
	}

	url = strings.TrimSpace(url)

	return
}

func getHeadingTitle(node *tree_sitter.Node, content string) (link string, ok bool) {

	link, ok = fieldText(node, "title", content)

	if !ok {
		return
	}

	link = strings.TrimSpace(link)

	return
}

func getWikilinkLink(node *tree_sitter.Node, content string) (link string, ok bool) {

	link, ok = fieldText(node, "target", content)

	if !ok {
		return
	}

	if strings.Contains(link, "|") {
		link = strings.Split(link, "|")[0]
	}

	link = strings.TrimSpace(link)

	return
}
