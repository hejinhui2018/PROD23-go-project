package service

import (
	"os"
	"prod23/internal/store"
	"testing"
)

func TestFlow(t *testing.T) {
	d := t.TempDir()
	s, _ := store.Open(d)
	x := New(s)
	f, _ := x.CreateFault("FDR-1", "SW-1")
	p, _ := x.Confirm(f.ID)
	for _, id := range p.Steps {
		x.Step(id, "crew", "claim", 1)
		x.Step(id, "crew", "ack", 2)
		x.Step(id, "crew", "complete", 3)
	}
	if x.Faults[f.ID].Status != "restored" {
		t.Fatal("not restored")
	}
	_ = os.RemoveAll(d)
}
