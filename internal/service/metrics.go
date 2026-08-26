package service

import (
	"sync/atomic"
	"time"
)

type ServiceMetrics struct {
	Faults, Plans, Claims, Acks, Completions, Failures, Conflicts atomic.Int64
	Started                                                       time.Time
}

func (m *ServiceMetrics) Snapshot() map[string]int64 {
	return map[string]int64{"faults": m.Faults.Load(), "plans": m.Plans.Load(), "claims": m.Claims.Load(), "acks": m.Acks.Load(), "completions": m.Completions.Load(), "failures": m.Failures.Load(), "conflicts": m.Conflicts.Load()}
}
func (m *ServiceMetrics) Uptime() time.Duration {
	if m.Started.IsZero() {
		return 0
	}
	return time.Since(m.Started)
}
func (m *ServiceMetrics) Reset() {
	m.Faults.Store(0)
	m.Plans.Store(0)
	m.Claims.Store(0)
	m.Acks.Store(0)
	m.Completions.Store(0)
	m.Failures.Store(0)
	m.Conflicts.Store(0)
}
func MetricNames() []string {
	return []string{"faults", "plans", "claims", "acks", "completions", "failures", "conflicts"}
}
