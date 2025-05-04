package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sourcegraph/jsonrpc2"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func (h *langHandler) handleTextDocumentDidOpen(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	content := params.TextDocument.Text

	slog.Info("Got uri " + string(params.TextDocument.URI))
	slog.Info("Got text " + string(params.TextDocument.Text))

	parser := tree_sitter.NewParser()
	defer parser.Close()

	md := goldmark.New()
	r := text.NewReader([]byte(content))
	doc := md.Parser().Parse(r)

	slog.Info(fmt.Sprintf("ChildCount %d", doc.ChildCount()))

	slog.Info(fmt.Sprintf("Type %d, Kind %s", doc.Type(), doc.Kind()))
	n:= doc.FirstChild()
	for {
		if n != nil {
			slog.Info(fmt.Sprintf("Type %d, Kind %s", n.Type(), n.Kind()))
		} else {
			break
		}
		n = n.NextSibling()
	}

	return nil, nil

}
