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
	now := time.Now().UTC()
	for _, msg := range s.Queue.Pending() {
		if r.Attempted >= limit {
			break
		}
		// Backoff: a message that has failed before must wait RetryDelay before
		// being retried, so a flapping downstream does not get hammered. The
		// message stays pending and is simply skipped on this pass until ready,
		// keeping the queue non-blocking and leaving room for other messages.
		if ready := readyAt(msg); ready.After(now) {
			continue
		}
		r.Attempted++
		if !s.Queue.BeginDelivery(msg.ID) {
			r.Failed++
			r.Errors = append(r.Errors, "notification no longer pending")
			continue
		}
		// BeginDelivery only records the attempt, so a failed Send is Nacked
		// back to the pending set for the next dispatch pass. This keeps the
		// message retryable across restarts: notifications.json stores
		// Sent=false, Attempts and LastAttemptAt, so a crashed process resumes
		// the retry instead of dropping the notification as if delivered.
		if err := n.Send(ctx, msg); err != nil {
			s.Queue.Nack(msg.ID, err.Error())
			r.Failed++
			r.Errors = append(r.Errors, err.Error())
			continue
		}
		s.Queue.Ack(msg.ID)
		r.Delivered++
	}
	return r
}

// readyAt reports the earliest time a message may be retried after a failure,
// applying exponential backoff via RetryDelay. A message that has never been
// attempted (or whose attempt timestamp was lost) is immediately ready.
func readyAt(n store.Notification) time.Time {
	if n.Attempts == 0 || n.LastAttemptAt.IsZero() {
		return time.Time{}
	}
	return n.LastAttemptAt.Add(RetryDelay(n.Attempts))
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
