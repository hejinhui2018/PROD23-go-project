package service

import (
	"prod23/internal/domain"
	"prod23/internal/store"
	"testing"
)

func TestFailedStepReviewSurvivesRestart(t *testing.T) {
	d := t.TempDir()
	firstStore, err := store.Open(d)
	if err != nil {
		t.Fatal(err)
	}
	first := New(firstStore)
	fault, err := first.CreateFault("FDR-1", "SW-1")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := first.Confirm(fault.ID)
	if err != nil {
		t.Fatal(err)
	}
	stepID := plan.Steps[0]
	if _, err = first.ClaimStep(stepID, "crew-7", "claim-review", 1); err != nil {
		t.Fatal(err)
	}
	if _, err = first.AcknowledgeStep(stepID, "crew-7", "ack-review", 2); err != nil {
		t.Fatal(err)
	}
	failed, err := first.FailStep(stepID, "crew-7", "fail-review", 3, "access road blocked")
	if err != nil {
		t.Fatal(err)
	}
	if failed.Status != domain.Failed || len(first.OpenReviews()) != 1 {
		t.Fatalf("failure was not queued for review: step=%+v reviews=%d", failed, len(first.OpenReviews()))
	}

	restartedStore, err := store.Open(d)
	if err != nil {
		t.Fatal(err)
	}
	restarted := New(restartedStore)
	recovered, err := restarted.StepByID(stepID)
	if err != nil {
		t.Fatal(err)
	}
	reviews := restarted.OpenReviews()
	if recovered.Status != domain.Failed || recovered.Result != "access road blocked" {
		t.Fatalf("recovered step lost failure details: %+v", recovered)
	}
	if len(reviews) != 1 || reviews[0].StepID != stepID || reviews[0].Reason != "access road blocked" {
		t.Fatalf("recovered review queue mismatch: %+v", reviews)
	}
}
