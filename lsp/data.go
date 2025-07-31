package lsp

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

type Target string // is like # Some heading
type GTarget string   // is full GLink target
type Tag string

func (h *LangHandler) addRootPathAndLoad(dir string) {
	h.rootPath = dir
	h.loadAllClosedDocsData()
}

func (h *LangHandler) loadAllClosedDocsData() {
	if h.rootPath == "" {
		slog.Error("h.rootPath is empty")
		return
	}

	filepath.WalkDir(h.rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.HasSuffix(path, ".git") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			h.loadDocData(path)
		}
		return nil
	})
}

func (h *LangHandler) loadDocData(mdDocPath string) {
	content := contentFromDocPath(mdDocPath)
	uri, err := uriFromPath(mdDocPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	tree := h.parse(content)
	defer tree.Close()
	h.store.loadData(uri, content, tree.RootNode())
}

func (h *LangHandler) onDocOpened(uri DocumentURI, content string) {
	tree := h.parse(content)
	h.openedDocs.addDoc(uri, Document(content), tree)
	doc := Document(content)

	h.openedDocs.addDoc(uri, doc, tree)
}
func (h *LangHandler) onDocClosed(uri DocumentURI) {
	// remove data into openedDocs
	_, found := h.openedDocs.removeDoc(uri)
	if !found {
		slog.Error("Document not in openedDocs")
		return
	}
}

func (h *LangHandler) onDocChanged(uri DocumentURI, changes TextDocumentContentChangeEvent) {

	// update data into openedDocs
	updatedDocData, oldDocData, ok := h.openedDocs.updateDoc(uri, changes, h)
	if !ok {
		slog.Info("Update doc failed.")
		return
	}

	// update openedDocsStore
	tempStoreOld := newStore()
	tempStoreOld.loadData(uri, string(oldDocData.Content), oldDocData.Tree.RootNode())

	tempStoreNew := newStore()
	tempStoreNew.loadData(uri, string(updatedDocData.Content), updatedDocData.Tree.RootNode())

	// syltodo TODO optimze this flow it
	// do the deltas
	h.store.subtractStore(&tempStoreOld)
	h.store.mergeStore(&tempStoreNew)

}
