package data

import (
	"sylmark-server/lsp"
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

// returns ok
func (glinkStore *GLinkStore) RemoveRef(target GTarget, loc lsp.Location) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glink, found := gs[target]

	if !found {
		return true
	}
	var newRefs []lsp.Location
	for _, ref := range glink.refs {
		if ref.URI == loc.URI && ref.Range.Start == loc.Range.Start {
			continue
		}
		newRefs = append(newRefs, ref)
	}

	glink.refs = newRefs
	gs[target] = glink
	return true
}

func (glinkStore *GLinkStore) AddDef(target GTarget, loc lsp.Location) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glink, found := gs[target]

	if !found {
		glink = newGlink()
	}
	glink.defs = append(glink.defs, loc)
	gs[target] = glink

	return true
}

// returns ok
func (glinkStore *GLinkStore) RemoveDef(target GTarget, loc lsp.Location) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glink, found := gs[target]

	if !found {
		return true
	}
	var newDefs []lsp.Location
	for _, def := range glink.defs {
		if def.URI == loc.URI && def.Range.Start == loc.Range.Start {
			continue
		}
		newDefs = append(newDefs, def)
	}

	glink.defs = newDefs
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

func (glinkStore *GLinkStore) GetTargets() []GTargetAndLoc {
	targets := []GTargetAndLoc{}
	if glinkStore == nil {
		return targets
	}
	gs := *glinkStore

	for k, v := range gs {
		var def *lsp.Location
		if len(v.defs) > 0 {
			def = &v.defs[0]
		}
		targets = append(targets, GTargetAndLoc{
			target: k,
			loc:    def,
		})
	}

	return targets
}
