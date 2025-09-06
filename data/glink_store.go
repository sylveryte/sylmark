package data

import (
	"sylmark/lsp"
)

type Glink struct {
	Defs []lsp.Location
	Refs []lsp.Location
}

func newGlink() Glink {
	return Glink{
		Defs: []lsp.Location{},
		Refs: []lsp.Location{},
	}
}

type GLinkStore map[GTarget]Glink

func NewGlinkStore() GLinkStore {
	return map[GTarget]Glink{}
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
	glink.Refs = append(glink.Refs, loc)
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
	for _, ref := range glink.Refs {
		if ref.URI == loc.URI && ref.Range.Start == loc.Range.Start {
			continue
		}
		newRefs = append(newRefs, ref)
	}

	if len(newRefs) == 0 && len(glink.Defs) == 0 {
		delete(gs, target)
		return true
	} else {
		glink.Refs = newRefs
		gs[target] = glink
		return true
	}
}

func (glinkStore *GLinkStore) AddDef(target GTarget, loc lsp.Location) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glink, found := gs[target]

	if !found {
		glink = newGlink()
	} else {
		// check if exists
		for _, k := range glink.Defs {
			if k.URI == loc.URI {
				return true
			}
		}
		// doesnt exists continue to add
	}
	glink.Defs = append(glink.Defs, loc)
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
	for _, def := range glink.Defs {
		if def.URI == loc.URI && def.Range.Start == loc.Range.Start {
			continue
		}
		newDefs = append(newDefs, def)
	}
	if len(newDefs) == 0 && len(glink.Refs) == 0 {
		delete(gs, target)
		return true
	} else {
		glink.Defs = newDefs
		gs[target] = glink
		return true
	}
}

func (glinkStore *GLinkStore) GetRefs(target GTarget) (refs []lsp.Location, found bool) {
	if glinkStore == nil {
		return refs, found
	}
	gs := *glinkStore

	glink, found := gs[target]

	return glink.Refs, found && len(glink.Refs) > 0
}

func (glinkStore *GLinkStore) GetDefs(target GTarget) (defs []lsp.Location, found bool) {
	if glinkStore == nil {
		return defs, found
	}
	gs := *glinkStore

	glink, found := gs[target]

	return glink.Defs, found && len(glink.Defs) > 0
}

func (glinkStore *GLinkStore) GetTargets() []GTargetAndLoc {
	targets := []GTargetAndLoc{}
	if glinkStore == nil {
		return targets
	}
	gs := *glinkStore

	for k, v := range gs {
		var def *lsp.Location
		if len(v.Defs) > 0 {
			def = &v.Defs[0]
		}
		targets = append(targets, GTargetAndLoc{
			target: k,
			loc:    def,
		})
	}

	return targets
}
