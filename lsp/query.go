package lsp

func (h *LangHandler) getTagRefs(tag Tag) int {
	locs, ok := h.inactiveStore.Tags[tag]
	if ok {
		return len(locs)
	}

	return 0
}
