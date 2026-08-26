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
	InFlight          bool
	Attempts          int
	CreatedAt         time.Time
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
		if !n.Sent && !n.InFlight {
			r = append(r, n)
		}
	}
	return r
}

// Claim marks a pending notification as being processed by a dispatcher.
func (q *Queue) Claim(id string) (Notification, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.items {
		if q.items[i].ID != id || q.items[i].Sent || q.items[i].InFlight {
			continue
		}
		q.items[i].InFlight = true
		q.items[i].Attempts++
		q.persist()
		return q.items[i], true
	}
	return Notification{}, false
}

// Release makes a claimed notification available for a later retry.
func (q *Queue) Release(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.items {
		if q.items[i].ID != id || q.items[i].Sent || !q.items[i].InFlight {
			continue
		}
		q.items[i].InFlight = false
		q.persist()
		return true
	}
	return false
}

func (q *Queue) Ack(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.items {
		if q.items[i].ID == id {
			q.items[i].Sent = true
			q.items[i].InFlight = false
			q.persist()
		}
	}
}
func (q *Queue) All() []Notification {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]Notification(nil), q.items...)
}
