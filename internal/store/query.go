package store

import (
	"prod23/internal/domain"
	"sort"
	"time"
)

func (s *Store) EventsSince(at time.Time) []domain.Event {
	all := s.Events()
	r := []domain.Event{}
	for _, e := range all {
		if e.At.After(at) {
			r = append(r, e)
		}
	}
	return r
}
func (s *Store) EventsFor(id string) []domain.Event {
	r := []domain.Event{}
	for _, e := range s.Events() {
		if e.AggregateID == id {
			r = append(r, e)
		}
	}
	sort.Slice(r, func(i, j int) bool { return r[i].Version < r[j].Version })
	return r
}
func (s *Store) LastVersion(id string) int {
	v := 0
	for _, e := range s.EventsFor(id) {
		if e.Version > v {
			v = e.Version
		}
	}
	return v
}
func (s *Store) Types() map[string]int {
	m := map[string]int{}
	for _, e := range s.Events() {
		m[e.Type]++
	}
	return m
}
