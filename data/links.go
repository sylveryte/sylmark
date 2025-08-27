package data

import (
	"log/slog"
	"strings"
	"sylmark-server/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func getLinkUrl(node *tree_sitter.Node, content string) (url string, ok bool) {

	url, ok = lsp.FieldText(node, "url", content)

	if !ok {
		return url, ok
	}

	url = strings.TrimSpace(url)

	return url, ok
}

func getHeadingTitle(node *tree_sitter.Node, content string) (link string, ok bool) {

	link, ok = lsp.FieldText(node, "title", content)

	if !ok {
		return
	}

	link = strings.TrimSpace(link)

	return
}

// handles node kind wikilink, wiktarget, heading, title
func GetWikilinkTarget(node *tree_sitter.Node, content string, uri lsp.DocumentURI) (target GTarget, ok bool) {

	var link string
	if node.Kind() == "wikilink" {
		link, ok = lsp.FieldText(node, "target", content)
		if !ok {
			return
		}
	} else if node.Kind() == "heading" || node.Kind() == "title" {
		var heading string
		if node.Kind() == "heading" {
			heading, ok = getHeadingTitle(node, content)
			if !ok {
				slog.Error("Could not extract heading")
				return "", false
			}
		} else {
			heading = lsp.GetNodeContent(*node, content)
			heading = strings.TrimSpace(heading) // important
		}
		gtarget, ok := getGTarget(heading, uri)
		if !ok {
			slog.Error("Could not form gtarget")
			return "", false
		}
		link = string(gtarget)
	} else if node.Parent().Kind() == "wikilink" {
		link = lsp.GetNodeContent(*node, content)
	} else {
		return "", false
	}

	if strings.Contains(link, "|") {
		link = strings.Split(link, "|")[0]
	}

	target = GTarget(strings.TrimSpace(link))

	return target, true
}
