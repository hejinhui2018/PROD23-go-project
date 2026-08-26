package service

import (
	"prod23/internal/domain"
	"sort"
)

type PlanSummary struct {
	Plan                       domain.Plan
	Completed, Failed, Pending int
}

func (s *Service) PlanSummaries() []PlanSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := make([]PlanSummary, 0, len(s.Plans))
	for _, p := range s.Plans {
		x := PlanSummary{Plan: *p}
		for _, id := range p.Steps {
			switch s.Steps[id].Status {
			case domain.Completed:
				x.Completed++
			case domain.Failed:
				x.Failed++
			default:
				x.Pending++
			}
		}
		r = append(r, x)
	}
	sort.Slice(r, func(i, j int) bool { return r[i].Plan.ID < r[j].Plan.ID })
	return r
}
func (s *Service) StepsForWorker(worker string) []domain.Step {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := []domain.Step{}
	for _, st := range s.Steps {
		if st.ClaimedBy == worker && !domain.IsTerminal(st.Status) {
			r = append(r, *st)
		}
	}
	sort.Slice(r, func(i, j int) bool { return r[i].Sequence < r[j].Sequence })
	return r
}
func (s *Service) EventCount(typ string) int {
	n := 0
	for _, e := range s.Store.Events() {
		if typ == "" || e.Type == typ {
			n++
		}
	}
	return n
}
