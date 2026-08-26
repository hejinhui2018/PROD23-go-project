package service

import (
	"prod23/internal/store"
	"testing"
)

func TestRestartReplay(t *testing.T) {
	d := t.TempDir()
	st, _ := store.Open(d)
	s := New(st)
	f, _ := s.CreateFault("FDR-1", "SW-1")
	p, _ := s.Confirm(f.ID)
	if len(p.Steps) != 1 {
		t.Fatal("expected one step")
	}
	st2, _ := store.Open(d)
	r := New(st2)
	if _, e := r.GetFault(f.ID); e != nil {
		t.Fatal(e)
	}
	if _, _, e := r.GetPlan(p.ID); e != nil {
		t.Fatal(e)
	}
}
