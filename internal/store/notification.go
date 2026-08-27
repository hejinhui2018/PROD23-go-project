package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Notification struct {
	ID, Type, Payload string
	Sent              bool
	Attempts          int
	LastError         string
	CreatedAt         time.Time
	LastAttemptAt     time.Time
}
type Queue struct {
	mu    sync.Mutex
	items []Notification
	dir   string
}

func LoadQueue(dir string) *Queue {
	q := &Queue{dir: dir}
	b, e := os.ReadFile(filepath.Join(dir, "notifications.json"))
	if e == nil {
		_ = json.Unmarshal(b, &q.items)
	}
	return q
}
func (q *Queue) persist() {
	if q.dir == "" {
		return
	}
	b, _ := json.MarshalIndent(q.items, "", "  ")
	_ = os.WriteFile(filepath.Join(q.dir, "notifications.json"), b, 0644)
}
func (q *Queue) Enqueue(n Notification) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, x := range q.items {
		if x.ID == n.ID {
			return
		}
	}
	q.items = append(q.items, n)
	q.persist()
}
func (q *Queue) Pending() []Notification {
	q.mu.Lock()
	defer q.mu.Unlock()
	r := []Notification{}
	for _, n := range q.items {
		if !n.Sent {
			r = append(r, n)
		}
	}
	return r
}

// BeginDelivery records a durable delivery attempt. The message stays pending so
// that a downstream failure leaves it retryable; Ack marks it delivered and
// removing it from the pending set. LastAttemptAt timestamps the attempt so the
// dispatcher can apply backoff before retrying without dropping the message.
func (q *Queue) BeginDelivery(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.items {
		if q.items[i].ID != id || q.items[i].Sent {
			continue
		}
		q.items[i].Attempts++
		q.items[i].LastAttemptAt = time.Now().UTC()
		q.persist()
		return true
	}
	return false
}

// Nack returns a message to the retryable pending set after a failed delivery.
// The attempt was already counted by BeginDelivery (which also set LastAttemptAt),
// so the dispatcher applies backoff via RetryDelay(Attempts) before retrying.
// LastError records the most recent failure so operators can inspect it without
// losing the message; it stays retryable across restarts via notifications.json.
func (q *Queue) Nack(id, cause string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.items {
		if q.items[i].ID != id {
			continue
		}
		q.items[i].Sent = false
		q.items[i].LastError = cause
		q.persist()
		return
	}
}
func (q *Queue) Ack(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.items {
		if q.items[i].ID == id {
			q.items[i].Sent = true
			if q.items[i].Attempts == 0 {
				q.items[i].Attempts++
			}
			q.persist()
		}
	}
}
func (q *Queue) All() []Notification {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]Notification(nil), q.items...)
}
