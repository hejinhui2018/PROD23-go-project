package service

import (
	"prod23/internal/store"
	"testing"
)

func TestStepCompletionCanBeRetriedAfterRestart(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	s := New(st)
	fault, err := s.CreateFault("FDR-1", "SW-1")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := s.Confirm(fault.ID)
	if err != nil {
		t.Fatal(err)
	}
	stepID := plan.Steps[0]
	if _, err = s.ClaimStep(stepID, "crew-1", "claim-request", 1); err != nil {
		t.Fatal(err)
	}
	if _, err = s.AcknowledgeStep(stepID, "crew-1", "ack-request", 2); err != nil {
		t.Fatal(err)
	}
	completed, err := s.CompleteStep(stepID, "crew-1", "completion-request", 3, "ok")
	if err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	restarted := New(reopened)
	retried, err := restarted.CompleteStep(stepID, "crew-1", "completion-request", 3, "ok")
	if err != nil {
		t.Fatal(err)
	}
	if retried.ID != completed.ID || retried.Status != completed.Status || retried.Version != completed.Version {
		t.Fatalf("retry changed completed step: first=%+v retry=%+v", completed, retried)
	}
	if got := len(restarted.Audit()); got != len(s.Audit()) {
		t.Fatalf("retry appended an event: before=%d after=%d", len(s.Audit()), got)
	}
}
