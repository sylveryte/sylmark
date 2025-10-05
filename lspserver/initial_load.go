package lspserver

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"sylmark/data"
	"sylmark/lsp"
	"time"
)

func (h *LangHandler) addRootPathAndLoad(dir string) {
	h.Store.Config.RootPath = dir
	t := time.Now()
	h.loadAllClosedDocsData()
	slog.Info(fmt.Sprintf("=====>Initial Load time is [[%dms]]<=====", time.Since(t).Milliseconds()))
	h.Store.Config.CreatDirsIfNeeded()
}

func (h *LangHandler) loadAllClosedDocsData() {
	if h.Store.Config.RootPath == "" {
		slog.Error("h.rootPath is empty")
		return
	}

	parallels := 25 // 8 seems to give best results

	in := make(chan string, parallels)
	defer close(in)
	out := make(chan *TreesContent, parallels)
	defer close(out)

	// processing goroutines
	for range parallels {
		go func() {
			parsers := getParsers()
			parse := getParseFunction(parsers)
			for mdFilePath := range in {
				uri, content, trees, err := TreesFromMdDocPath(mdFilePath, parse)
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
			parsers[0].Close()
			parsers[1].Close()
		}()
	}

	// input prepare
	filepath.WalkDir(h.Store.Config.RootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && (strings.HasSuffix(path, ".git") || strings.HasSuffix(path, "node_modules")) {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			if strings.HasSuffix(path, ".md") {
				h.Store.MdFiles = append(h.Store.MdFiles, path)
			} else {
				h.Store.OtherFiles = append(h.Store.OtherFiles, path)
			}
		}
		return nil
	})

	// input goroutine
	go func() {
		for _, path := range h.Store.MdFiles {
			in <- path
		}
	}()

	// collect out goroutine
	total := len(h.Store.MdFiles)
	for val := range out {
		if val.ok {
			h.Store.LoadData(val.uri, val.content, val.trees)
			// clean up trees
			val.trees[0].Close()
			val.trees[1].Close()
		} else {
			slog.Error("Could not process " + string(val.uri.GetFileName()))
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

func TreesFromMdDocPath(mdDocPath string, parse lsp.ParseFunction) (uri lsp.DocumentURI, content string, trees *lsp.Trees, err error) {
	content = data.ContentFromDocPath(mdDocPath)
	uri, err = data.UriFromPath(mdDocPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	trees = parse(content, nil)
	return
}
