package lspserver

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"sylmark/data"
	"sylmark/lsp"
)

func (h *LangHandler) addRootPathAndLoad(dir string) {
	h.loadAllClosedDocsData()
	h.Store.Config.CreatDirsIfNeeded()
}

func (h *LangHandler) loadAllClosedDocsData() {
	if h.Store.Config.RootPath == "" {
		slog.Error("h.rootPath is empty")
		return
	}

	parallels := 2500

	in := make(chan string, parallels)
	defer close(in)
	out := make(chan *TreesContent, parallels)
	defer close(out)

	// processing goroutines
	for range parallels {
		go func() {
			parsers := getParsers()
			defer parsers[0].Close()
			defer parsers[1].Close()
			parse := getParseFunction(parsers)
			for mdFilePath := range in {
				uri, content, trees, err := TreesFromUri(mdFilePath, parse)
				if err != nil {
					slog.Error("Some error while parsing " + err.Error())
					out <- &TreesContent{
						uri:     uri,
						content: content,
						trees:   nil,
						ok:      false,
					}
					continue
				}
				out <- &TreesContent{
					uri:     uri,
					content: content,
					trees:   trees,
					ok:      true,
				}
			}
		}()
	}

	var mdFiles []string

	// input prepare
	filepath.WalkDir(h.Store.Config.RootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && (strings.HasSuffix(path, ".") || strings.HasSuffix(path, "node_modules")) {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			if data.IsMdFile(path) {
				mdFiles = append(mdFiles, path)
			} else {
				h.Store.OtherFiles = append(h.Store.OtherFiles, path)
			}
		}
		return nil
	})

	// input goroutine
	go func() {
		for _, path := range mdFiles {
			in <- path
		}
	}()

	// collect out goroutine
	total := len(mdFiles)
	for val := range out {
		if val.ok {
			id := h.Store.GetIdFromURI(val.uri)
			// utils.Sprintf("jiko id is %d uri was %s", id, val.uri)
			h.Store.LoadData(id, val.content, val.trees)
			// clean up trees
			val.trees[0].Close()
			val.trees[1].Close()
		} else {
			slog.Error("Could not process " + string(val.uri))
		}
		total -= 1
		if total == 0 {
			break
		}
	}
}

type TreesContent struct {
	ok      bool
	uri     lsp.DocumentURI
	content string
	trees   *lsp.Trees
}

func TreesFromUri(mdDocPath string, parse lsp.ParseFunction) (uri lsp.DocumentURI, content string, trees *lsp.Trees, err error) {
	content = data.ContentFromDocPath(mdDocPath)
	uri, err = data.UriFromPath(mdDocPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	trees = parse(content, nil)
	return
}
