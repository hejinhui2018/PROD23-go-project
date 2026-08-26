package service

import (
	"context"
	"prod23/internal/domain"
	"sync"
	"time"
)

type ClaimCoordinator struct {
	mu     sync.Mutex
	owners map[string]string
	leases map[string]time.Time
}

func NewClaimCoordinator() *ClaimCoordinator {
	return &ClaimCoordinator{owners: map[string]string{}, leases: map[string]time.Time{}}
}
func (c *ClaimCoordinator) Acquire(ctx context.Context, step, worker string, ttl time.Duration) error {
	for {
		c.mu.Lock()
		now := time.Now()
		if owner, ok := c.owners[step]; !ok || owner == worker || now.After(c.leases[step]) {
			c.owners[step] = worker
			c.leases[step] = now.Add(ttl)
			c.mu.Unlock()
			return nil
		}
		c.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}
	}
}
func (c *ClaimCoordinator) Release(step, worker string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.owners[step] != worker {
		return false
	}
	delete(c.owners, step)
	delete(c.leases, step)
	return true
}
func (c *ClaimCoordinator) Owner(step string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.owners[step]
	return v, ok
}
func (s *Service) ValidateOrdering(planID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.Plans[planID]
	if p == nil {
		return domain.ErrNotFound
	}
	last := 0
	for _, id := range p.Steps {
		if s.Steps[id].Sequence <= last {
			return domain.ErrInvalid
		}
		last = s.Steps[id].Sequence
	}
	return nil
}
