package metrics

import "sync/atomic"

type Counter struct{ n atomic.Int64 }

func (c *Counter) Inc()         { c.n.Add(1) }
func (c *Counter) Value() int64 { return c.n.Load() }
