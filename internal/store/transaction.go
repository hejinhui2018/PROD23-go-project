package store

import (
	"prod23/internal/domain"
	"sync"
)

type Transaction struct {
	mu     sync.Mutex
	store  *Store
	events []domain.Event
	closed bool
}

func (s *Store) Begin() *Transaction { return &Transaction{store: s} }
func (t *Transaction) Append(e domain.Event) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return domain.ErrInvalid
	}
	t.events = append(t.events, e)
	return nil
}
func (t *Transaction) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return domain.ErrInvalid
	}
	for _, e := range t.events {
		if err := t.store.Append(e); err != nil {
			return err
		}
	}
	t.closed = true
	return nil
}
func (t *Transaction) Rollback() { t.mu.Lock(); defer t.mu.Unlock(); t.events = nil; t.closed = true }
func (t *Transaction) Len() int  { t.mu.Lock(); defer t.mu.Unlock(); return len(t.events) }
