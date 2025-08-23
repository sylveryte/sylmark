package data

import (
	"fmt"
	"iter"
	"log/slog"
	"maps"
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

// returns ok
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
			return true
		} else {
			glink := newGlink()
			glink.refs = append(glink.refs, loc)
			glinkmap[loc.URI] = glink
			return true
		}
	} else {
		glinkmap = newGlinkMap()
		glink := newGlink()
		glink.refs = append(glink.refs, loc)
		glinkmap[loc.URI] = glink
		gs[target] = glinkmap
		return true
	}
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

	slog.Info(fmt.Sprintf("addingdef GTarget=[%s] and uri=[%s]", target, uri))

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
			glink.def = lsp.Location{
				URI:   uri,
				Range: rng,
			}
			glinkmap[uri] = glink
		} else {
			glink.def = lsp.Location{
				URI:   uri,
				Range: rng,
			}
		}
	}

	// syldoto delete this block
	nglink, nfound := glinkmap[uri]
	if nfound {
		slog.Info(fmt.Sprintf("newlyadded glink.uri=[%s] and uri=[%s]", nglink.def.URI, uri))
	}

	return true
}

func (glinkStore *GLinkStore) GetGLinks(target GTarget) (glinks iter.Seq[glink], count int, gfound bool) {
	if glinkStore == nil {
		return glinks, 0, gfound
	}
	gs := *glinkStore

	glinkmap, gfound := gs[target]
	if !gfound {
		return glinks, 0, gfound
	}

	glinks = maps.Values(glinkmap)
	linksCount := len(glinkmap)

	return glinks, linksCount, gfound
}
func (glinkStore *GLinkStore) GetRefs(target GTarget) (locs []lsp.Location, gfound bool) {
	if glinkStore == nil {
		return locs, gfound
	}
	gs := *glinkStore

	glinkmap, gfound := gs[target]
	if !gfound {
		return locs, gfound
	}
	for k, v := range glinkmap {
		if k != "" {
			locs = append(locs, v.refs...)
		}
	}

	return locs, len(locs) > 0
}

func (glinkStore *GLinkStore) GetDefs(target GTarget) (locs []lsp.Location, gfound bool) {
	if glinkStore == nil {
		return locs, gfound
	}
	gs := *glinkStore

	glinkmap, gfound := gs[target]
	if !gfound {
		return locs, gfound
	}

	for k, v := range glinkmap {
		gfound = true
		if k != "" && v.def.URI != "" {
			locs = append(locs, v.def)
		}
	}

	return locs, gfound
}
