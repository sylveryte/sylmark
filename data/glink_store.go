package data

import (
	"sylmark/lsp"
)

type glink struct {
	def  lsp.Location
	refs []lsp.Location
}

func newGlink() glink {
	return glink{
		def:  lsp.Location{},
		refs: []lsp.Location{},
	}
}

type GLinkStore map[GTarget]map[lsp.DocumentURI]glink

func newGlinkMap() map[lsp.DocumentURI]glink {
	return map[lsp.DocumentURI]glink{}
}
func NewGlinkStore() GLinkStore {
	return map[GTarget]map[lsp.DocumentURI]glink{}
}
func (glinkStore *GLinkStore) AddRef(target GTarget, loc lsp.Location) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glinkmap, found := gs[target]
	if found {
		glink, found := glinkmap[loc.URI]

		if found {
			glink.refs = append(glink.refs, loc)
		} else {
			glink := newGlink()
			glink.refs = append(glink.refs, loc)
			glinkmap[loc.URI] = glink
		}
	} else {
		glinkmap = newGlinkMap()
		glink := newGlink()
		glink.refs = append(glink.refs, loc)
		glinkmap[loc.URI] = glink
		gs[target] = glinkmap
	}

	return false
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

func (glinkStore *GLinkStore) AddDef(target GTarget, uri lsp.DocumentURI, rng lsp.Range) bool {
	if glinkStore == nil {
		return false
	}
	gs := *glinkStore

	glinkmap, found := gs[target]
	if !found {
		glinkmap = newGlinkMap()
		glink := newGlink()
		glink.def = lsp.Location{
			URI:   uri,
			Range: rng,
		}
		glinkmap[uri] = glink
		gs[target] = glinkmap
	} else {
		glink, found := glinkmap[uri]
		if !found {
			glink = newGlink()
			glinkmap[uri] = glink
		}
		glink.def = lsp.Location{
			URI:   uri,
			Range: rng,
		}
	}

	return false
}

func (glinkStore *GLinkStore) GetDefs(target GTarget) (locs []lsp.Location, gfound bool) {
	if glinkStore == nil {
		return locs, gfound
	}
	gs := *glinkStore

	glinkmap, found := gs[target]
	if !found {
		return locs, gfound
	}

	for k, v := range glinkmap {
		found = true
		if k != "" {
			locs = append(locs, v.def)
		}
	}

	return locs, gfound
}
