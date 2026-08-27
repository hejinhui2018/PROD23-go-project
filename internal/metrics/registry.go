package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
	gauges   map[string]*atomic.Int64
}

func NewRegistry() *Registry {
	return &Registry{counters: map[string]*Counter{}, gauges: map[string]*atomic.Int64{}}
}
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c := r.counters[name]; c != nil {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}
func (r *Registry) Gauge(name string) *atomic.Int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g := r.gauges[name]; g != nil {
		return g
	}
	g := &atomic.Int64{}
	r.gauges[name] = g
	return g
}
func (r *Registry) Prometheus() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := ""
	for n, c := range r.counters {
		out += fmt.Sprintf("%s_total %d\\n", n, c.Value())
	}
	for n, g := range r.gauges {
		out += fmt.Sprintf("%s %d\\n", n, g.Load())
	}
	return out
}
