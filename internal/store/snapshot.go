package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"prod23/internal/domain"
)

type Snapshot struct {
	Faults        []domain.Fault
	Plans         []domain.Plan
	Steps         []domain.Step
	Notifications []Notification
}

func (s *Store) SaveSnapshot(dir string, v Snapshot) error {
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return e
	}
	return os.WriteFile(filepath.Join(dir, "snapshot.json"), b, 0644)
}
func LoadSnapshot(dir string) (Snapshot, error) {
	var v Snapshot
	b, e := os.ReadFile(filepath.Join(dir, "snapshot.json"))
	if e != nil {
		return v, e
	}
	e = json.Unmarshal(b, &v)
	return v, e
}
