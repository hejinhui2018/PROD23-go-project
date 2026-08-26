package metrics

import (
	"sync/atomic"
	"time"
)

type Timer struct {
	started time.Time
	sum     atomic.Int64
	count   atomic.Int64
}

func StartTimer() *Timer { return &Timer{started: time.Now()} }
func (t *Timer) Stop() {
	if t.started.IsZero() {
		return
	}
	t.sum.Add(time.Since(t.started).Microseconds())
	t.count.Add(1)
}
func (t *Timer) Count() int64       { return t.count.Load() }
func (t *Timer) TotalMicros() int64 { return t.sum.Load() }
func (t *Timer) AverageMicros() int64 {
	n := t.Count()
	if n == 0 {
		return 0
	}
	return t.TotalMicros() / n
}
