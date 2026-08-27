package service

import (
	"prod23/internal/domain"
	"time"
)

type Lifecycle struct {
	StartedAt time.Time
	StoppedAt time.Time
	Running   bool
}

func (l *Lifecycle) Start() {
	l.StartedAt = time.Now().UTC()
	l.Running = true
	l.StoppedAt = time.Time{}
}
func (l *Lifecycle) Stop() { l.StoppedAt = time.Now().UTC(); l.Running = false }
func (l Lifecycle) Uptime() time.Duration {
	if l.StartedAt.IsZero() {
		return 0
	}
	end := l.StoppedAt
	if end.IsZero() {
		end = time.Now().UTC()
	}
	return end.Sub(l.StartedAt)
}
func (s *Service) FaultStatus(id string) (domain.FaultStatus, error) {
	f, e := s.GetFault(id)
	return f.Status, e
}
