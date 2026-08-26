package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"prod23/internal/domain"
)

func (s *Store) Compact() error {
	tmp := s.path + ".compact"
	f, e := os.Create(tmp)
	if e != nil {
		return e
	}
	enc := json.NewEncoder(f)
	for _, ev := range s.events {
		if e = enc.Encode(ev); e != nil {
			f.Close()
			return e
		}
	}
	if e = f.Sync(); e != nil {
		f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	return os.Rename(tmp, s.path)
}
func (s *Store) CompactUntil(version int) error {
	keep := []domain.Event{}
	for _, e := range s.events {
		if e.Version >= version {
			keep = append(keep, e)
		}
	}
	old := s.events
	s.events = keep
	if e := s.Compact(); e != nil {
		s.events = old
		return e
	}
	return nil
}
func (s *Store) EventPath() string { return filepath.Clean(s.path) }
