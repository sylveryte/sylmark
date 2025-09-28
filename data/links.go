package data

import (
	"log/slog"
	"strings"
	"sylmark/lsp"

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

func getHeadingTarget(node *tree_sitter.Node, content string) (link string, ok bool) {
	link, ok = getHeadingTitle(node, content)
	if ok {
		link = "#" + link
	}
	return link, ok
}

// this is to overcome parser issue of recognising `conten # not heading` as heading
func isValidHeading(node *tree_sitter.Node, content string) bool {
	// check if starting from 0
	startPoint := node.StartPosition()
	if startPoint.Column == 0 {
		return true
	}

	// since # should start within first 4 columns
	if startPoint.Column > 3 {
		return false
	}
	// check if from 0 to node start column is blank string
	line, ok := GetLineFromContent(content, int(startPoint.Row))
	if !ok {
	}

	trimmed := strings.TrimSpace(line[0:startPoint.Column])
	return len(trimmed) == 0
}

func getHeadingTitle(node *tree_sitter.Node, content string) (link string, ok bool) {
	ok = isValidHeading(node, content)
	if !ok {
		return
	}

	link, ok = lsp.FieldText(node, "heading_content", content)

	if !ok {
		return
	}

	link = strings.TrimSpace(link)

	return
}

// handles node kind wikilink, wiktarget, heading, title
func GetWikilinkTarget(node *tree_sitter.Node, content string, uri lsp.DocumentURI) (target GTarget, ok bool) {

	var link string
	if node.Kind() == "link_destination" {
		link = lsp.GetNodeContent(*node, content)
	} else if node.Kind() == "wiki_link" {
		ch := node.Child(2)
		link = lsp.GetNodeContent(*ch, content)
	} else if node.Kind() == "atx_heading" || node.Kind() == "heading_content" {
		var heading string
		if node.Kind() == "atx_heading" {
			heading, ok = getHeadingTitle(node, content)
			if !ok {
				slog.Error("GetWikilinkTarget Could not extract heading")
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
	} else {
		return "", false
	}

	// syltodo delete
	// if strings.Contains(link, "|") {
	// 	link = strings.Split(link, "|")[0]
	// }
	//

	target = GTarget(strings.TrimSpace(link))

	return target, true
}
