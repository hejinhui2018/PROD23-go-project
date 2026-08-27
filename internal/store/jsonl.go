package store

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"prod23/internal/domain"
	"sync"
	"time"
)

type Store struct {
	mu     sync.Mutex
	path   string
	events []domain.Event
}

// persistedEvent is the JSONL wire shape used by the append-only log.
type persistedEvent struct {
	ID, Type    string
	AggregateID string
	Version     int
	Data        any
	At          time.Time
	Idempotency string
}

func (s *Store) Dir() string { return filepath.Dir(s.path) }

func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	s := &Store{path: filepath.Join(dir, "events.jsonl")}
	f, e := os.OpenFile(s.path, os.O_CREATE|os.O_RDONLY, 0644)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		var v persistedEvent
		if json.Unmarshal(line, &v) == nil {
			s.events = append(s.events, domain.Event{ID: v.ID, Type: v.Type, AggregateID: v.AggregateID, Version: v.Version, Data: v.Data, At: v.At, Idempotency: v.Idempotency})
		} else if _, x := f.Seek(0, io.SeekCurrent); x == nil {
			break
		}
	}
	if e := sc.Err(); e != nil {
		return nil, e
	}
	return s, nil
}
func (s *Store) Append(e domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(persistedEvent{ID: e.ID, Type: e.Type, AggregateID: e.AggregateID, Version: e.Version, Data: e.Data, At: e.At, Idempotency: e.Idempotency})
	if err != nil {
		return err
	}
	if _, err = f.Write(append(b, '\n')); err != nil {
		return err
	}
	s.events = append(s.events, e)
	return f.Sync()
}
func (s *Store) Events() []domain.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]domain.Event(nil), s.events...)
}
