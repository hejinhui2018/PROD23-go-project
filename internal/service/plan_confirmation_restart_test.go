package service

import (
	"prod23/internal/store"
	"testing"
)

func TestPlanConfirmationRetryAfterRestartDoesNotDuplicatePlan(t *testing.T) {
	d := t.TempDir()
	firstStore, err := store.Open(d)
	if err != nil {
		t.Fatal(err)
	}
	first := New(firstStore)
	fault, err := first.CreateFaultWithKey("FDR-1", "SW-1", "fault-request")
	if err != nil {
		t.Fatal(err)
	}
	written, err := first.ConfirmWithKey(fault.ID, "confirm-request")
	if err != nil {
		t.Fatal(err)
	}

	secondStore, err := store.Open(d)
	if err != nil {
		t.Fatal(err)
	}
	second := New(secondStore)
	retried, err := second.ConfirmWithKey(fault.ID, "confirm-request")
	if err != nil {
		t.Fatal(err)
	}

	if retried.ID != written.ID {
		t.Fatalf("confirmation retry created plan %s instead of returning %s", retried.ID, written.ID)
	}
	if second.PlanCount() != 1 || second.StepCount() != len(written.Steps) {
		t.Fatalf("restart retry duplicated execution state: plans=%d steps=%d", second.PlanCount(), second.StepCount())
	}
	if second.EventCount("plan.confirmed") != 1 {
		t.Fatalf("confirmation retry appended duplicate event: %d", second.EventCount("plan.confirmed"))
	}
}
