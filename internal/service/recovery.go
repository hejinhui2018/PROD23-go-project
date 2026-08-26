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
				// Index confirmed plans by fault id so a post-restart
				// confirm for an already-confirmed fault returns the
				// existing plan instead of creating a duplicate. The
				// prior code keyed this index by the idempotency key,
				// which ConfirmWithKey never consults (it queries by
				// fault id), so the dedup index was lost on restart.
				s.confirmedPlans[p.FaultID] = p.ID
				// Rebuild the in-memory idempotency cache from the
				// persisted key so a retried confirm carrying the
				// original request identifier replays to the same plan
				// rather than minting a new one. This is what makes
				// confirmation idempotent across restarts.
				if e.Idempotency != "" {
					s.idem[e.Idempotency] = &p
				}
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
