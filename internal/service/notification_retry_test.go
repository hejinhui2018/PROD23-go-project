package service

import (
	"context"
	"testing"
	"time"

	"prod23/internal/store"
)

func TestNotificationRetryAfterDownstreamError(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(st)
	svc.Queue.Enqueue(store.Notification{ID: "restoration-1", Type: "restoration.completed", Payload: "fault-1", CreatedAt: time.Now().UTC()})

	report := svc.Dispatch(context.Background(), &MemoryNotifier{Fail: true}, 1)
	if report.Failed != 1 || report.Delivered != 0 {
		t.Fatalf("unexpected dispatch report: %+v", report)
	}
	if got := len(svc.Queue.Pending()); got != 1 {
		t.Fatalf("notification was removed from pending queue: got %d", got)
	}

	reopened, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(New(reopened).Queue.Pending()); got != 1 {
		t.Fatalf("notification did not survive restart: got %d", got)
	}
}
