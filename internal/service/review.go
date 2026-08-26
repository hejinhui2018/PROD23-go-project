package service

import (
	"prod23/internal/domain"
	"sort"
	"time"
)

func (s *Service) OpenReview(stepID, reason, worker string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Steps[stepID] == nil {
		return domain.ErrNotFound
	}
	s.Reviews = append(s.Reviews, domain.ReviewRecord{ID: newID(), StepID: stepID, Reason: reason, Worker: worker, CreatedAt: time.Now().UTC()})
	return nil
}
func (s *Service) ResolveReview(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Reviews {
		if s.Reviews[i].ID == id {
			s.Reviews[i].Resolved = true
			return nil
		}
	}
	return domain.ErrNotFound
}
func (s *Service) OpenReviews() []domain.ReviewRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := []domain.ReviewRecord{}
	for _, x := range s.Reviews {
		if !x.Resolved {
			r = append(r, x)
		}
	}
	sort.Slice(r, func(i, j int) bool { return r[i].CreatedAt.Before(r[j].CreatedAt) })
	return r
}
