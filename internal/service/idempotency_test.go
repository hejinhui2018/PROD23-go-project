package service

import (
	"errors"
	"prod23/internal/domain"
	"prod23/internal/store"
	"testing"
)

func TestIdempotencyAndConflict(t *testing.T) {
	st, _ := store.Open(t.TempDir())
	s := New(st)
	a, _ := s.CreateFaultWithKey("FDR-1", "SW-1", "k")
	b, _ := s.CreateFaultWithKey("FDR-1", "SW-1", "k")
	if a.ID != b.ID {
		t.Fatal("duplicate")
	}
	p, _ := s.Confirm(a.ID)
	stp := p.Steps[0]
	if _, e := s.ClaimStep(stp, "w", "", 99); !errors.Is(e, domain.ErrConflict) {
		t.Fatal(e)
	}
}
