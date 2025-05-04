package lsp

import (
	"fmt"
	"log/slog"
	"unsafe"

	"github.com/ebitengine/purego"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (h *LangHandler) SetupGrammars() {
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(loadGrammar())
	parser.SetLanguage(language)

	h.MarkdownParser = parser

	parserInline := tree_sitter.NewParser()
	languageInline := tree_sitter.NewLanguage(loadInlineGrammar())
	parserInline.SetLanguage(languageInline)
	h.InlineMarkdownParser = parserInline
	slog.Info("Grammars are set")
}

func (h *LangHandler) parseTreesitter(content string) {

	tree := h.MarkdownParser.Parse([]byte(content), nil)
	rootNode := tree.RootNode()
	slog.Info(fmt.Sprintf("n child %d ", rootNode.ChildCount()))
	h.printTsTree(*rootNode, 0, content)

	slog.Info("-------------------------------inline------------------------------------------------")

	itree := h.InlineMarkdownParser.Parse([]byte(content), nil)
	irootNode := itree.RootNode()
	slog.Info(fmt.Sprintf("n child %d ", irootNode.ChildCount()))
	h.printTsTree(*irootNode, 0, content)


}

func (h *LangHandler) printTsTree(node tree_sitter.Node, depth int, cont string) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	nodeConte := cont[node.StartByte():node.EndByte()]
	slog.Info(fmt.Sprintf("%sNode: %s (%s), Range: (%d,%d)-(%d,%d) Text: (%s)", indent, node.Kind(), node.GrammarName(), node.StartByte(), node.StartPosition(), node.EndByte(), node.EndPosition(), nodeConte))

	// if node.Kind() == "inline" {
	// 	slog.Info("Diving into inline instead of printing children")
	//
	// 	nodeTree := h.InlineMarkdownParser.Parse([]byte(nodeConte), nil)
	// 	h.printTsTree(*nodeTree.RootNode(),0,nodeConte)
	// } else {
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(uint(i))
			h.printTsTree(*child, depth+1, cont)
		}
	// }
}

// func parseGoldmark(content string) ast.Node {
//
// 	md := goldmark.New()
// 	r := text.NewReader([]byte(content))
// 	doc := md.Parser().Parse(r)
//
// 	slog.Info(fmt.Sprintf("ChildCount %d", doc.ChildCount()))
//
// 	slog.Info(fmt.Sprintf("Type %d, Kind %s", doc.Type(), doc.Kind()))
// 	n := doc.FirstChild()
// 	slog.Info("n has ", n.FirstChild().Kind())
// 	for {
// 		if n != nil {
// 			slog.Info(fmt.Sprintf("Type %d, Kind %s", n.Type(), n.Kind()))
// 		} else {
// 			break
// 		}
// 		n = n.NextSibling()
// 	}
// 	return doc
// }

func loadGrammar() unsafe.Pointer {
	path := "/home/sylveryte/projects/sylmark/markdown.so"
	lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		// handle error
		slog.Info("panic :" + err.Error())
		panic(err)
	}

	var k func() uintptr
	purego.RegisterLibFunc(&k, lib, "tree_sitter_markdown")

	lang := unsafe.Pointer(k())

	return lang
}

func loadInlineGrammar() unsafe.Pointer {
	path := "/home/sylveryte/projects/sylmark/markdown-inline.so"
	lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		// handle error
		slog.Info("panic :" + err.Error())
		panic(err)
	}

	var k func() uintptr
	purego.RegisterLibFunc(&k, lib, "tree_sitter_markdown_inline")

	lang := unsafe.Pointer(k())

	return lang
}
