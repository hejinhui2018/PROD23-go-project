package service

import (
	"encoding/json"
	"fmt"
	"prod23/internal/domain"
	"strings"
	"time"
)

type Report struct {
	GeneratedAt          time.Time
	Faults, Plans, Steps int
	Events               int
	Notifications        int
	Status               map[string]int
}

func (s *Service) Report() Report {
	return Report{GeneratedAt: time.Now().UTC(), Faults: s.FaultCount(), Plans: s.PlanCount(), Steps: s.StepCount(), Events: len(s.Audit()), Notifications: s.NotificationCount(), Status: s.StatusCounts()}
}
func (r Report) Text() string {
	parts := []string{fmt.Sprintf("faults=%d", r.Faults), fmt.Sprintf("plans=%d", r.Plans), fmt.Sprintf("steps=%d", r.Steps), fmt.Sprintf("events=%d", r.Events), fmt.Sprintf("notifications=%d", r.Notifications)}
	for k, v := range r.Status {
		parts = append(parts, fmt.Sprintf("%s=%d", k, v))
	}
	return strings.Join(parts, " ")
}
func (r Report) Healthy() bool { return r.Events >= 0 && r.Faults >= 0 && r.Plans >= 0 && r.Steps >= 0 }
func (s *Service) EventTypes() []string {
	m := map[string]bool{}
	for _, e := range s.Audit() {
		m[e.Type] = true
	}
	r := []string{}
	for k := range m {
		r = append(r, k)
	}
	return r
}
func (s *Service) LastEventAt() time.Time {
	all := s.Audit()
	if len(all) == 0 {
		return time.Time{}
	}
	return all[len(all)-1].At
}
func (s *Service) SummaryLines() []string {
	r := s.Report()
	return []string{r.Text(), "healthy=" + fmt.Sprint(r.Healthy()), "generated_at=" + r.GeneratedAt.Format(time.RFC3339), "last_event=" + s.LastEventAt().Format(time.RFC3339)}
}
func (s *Service) HasPendingWork() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, st := range s.Steps {
		if !domain.IsTerminal(st.Status) {
			return true
		}
	}
	return false
}
func (s *Service) ReportJSON() ([]byte, error) { return json.Marshal(s.Report()) }
func (s *Service) ReportAt(at time.Time) Report {
	r := s.Report()
	if at.Before(r.GeneratedAt) {
		r.GeneratedAt = at
	}
	return r
}
func (s *Service) ReportHealthyText() string {
	if s.Report().Healthy() {
		return "healthy"
	}
	return "degraded"
}
func (s *Service) ActiveFaults() []domain.Fault {
	return s.ListFaults(FaultFilter{Status: string(domain.Assessed)})
}
func (s *Service) RestoredFaults() []domain.Fault {
	return s.ListFaults(FaultFilter{Status: string(domain.Restored)})
}
func (s *Service) FailedSteps() []domain.Step {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Step{}
	for _, st := range s.Steps {
		if st.Status == domain.Failed {
			out = append(out, *st)
		}
	}
	return out
}
func (s *Service) CompletedSteps() []domain.Step {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Step{}
	for _, st := range s.Steps {
		if st.Status == domain.Completed {
			out = append(out, *st)
		}
	}
	return out
}
func (s *Service) WorkerNames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seen := map[string]bool{}
	for _, st := range s.Steps {
		if st.ClaimedBy != "" {
			seen[st.ClaimedBy] = true
		}
	}
	out := []string{}
	for w := range seen {
		out = append(out, w)
	}
	return out
}
func (s *Service) PendingFaults() []domain.Fault {
	return s.ListFaults(FaultFilter{Status: string(domain.Assessed)})
}
func (s *Service) ReviewCount() int              { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.Reviews) }
func (s *Service) EventTypeCount(typ string) int { return s.EventCount(typ) }
func (s *Service) IsReady() bool                 { return s.Report().Healthy() && s.Store.Integrity() }

// Reports are built from locked state and immutable event snapshots.
// Operators can consume Text while APIs use the structured Report value.
// The report is intentionally side-effect free.
// Event counts remain monotonic across restarts.
// This keeps smoke output deterministic for operators.
// No report operation mutates recovery state.
