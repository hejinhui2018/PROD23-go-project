package service

import (
	"context"
	"errors"
	"prod23/internal/store"
	"time"
)

type Notifier interface {
	Send(context.Context, store.Notification) error
}
type DispatchReport struct {
	Attempted, Delivered, Failed int
	Errors                       []string
}

func (s *Service) Dispatch(ctx context.Context, n Notifier, limit int) DispatchReport {
	if limit <= 0 {
		limit = 100
	}
	r := DispatchReport{}
	for _, msg := range s.Queue.Pending() {
		if r.Attempted >= limit {
			break
		}
		claimed, ok := s.Queue.Claim(msg.ID)
		if !ok {
			continue
		}
		r.Attempted++
		// Claim holds the lease (InFlight) for the duration of the send;
		// releasing before the send would re-expose the message to a
		// concurrent dispatcher and produce duplicate deliveries.
		if err := n.Send(ctx, claimed); err != nil {
			r.Failed++
			r.Errors = append(r.Errors, err.Error())
			s.Queue.Release(claimed.ID)
			continue
		}
		s.Queue.Ack(claimed.ID)
		r.Delivered++
	}
	return r
}

type MemoryNotifier struct {
	Messages []store.Notification
	Fail     bool
}

func (m *MemoryNotifier) Send(ctx context.Context, n store.Notification) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if m.Fail {
		return errors.New("downstream unavailable")
	}
	m.Messages = append(m.Messages, n)
	return nil
}
func (s *Service) NotificationCount() int { return len(s.Queue.All()) }
func RetryDelay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 6 {
		attempt = 6
	}
	return time.Duration(1<<attempt) * time.Second
}
