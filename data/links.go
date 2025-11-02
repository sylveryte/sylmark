package data

import (
	"strings"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// any atx heading node, # included in SubTarget
func GetSubTarget(node *tree_sitter.Node, content string) (subTarget SubTarget, ok bool) {
	if node.Kind() != "atx_heading" {
		node = node.Parent()
	}
	if node.Kind() != "atx_heading" {
		return "", false
	}

	if node.NamedChildCount() == 2 {
		lineNode := node.NamedChild(1)
		t := lsp.GetNodeContent(*lineNode, content)
		subTarget = SubTarget("#" + t)
		return subTarget, true
	}
	return
}

// wiki-link node only, subTarget includes '#'
func GetWikilinkTargets(node *tree_sitter.Node, content string) (target Target, subTarget SubTarget, isSubTarget bool, ok bool) {
	linkDestNode := node.NamedChild(0)
	t := lsp.GetNodeContent(*linkDestNode, content)
	ok = true
	if strings.ContainsRune(t, '#') {
		isSubTarget = true
		hashI := strings.IndexRune(t, '#')
		subTarget = SubTarget(t[hashI:])
		target = Target(t[:hashI])
	} else {
		target = Target(t)
	}
	return
}

// syltodo work needed
// could we unload and load entire file instead ??
func (s *Store) ReplaceTarget(id Id, oldTarget Target, target Target) {
	if oldTarget == target {
		return
	}
	ts := s.TargetStore
	ids := ts[oldTarget]
	switch len(ids) {
	case 0:
		{
			ts[target] = []Id{id}
		}
	case 1:
		{
			delete(ts, oldTarget)
			ts[target] = []Id{id}
		}
	default:
		{
			var newIds []Id
			for _, i := range ids {
				if i == id {
					continue
				}
				newIds = append(newIds, i)
			}
			ts[oldTarget] = newIds
			s.addTargetEntry(target, id)

		}
	}
}

// gets uri from id
func (s *Store) GetUri(id Id) (lsp.DocumentURI, bool) {
	uri, ok := s.IdStore.Id[id]
	return uri, ok
}

func (s *Store) FillInLocations(locs *[]lsp.Location, idLocs *[]IdLocation) *[]lsp.Location {
	for _, il := range *idLocs {
		l, ok := s.LocationFromIdLocation(il)
		if ok {
			*locs = append(*locs, *l)
		}
	}
	return locs
}

func (s *Store) IdLocationFromLocation(loc lsp.Location) IdLocation {
	return IdLocation{
		Id:    s.GetIdFromURI(loc.URI),
		Range: loc.Range,
	}
}

func (s *Store) LocationFromIdLocation(loc IdLocation) (*lsp.Location, bool) {
	uri, ok := s.GetUri(loc.Id)
	if ok {
		return &lsp.Location{
			URI:   uri,
			Range: loc.Range,
		}, true
	}
	return nil, false
}

// matches Id case insensitve and returns id
func (s *Store) findIdFromURIFold(uri lsp.DocumentURI) (id Id, found bool) {
	id, found = s.findIdFromURI(uri)
	if found {
		return
	}
	for u, fid := range s.IdStore.uri {
		m := strings.EqualFold(string(uri), string(u))
		if m {
			id = fid
			break
		}
	}

	return id, id != 0
}

func (s *Store) findIdFromURI(uri lsp.DocumentURI) (id Id, found bool) {
	id, found = s.IdStore.uri[uri]
	return
}

/// Store and link add business

// gets new id creates if not already, makes shadow id real
func (s *Store) GetIdFromURI(uri lsp.DocumentURI) Id {
	id, ok := s.IdStore.uri[uri]
	// utils.Sprintf(">>>GetIdFromURI uri=[%s] id=[%d] [%v]", uri, id, ok)
	if ok {
		return id
	} else {

		// Fact: down variants have same shadow id, because of addDownVariants
		// So better start at top, if found then that is infact defacto id for uri

		target, _ := GetTarget(uri)
		vaultTarget, _ := s.GetVaultTarget(uri)
		oneUpTarget, _ := GetOneUpTarget(vaultTarget)
		var shadowId Id
		var isShadow bool

		// utils.Sprintf("   target=[%s] vaultTarget=[%s] oneUpTarget=[%s]", target, vaultTarget, oneUpTarget)

		ids, ok := s.TargetStore.fetchIds(vaultTarget)
		// utils.Sprintf("  trying vaultTarget %d", ids)
		shadowId, isShadow = s.IdStore.findShadowId(ids)
		if ok && isShadow {
			goto proceed
		}

		ids, ok = s.TargetStore.fetchIds(oneUpTarget)
		// utils.Sprintf("  trying oneUpTarget %d", ids)
		shadowId, isShadow = s.IdStore.findShadowId(ids)
		if ok && isShadow {
			goto proceed
		}

		ids, ok = s.TargetStore.fetchIds(target)
		// utils.Sprintf("  trying target %d", ids)
		shadowId, isShadow = s.IdStore.findShadowId(ids)
		if ok && isShadow {
			goto proceed
		}

	proceed:
		if isShadow && s.isClaimableShadow(shadowId, uri) {
			// utils.Sprintf("     Decision isShadow ReplaceUri shadowId=%d", shadowId)
			s.IdStore.ReplaceUri(shadowId, uri)
			return shadowId
		}

		// no shadow add new
		newId := s.IdStore.addEntry(uri)
		// utils.Sprintf("     Decision NO shadow newId=%d ids=%d", newId, ids)
		if len(ids) > 0 {
			// some exists add new and existings' variants
			// utils.Sprintf("     Decision has hasSomeIds ids=%d", ids)

			// variants if only one other existing, else will have already
			// by this same algo
			if len(ids) == 1 {
				s.addVariants(ids[0])
			}

			// addVariants
			s.addTargetEntry(vaultTarget, newId)
			s.addTargetEntry(oneUpTarget, newId)
			s.addTargetEntry(target, newId)
		} else {
			// utils.Sprintf("     Totally new target=[%s] newId=[%d]", target, newId)
			s.addTargetEntry(target, newId)
		}
		return newId
	}
}

// filtered real ids from getIds
func (s *Store) getValidIds(target Target) ([]Id, bool) {
	ids := s.getIds(target)
	var validIds []Id
	for _, id := range ids {
		if !s.IdStore.isShadowId(id) {
			validIds = append(validIds, id)
		}
	}
	return validIds, len(validIds) > 0
}

// creates id if doesn't exists, adds variants in case of non plain
func (s *Store) getIds(target Target) []Id {
	ts := s.TargetStore
	ids, ok := ts[target]
	// utils.Sprintf(">>>getIds target=[%s] ids=[%d] [%v]", target, ids, ok)
	if !ok {
		// utils.Sprintf("   getIds target=[%s]", target)
		if strings.ContainsRune(string(target), '/') {
			// variant route
			// check if downVariants have any shadowId
			// assuming target to be vaultTarget
			ids, _ := s.TargetStore.fetchIds(target)
			id, isShadow := s.IdStore.findShadowId(ids)
			// utils.Sprintf("   getIds variants route ids=%d isShadow=%v", ids, isShadow)
			if !isShadow {
				oneUpTarget, isDiff := GetOneUpTarget(target)
				// utils.Sprintf("   getIds variants route target=%s  oneUpTarget=%s  isDiff=%v ids=%d isShadow=%v", target, oneUpTarget, isDiff, ids, isShadow)
				if isDiff {
					ids, ok = s.TargetStore.fetchIds(oneUpTarget)
					id, isShadow = s.IdStore.findShadowId(ids)
					if !isShadow {
						plainTarget, isDiff := GetPlainTarget(target)
						// utils.Sprintf("         getIds variants route target=%s  plainTarget=%s  isDiff=%v ids=%d isShadow=%v", target, plainTarget, isDiff, ids, isShadow)
						if isDiff {
							ids, ok = s.TargetStore.fetchIds(plainTarget)
							id, isShadow = s.IdStore.findShadowId(ids)
						}
					}
				}
			}
			// second time check since possibly changed value of isShadow
			if !isShadow {
				id = s.IdStore.addEntry(lsp.DocumentURI(""))
			}
			// utils.Sprintf("   getIds id=[%d]", id)
			// add down variants
			ts[target] = []Id{id}
			s.addDownVariants(target, id)
		} else {
			id := s.IdStore.addEntry(lsp.DocumentURI(""))
			ids = []Id{id}
			// new id get from IdsStore
			// utils.Sprintf("    getIds newEntry [%s]=[%d]", target, id)
			s.addTargetEntry(target, id)
		}
	}
	return ids
}

func (s *Store) addVariants(id Id) {
	uri, found := s.GetUri(id)
	// utils.Sprintf("targetStoreAddFullPathVariant id=[%d] uri=[%s] [%v]", id, uri, found)
	if !found {
		return
	}
	target, ok := GetTarget(uri)
	if !ok {
		return
	}
	s.addTargetEntry(target, id)
	vaultTarget, isDiff := s.GetVaultTarget(uri)
	if isDiff {
		s.addTargetEntry(vaultTarget, id)
	}
	oneUpTarget, isDiff := GetOneUpTarget(vaultTarget)
	if isDiff {
		s.addTargetEntry(oneUpTarget, id)
	}
}
