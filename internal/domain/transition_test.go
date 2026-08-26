package domain

import "testing"

func TestTransitionRules(t *testing.T) {
	if !ValidTransition(Pending, Dispatched) || ValidTransition(Pending, Completed) {
		t.Fatal("transition matrix")
	}
}
