package data

import (
	"fmt"
	"log/slog"
	"slices"
	"sylmark/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Id uint

type IdStore struct {
	Id            map[Id]lsp.DocumentURI
	uri           map[lsp.DocumentURI]Id
	ShadowTargets map[Id][]Target
}

func NewIdStore() IdStore {
	return IdStore{
		Id:            map[Id]lsp.DocumentURI{},
		uri:           map[lsp.DocumentURI]Id{},
		ShadowTargets: map[Id][]Target{},
	}
}
func (s *IdStore) Print() {
	slog.Info("IdStore===============>>>>>>>>>>")
	slog.Info("IdStore====id")
	for k, v := range s.Id {
		slog.Info(fmt.Sprintf("[%d]=[%s]", k, v))
	}
	slog.Info("IdStore====uri")
	for k, v := range s.uri {
		slog.Info(fmt.Sprintf("[%s]=[%d]", k, v))
	}
	slog.Info("IdStore====shadowTargets")
	for k, v := range s.ShadowTargets {
		slog.Info(fmt.Sprintf("[%d]=[%s]", k, v))
	}
	slog.Info("IdStore===============<<<<<<<<<<<<")
}
func (s *IdStore) ReplaceUri(id Id, uri lsp.DocumentURI) {
	// utils.Sprintf("ReplaceUri       id=[%d] uri=[%s]", id, uri)
	s.Id[id] = uri
	s.uri[uri] = id

	// cleanup shadow
	delete(s.ShadowTargets, id)
}

func (s *IdStore) addShadowEntry(target Target, id Id) {
	targets, found := s.ShadowTargets[id]
	// utils.Sprintf("addShadowEntry target=[%s] id=[%d] [%v]", target, id, found)
	if !found {
		// doesn't exists, no need to add, since shadowTargets are already in
		return
	}

	alreadyIdExists := slices.Contains(targets, target)
	if !alreadyIdExists {
		targets = append(targets, target)
		s.ShadowTargets[id] = targets
	}
}

// find if all shadowId refs are compatible with uri
func (s *Store) isClaimableShadow(shadowId Id, uri lsp.DocumentURI) bool {
	targets, found := s.IdStore.ShadowTargets[shadowId]
	if !found {
		return false
	}
	vaultTarget, ok := s.GetVaultTarget(uri)
	if ok {
		i := slices.Index(targets, vaultTarget)
		if i > -1 {
			targets = slices.Delete(targets, i, i+1)
		}
	}

	oneUpTarget, diff := GetOneUpTarget(vaultTarget)
	if diff {
		i := slices.Index(targets, oneUpTarget)
		if i > -1 {
			targets = slices.Delete(targets, i, i+1)
		}
	}

	plainTarget, diff := GetPlainTarget(vaultTarget)
	if diff {
		i := slices.Index(targets, plainTarget)
		if i > -1 {
			targets = slices.Delete(targets, i, i+1)
		}
	}

	return len(targets) == 0
}

// doesn't check if exists same uri, only store
func (s *IdStore) addEntry(uri lsp.DocumentURI) Id {
	id := Id(len(s.Id) + 1)
	// utils.Sprintf(" IdStore addEntry uri=[%s] id=[%d]", uri, id)
	// utils.Sprintf("addEntry uri=[%s] id=[%d] ", uri, id)
	s.Id[id] = uri
	if len(uri) > 0 {
		s.uri[uri] = id
	} else {
		s.ShadowTargets[id] = []Target{}
	}
	return id
}

// id with no URI
func (s *IdStore) isShadowId(id Id) bool {
	_, ok := s.ShadowTargets[id]
	return ok
}

// id with no URI
func (s *IdStore) findShadowId(ids []Id) (Id, bool) {
	for _, id := range ids {
		if s.isShadowId(id) {
			return id, true
		}
	}
	return 0, false
}

// UTILS for IdStore
type IdLocation struct {
	Id    Id
	Range lsp.Range
}

func (d Id) LocationOf(node *tree_sitter.Node) IdLocation {
	return IdLocation{
		Id:    d,
		Range: lsp.GetRange(node),
	}
}
