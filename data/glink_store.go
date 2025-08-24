package data

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"
)

type glink struct {
	defs []lsp.Location
	refs []lsp.Location
}

func newGlink() glink {
	return glink{
		defs: []lsp.Location{},
		refs: []lsp.Location{},
	}
}

type GLinkStore map[GTarget]glink

func NewGlinkStore() GLinkStore {
	return map[GTarget]glink{}
}

// returns ok
func (glinkStore *GLinkStore) AddRef(target GTarget, loc lsp.Location) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glink, found := gs[target]

	if !found {
		glink = newGlink()
	}
	glink.refs = append(glink.refs, loc)
	gs[target] = glink
	return true
}


func (glinkStore *GLinkStore) AddDef(target GTarget, uri lsp.DocumentURI, rng lsp.Range) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	slog.Info(fmt.Sprintf("addingdef GTarget=[%s] and uri=[%s]", target, uri))

	glink, found := gs[target]

	if !found {
		glink = newGlink()
	}
	glink.defs = append(glink.defs,
		lsp.Location{
			URI:   uri,
			Range: rng,
		})
	gs[target] = glink

	return true
}

func (glinkStore *GLinkStore) GetRefs(target GTarget) (refs []lsp.Location, found bool) {
	if glinkStore == nil {
		return refs, found
	}
	gs := *glinkStore

	glink, found := gs[target]

	return glink.refs, found && len(glink.refs) > 0
}

func (glinkStore *GLinkStore) GetDefs(target GTarget) (defs []lsp.Location, found bool) {
	if glinkStore == nil {
		return defs, found
	}
	gs := *glinkStore

	glink, found := gs[target]

	return glink.defs, found && len(glink.defs) > 0
}

func (glinkStore *GLinkStore) GetTargets() []GTarget {
	targets := []GTarget{}
	if glinkStore == nil {
		return targets
	}
	gs := *glinkStore

	for k := range gs {
		targets = append(targets, k)
	}

	return targets
}
