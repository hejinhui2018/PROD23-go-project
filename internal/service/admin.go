package service

import (
	"prod23/internal/domain"
	"sort"
	"strings"
	"time"
)

type FaultFilter struct {
	Feeder, Status, Device string
	Since                  time.Time
	Limit                  int
}

func (s *Service) ListFaults(f FaultFilter) []domain.Fault {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := []domain.Fault{}
	for _, x := range s.Faults {
		if f.Feeder != "" && x.Feeder != f.Feeder {
			continue
		}
		if f.Status != "" && string(x.Status) != f.Status {
			continue
		}
		if f.Device != "" && !strings.Contains(x.Device, f.Device) {
			continue
		}
		if !f.Since.IsZero() && x.CreatedAt.Before(f.Since) {
			continue
		}
		r = append(r, *x)
	}
	sort.Slice(r, func(i, j int) bool { return r[i].CreatedAt.Before(r[j].CreatedAt) })
	if f.Limit > 0 && len(r) > f.Limit {
		r = r[:f.Limit]
	}
	return r
}

type EventFilter struct {
	Type, Aggregate string
	Since           time.Time
	Limit           int
}

func (s *Service) ListEvents(f EventFilter) []domain.Event {
	all := s.Store.Events()
	r := []domain.Event{}
	for _, e := range all {
		if f.Type != "" && e.Type != f.Type {
			continue
		}
		if f.Aggregate != "" && e.AggregateID != f.Aggregate {
			continue
		}
		if !f.Since.IsZero() && e.At.Before(f.Since) {
			continue
		}
		r = append(r, e)
	}
	if f.Limit > 0 && len(r) > f.Limit {
		r = r[len(r)-f.Limit:]
	}
	return r
}
func (s *Service) PlanProgress(id string) (int, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.Plans[id]
	if p == nil {
		return 0, 0, domain.ErrNotFound
	}
	done := 0
	for _, sid := range p.Steps {
		if s.Steps[sid].Status == domain.Completed {
			done++
		}
	}
	return done, len(p.Steps), nil
}
func (s *Service) StepByID(id string) (domain.Step, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x := s.Steps[id]
	if x == nil {
		return domain.Step{}, domain.ErrNotFound
	}
	return *x, nil
}
func (s *Service) FindStepBySection(plan, section string) (domain.Step, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.Plans[plan]
	if p == nil {
		return domain.Step{}, domain.ErrNotFound
	}
	for _, id := range p.Steps {
		if s.Steps[id].Section == section {
			return *s.Steps[id], nil
		}
	}
	return domain.Step{}, domain.ErrNotFound
}
func (s *Service) RebuildIndexes() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.Plans {
		for i, id := range p.Steps {
			if st := s.Steps[id]; st != nil {
				st.PlanID = p.ID
				st.Sequence = i + 1
			}
		}
	}
}
func (s *Service) Verify() error {
	if !s.Store.Integrity() {
		return domain.ErrConflict
	}
	for _, f := range s.Faults {
		if err := domain.ValidateFault(*f); err != nil {
			return err
		}
	}
	for _, p := range s.Plans {
		if err := domain.ValidatePlan(*p); err != nil {
			return err
		}
	}
	for _, st := range s.Steps {
		if err := domain.ValidateStep(*st); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) TouchReview(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Reviews {
		if s.Reviews[i].ID == id {
			s.Reviews[i].CreatedAt = time.Now().UTC()
			return nil
		}
	}
	return domain.ErrNotFound
}
func (s *Service) FaultCount() int { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.Faults) }
func (s *Service) PlanCount() int  { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.Plans) }
func (s *Service) StepCount() int  { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.Steps) }
func (s *Service) StatusCounts() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := map[string]int{}
	for _, f := range s.Faults {
		m[string(f.Status)]++
	}
	for _, st := range s.Steps {
		m[string(st.Status)]++
	}
	return m
}
