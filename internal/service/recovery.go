package service

import (
	"encoding/json"
	"prod23/internal/domain"
)

func (s *Service) replay() {
	for _, e := range s.Store.Events() {
		b, _ := json.Marshal(e.Data)
		switch e.Type {
		case "fault.assessed", "fault.restored":
			var f domain.Fault
			if json.Unmarshal(b, &f) == nil {
				s.Faults[f.ID] = &f
			}
		case "plan.confirmed":
			var p domain.Plan
			if json.Unmarshal(b, &p) == nil {
				s.Plans[p.ID] = &p
				for i, id := range p.Steps {
					if _, ok := s.Steps[id]; !ok {
						s.Steps[id] = &domain.Step{ID: id, PlanID: p.ID, Sequence: i + 1, Status: domain.Pending, Version: 1}
					}
				}
			}
		case "step.dispatched", "step.acknowledged", "step.completed", "step.failed":
			var st domain.Step
			if json.Unmarshal(b, &st) == nil {
				s.Steps[st.ID] = &st
			}
		}
	}
}
func (s *Service) Replay() int { return len(s.Store.Events()) }
