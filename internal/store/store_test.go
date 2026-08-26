package store

import (
	"prod23/internal/domain"
	"testing"
	"time"
)

func TestAppendReload(t *testing.T) {
	d := t.TempDir()
	s, _ := Open(d)
	s.Append(domain.Event{ID: "1", Type: "x", At: time.Now()})
	r, _ := Open(d)
	if len(r.Events()) != 1 {
		t.Fatal("reload")
	}
}
