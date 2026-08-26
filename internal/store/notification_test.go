package store

import (
	"testing"
	"time"
)

func TestQueueReplayAndAck(t *testing.T) {
	d := t.TempDir()
	q := LoadQueue(d)
	q.Enqueue(Notification{ID: "n", Type: "restoration.completed", CreatedAt: time.Now()})
	q.Ack("n")
	q2 := LoadQueue(d)
	if len(q2.Pending()) != 0 || len(q2.All()) != 1 {
		t.Fatal("queue persistence")
	}
}
