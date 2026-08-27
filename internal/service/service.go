package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"prod23/internal/domain"
	"prod23/internal/store"
	"sync"
	"time"
)

func newID() string { b := make([]byte, 12); _, _ = rand.Read(b); return hex.EncodeToString(b) }

type Service struct {
	mu       sync.RWMutex
	Store    *store.Store
	Topology domain.Topology
	Faults   map[string]*domain.Fault
	Plans    map[string]*domain.Plan
	Steps    map[string]*domain.Step
	Queue    *store.Queue
	Reviews  []domain.ReviewRecord
	idem     map[string]any
	Metrics  *ServiceMetrics
}

func New(st *store.Store) *Service {
	s := &Service{Store: st, Topology: domain.DefaultTopology(), Faults: map[string]*domain.Fault{}, Plans: map[string]*domain.Plan{}, Steps: map[string]*domain.Step{}, Queue: store.LoadQueue(st.Dir()), idem: map[string]any{}, Metrics: &ServiceMetrics{Started: time.Now().UTC()}}
	s.replay()
	return s
}
func (s *Service) event(t, id string, v int, d any, key string) error {
	return s.Store.Append(domain.Event{ID: newID(), Type: t, AggregateID: id, Version: v, Data: d, At: time.Now().UTC(), Idempotency: key})
}
func (s *Service) CreateFault(a, b string) (*domain.Fault, error) {
	return s.CreateFaultWithKey(a, b, "")
}
func (s *Service) CreateFaultWithKey(feeder, device, key string) (*domain.Fault, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if feeder == "" || device == "" {
		return nil, domain.ErrInvalid
	}
	if key != "" {
		if v, ok := s.idem[key].(*domain.Fault); ok {
			return v, nil
		}
	}
	f := &domain.Fault{ID: newID(), Feeder: feeder, Device: device, Sections: s.Topology.Affected(feeder, device), Status: domain.Assessed, Version: 1, CreatedAt: time.Now().UTC()}
	if len(f.Sections) == 0 {
		return nil, fmt.Errorf("unknown feeder or device")
	}
	s.Faults[f.ID] = f
	s.Metrics.Faults.Add(1)
	if key != "" {
		s.idem[key] = f
	}
	return f, s.event("fault.assessed", f.ID, 1, *f, key)
}
func (s *Service) GetFault(id string) (domain.Fault, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f := s.Faults[id]
	if f == nil {
		return domain.Fault{}, domain.ErrNotFound
	}
	return *f, nil
}
func (s *Service) Confirm(id string) (*domain.Plan, error) { return s.ConfirmWithKey(id, "") }
func (s *Service) ConfirmWithKey(id, key string) (*domain.Plan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key != "" {
		if p, ok := s.idem[key].(*domain.Plan); ok {
			return p, nil
		}
	}
	f := s.Faults[id]
	if f == nil {
		return nil, domain.ErrNotFound
	}
	for _, p := range s.Plans {
		if p.FaultID == id {
			return p, nil
		}
	}
	p := &domain.Plan{ID: newID(), FaultID: id, Status: "confirmed", Version: 1}
	for i, sec := range f.Sections {
		st := &domain.Step{ID: newID(), PlanID: p.ID, Section: sec, Sequence: i + 1, Status: domain.Pending, Version: 1}
		s.Steps[st.ID] = st
		p.Steps = append(p.Steps, st.ID)
	}
	s.Plans[p.ID] = p
	s.Metrics.Plans.Add(1)
	if key != "" {
		s.idem[key] = p
	}
	_ = s.event("plan.confirmed", p.ID, 1, *p, key)
	return p, nil
}
func (s *Service) GetPlan(id string) (domain.Plan, []domain.Step, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.Plans[id]
	if p == nil {
		return domain.Plan{}, nil, domain.ErrNotFound
	}
	r := make([]domain.Step, 0, len(p.Steps))
	for _, x := range p.Steps {
		r = append(r, *s.Steps[x])
	}
	return *p, r, nil
}
func (s *Service) ClaimStep(id, w, k string, v int) (*domain.Step, error) {
	return s.mutateStep(id, w, "claim", k, v, "")
}
func (s *Service) AcknowledgeStep(id, w, k string, v int) (*domain.Step, error) {
	return s.mutateStep(id, w, "ack", k, v, "")
}
func (s *Service) CompleteStep(id, w, k string, v int, r string) (*domain.Step, error) {
	return s.mutateStep(id, w, "complete", k, v, r)
}
func (s *Service) FailStep(id, w, k string, v int, r string) (*domain.Step, error) {
	return s.mutateStep(id, w, "fail", k, v, r)
}
func (s *Service) Step(id, w, a string, v int) (*domain.Step, error) {
	return s.mutateStep(id, w, a, "", v, a)
}
func (s *Service) mutateStep(id, w, a, k string, v int, res string) (*domain.Step, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if k != "" {
		if x, ok := s.idem[k].(*domain.Step); ok {
			return x, nil
		}
	}
	st := s.Steps[id]
	if st == nil {
		return nil, domain.ErrNotFound
	}
	if w == "" || st.Version != v {
		if st.Version != v {
			s.Metrics.Conflicts.Add(1)
		}
		return nil, domain.ErrConflict
	}
	var openedReview *domain.ReviewRecord
	if a == "claim" {
		if st.Status != domain.Pending {
			return nil, domain.ErrInvalid
		}
		p := s.Plans[st.PlanID]
		for _, sid := range p.Steps {
			if s.Steps[sid].Sequence < st.Sequence && s.Steps[sid].Status != domain.Completed {
				return nil, domain.ErrPrecondition
			}
		}
		st.Status = domain.Dispatched
		s.Metrics.Claims.Add(1)
		st.ClaimedBy = w
	} else if a == "ack" {
		if st.Status != domain.Dispatched || st.ClaimedBy != w {
			return nil, domain.ErrPrecondition
		}
		st.Status = domain.Acknowledged
		s.Metrics.Acks.Add(1)
	} else if a == "complete" || a == "fail" {
		if st.Status != domain.Acknowledged || st.ClaimedBy != w {
			return nil, domain.ErrPrecondition
		}
		st.Result = res
		if a == "complete" {
			st.Status = domain.Completed
			s.Metrics.Completions.Add(1)
		} else {
			st.Status = domain.Failed
			s.Metrics.Failures.Add(1)
			review := domain.ReviewRecord{ID: newID(), StepID: id, Reason: res, Worker: w, CreatedAt: time.Now().UTC()}
			s.Reviews = append(s.Reviews, review)
			openedReview = &review
		}
	} else {
		return nil, domain.ErrInvalid
	}
	st.Version++
	cp := *st
	if k != "" {
		s.idem[k] = &cp
	}
	eventData := any(cp)
	if openedReview != nil {
		eventData = store.NewStepFailureEvent(cp, *openedReview)
	}
	_ = s.event("step."+string(st.Status), st.ID, st.Version, eventData, k)
	if st.Status == domain.Completed {
		s.checkRestored(st.PlanID)
	}
	return st, nil
}
func (s *Service) checkRestored(pid string) {
	p := s.Plans[pid]
	for _, sid := range p.Steps {
		if s.Steps[sid].Status != domain.Completed {
			return
		}
	}
	f := s.Faults[p.FaultID]
	if f.Status == domain.Restored {
		return
	}
	f.Status = domain.Restored
	f.Version++
	_ = s.event("fault.restored", f.ID, f.Version, *f, "")
	s.Queue.Enqueue(store.Notification{ID: newID(), Type: "restoration.completed", Payload: f.ID, CreatedAt: time.Now().UTC()})
}
func (s *Service) ReviewsList() []domain.ReviewRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.ReviewRecord(nil), s.Reviews...)
}
