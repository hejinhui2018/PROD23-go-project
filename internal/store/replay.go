package store

import "prod23/internal/domain"

type StepFailureEvent struct {
	Step     domain.Step         `json:"step"`
	ReviewID string              `json:"review_id"`
	Review   domain.ReviewRecord `json:"review"`
}

func NewStepFailureEvent(step domain.Step, review domain.ReviewRecord) StepFailureEvent {
	return StepFailureEvent{Step: step, ReviewID: review.ID, Review: review}
}

func (s *Store) ReplayCount() int { return len(s.events) }
