package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"prod23/internal/domain"
)

func EventDigest(e domain.Event) string {
	b, _ := json.Marshal(e)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func ChainDigest(events []domain.Event) string {
	h := sha256.New()
	for _, e := range events {
		h.Write([]byte(EventDigest(e)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
func (s *Store) Integrity() bool { return ValidateSequence(s.Events()) }
func (s *Store) AggregateCount() int {
	m := map[string]bool{}
	for _, e := range s.Events() {
		m[e.AggregateID] = true
	}
	return len(m)
}
