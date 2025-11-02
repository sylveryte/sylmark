package data

import (
	"fmt"
	"log/slog"
	"sylmark/lsp"
)

type Link struct {
	Def  map[SubTarget]lsp.Range
	Refs map[SubTarget][]IdLocation
}

func newLink() Link {
	return Link{
		Def:  map[SubTarget]lsp.Range{},
		Refs: map[SubTarget][]IdLocation{},
	}
}

type LinkStore map[Id]Link

func NewlinkStore() LinkStore {
	return LinkStore{}
}

func (linkStore *LinkStore) Print() {
	slog.Info("LinkStore===------------------------------------------------------------")
	for k, v := range *linkStore {
		slog.Info(fmt.Sprintf("\n====[%d]>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", k))
		slog.Info("Refs>>>>>>>>>>>>>")
		for k, j := range v.Refs {
			slog.Info(fmt.Sprintf("\n    [%s]=%d", k, j))
		}
		slog.Info("Defs>>>>>>>>>>>>>")
		for k, j := range v.Def {
			slog.Info(fmt.Sprintf("\n    [%s]=%d", k, j))
		}
		slog.Info(fmt.Sprintf("\n====[%d]^^^^^^^^^^^^^^^^^^^^^^^^^^^^", k))
	}
	slog.Info("LinkStore===END------------------------------------------------------------")
}

// returns ok
func (linkStore *LinkStore) AddRef(id Id, subTarget SubTarget, loc IdLocation) bool {
	if linkStore == nil {
		return false
	}
	ls := *linkStore

	link, found := ls[id]

	if !found {
		link = newLink()
	}
	locs, found := link.Refs[subTarget]
	locs = append(locs, loc)
	link.Refs[subTarget] = locs

	ls[id] = link
	return true
}

func (linkStore *LinkStore) RemoveRef(id Id, subTarget SubTarget, loc IdLocation) bool {
	if linkStore == nil {
		return false
	}
	ls := *linkStore

	l, found := ls[id]

	if !found {
		return true
	}
	refs, found := l.Refs[subTarget]
	// utils.Sprintf("Removing ref id=%d %s refs=%d", id, subTarget, len(refs))
	if len(refs) > 0 {
		var newRefs []IdLocation
		for _, ref := range refs {
			// utils.Sprintf("%d Ref L=%d C=%d", ref.Id, ref.Range.Start.Line, ref.Range.Start.Character)
			// utils.Sprintf("%d loc L=%d C=%d", loc.Id, loc.Range.Start.Line, loc.Range.Start.Character)
			if ref.Id == loc.Id && ref.Range.Start.Line == loc.Range.Start.Line && ref.Range.Start.Character == loc.Range.Start.Character {
				// utils.Sprintf("%d bitesssssssssssssssssssssssssss dust", id)
				continue
			}
			newRefs = append(newRefs, ref)
			// utils.Sprintf("%d nooooooooooooooooooooooooooooooo bitesssssssssssssssssssssssssss dust", id)
		}

		l.Refs[subTarget] = newRefs
	}

	if len(l.Refs) == 0 && len(l.Def) == 0 {
		delete(ls, id)
		return true
	}
	ls[id] = l
	return true
}

func (linkStore *LinkStore) AddDef(id Id, subTarget SubTarget, rng lsp.Range) bool {
	if linkStore == nil {
		return false
	}
	ls := *linkStore

	l, found := ls[id]
	if !found {
		l = newLink()
	}
	l.Def[subTarget] = rng

	ls[id] = l
	return true
}

// returns ok
func (linkStore *LinkStore) RemoveDef(id Id, subTarget SubTarget, rng lsp.Range) bool {
	if linkStore == nil {
		return false
	}
	ls := *linkStore

	l, found := ls[id]

	if !found {
		return true
	}
	delete(l.Def, subTarget)

	if l.Def == nil && len(l.Refs) == 0 {
		delete(ls, id)
		return true
	} else {
		ls[id] = l
		return true
	}
}

func (linkStore *LinkStore) GetRefs(id Id, subTarget SubTarget) (locs []IdLocation, found bool) {
	if linkStore == nil {
		return locs, found
	}
	ls := *linkStore

	l, found := ls[id]
	if !found {
		return locs, false
	}

	if subTarget == "" {
		// get all
		for _, v := range l.Refs {
			locs = append(locs, v...)
		}
	} else {
		locs, found = l.Refs[subTarget]
	}
	return locs, found
}
func (linkStore *LinkStore) GetDef(id Id, subTarget SubTarget) (def lsp.Range, found bool) {
	if linkStore == nil {
		return
	}
	ls := *linkStore

	l, found := ls[id]
	if !found {
		return
	}
	def, found = l.Def[subTarget]
	return
}

func (linkStore *LinkStore) GetSubTargetHover(id Id, subTarget SubTarget) string {
	refs, _ := linkStore.GetRefs(id, subTarget)
	totalRefs := len(refs)
	return fmt.Sprintf("%d references found", totalRefs)
}

func (linkStore *LinkStore) GetSubTargets(id Id) (subTargets []SubTarget) {
	if linkStore == nil {
		return
	}
	ls := *linkStore
	l, found := ls[id]
	if !found {
		return
	}

	sm := map[SubTarget]bool{}
	for subTarget := range l.Def {
		sm[subTarget] = true
	}
	for subTarget := range l.Refs {
		sm[subTarget] = true
	}
	for subTarget := range sm {
		subTargets = append(subTargets, subTarget)
	}
	return subTargets
}
func (linkStore *LinkStore) getSubTargetsAndRanges(id Id) (subTargets []SubTargetAndRanges) {
	if linkStore == nil {
		return
	}
	ls := *linkStore
	l, found := ls[id]
	if !found {
		return
	}

	for subTarget, rng := range l.Def {
		subTargets = append(subTargets, SubTargetAndRanges{
			subTarget: subTarget,
			rng:       &rng,
		})
	}
	return subTargets
}

func (store *Store) GetRefsFromTarget(target Target, subTarget SubTarget) (refs []IdLocation, refFound bool) {
	if store == nil {
		return
	}
	s := *store
	ids := s.getIds(target)
	for _, id := range ids {
		lrefs, found := s.LinkStore.GetRefs(id, subTarget)
		if found {
			refFound = true
			refs = append(refs, lrefs...)
		}
	}

	return
}

func (store *Store) GetDefsFromTarget(target Target, subTarget SubTarget) (defs []IdLocation, defFound bool) {
	if store == nil {
		return
	}
	s := *store
	ids, defFound := s.getValidIds(target)
	if len(subTarget) == 0 {
		for _, id := range ids {
			defs = append(defs, IdLocation{
				Id:    id,
				Range: lsp.Range{},
			})
		}
	} else {
		for _, id := range ids {
			def, found := s.LinkStore.GetDef(id, subTarget)
			if found {
				defFound = true
				defs = append(defs, IdLocation{
					Id:    id,
					Range: def,
				})
			}
		}
	}

	return
}

func (s *LinkStore) AddFileGTarget(id Id) {
	s.AddDef(id, "", lsp.Range{})
}
