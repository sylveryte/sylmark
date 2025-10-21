package data

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

// heading target
type SubTarget string

// file target
type Target string

// file + heading target
type FullTarget string
type TargetStore map[Target][]Id

func NewTargetStore() TargetStore {
	return TargetStore{}
}
func (t Target) GetFileName() string {
	return string(t) + ".md"
}
func (store *TargetStore) Print() {
	slog.Info("TargetStore>>>>>>>>>>>>>>>>>>>")
	for k, j := range *store {
		slog.Info(fmt.Sprintf("\n[%s]=%d", k, j))
	}
	slog.Info("TargetStore<<<<<<<<<<<<<<")
}

// simple map operation
func (ws *TargetStore) fetchIds(target Target) ([]Id, bool) {
	s := *ws
	ids, ok := s[target]
	return ids, ok
}

func (s *Store) addTargetEntry(target Target, id Id) (ids []Id, isMultiple bool) {
	vt := s.TargetStore
	s.IdStore.addShadowEntry(target, id)
	ids, found := vt[target]
	// utils.Sprintf("addTargetStoreEntry target=[%s] id=[%d] [%v]", target, id, found)
	if !found {
		// doesn't exists
		ids = []Id{id}
		vt[target] = ids
		return ids, false
	}

	alreadyIdExists := slices.Contains(ids, id)
	if !alreadyIdExists {
		ids = append(ids, id)
		vt[target] = ids
	}
	return ids, true
}

// if target has '/' then addd down variants
func (s *Store) addDownVariants(target Target, id Id) {
	// utils.Sprintf(" addDownVariants target is [%s]", target)
	if strings.ContainsRune(string(target), '/') {
		s.addTargetEntry(target, id)

		// utils.Sprintf("  Decision need to addDownVariants target is [%s]", target)
		// assume it's vault target
		oneUpTarget, _ := GetOneUpTarget(target)
		s.addTargetEntry(oneUpTarget, id)
		// get plain target as well
		target, _ := GetPlainTarget(oneUpTarget)
		s.addTargetEntry(target, id)
	}
}
