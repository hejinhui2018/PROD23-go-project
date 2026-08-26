package store

import (
	"os"
	"path/filepath"
	"prod23/internal/domain"
	"time"
)

type Manager struct {
	Root   string
	Events *Store
	Queue  *Queue
}

func OpenManager(root string) (*Manager, error) {
	if e := os.MkdirAll(root, 0755); e != nil {
		return nil, e
	}
	s, e := Open(root)
	if e != nil {
		return nil, e
	}
	return &Manager{Root: root, Events: s, Queue: LoadQueue(root)}, nil
}
func (m *Manager) Append(e domain.Event) error { return m.Events.Append(e) }
func (m *Manager) EventCount() int             { return len(m.Events.Events()) }
func (m *Manager) LastEvent() (domain.Event, bool) {
	all := m.Events.Events()
	if len(all) == 0 {
		return domain.Event{}, false
	}
	return all[len(all)-1], true
}
func (m *Manager) RotateSnapshot(v Snapshot) error {
	tmp := filepath.Join(m.Root, "snapshot.tmp")
	if e := os.MkdirAll(tmp, 0755); e != nil {
		return e
	}
	if e := m.Events.SaveSnapshot(tmp, v); e != nil {
		return e
	}
	return os.Rename(filepath.Join(tmp, "snapshot.json"), filepath.Join(m.Root, "snapshot.json"))
}
func (m *Manager) Age() time.Duration {
	st, e := os.Stat(filepath.Join(m.Root, "events.jsonl"))
	if e != nil {
		return 0
	}
	return time.Since(st.ModTime())
}
func (m *Manager) Health() map[string]any {
	return map[string]any{"events": m.EventCount(), "queue": len(m.Queue.All()), "age_seconds": m.Age().Seconds()}
}
