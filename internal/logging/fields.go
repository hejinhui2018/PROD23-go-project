package logging

import (
	"context"
	"fmt"
	"sync"
)

type ctxKey struct{}

func WithRequest(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}
func RequestID(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

type Buffer struct {
	mu    sync.Mutex
	Lines []string
}

func (b *Buffer) Add(line string)              { b.mu.Lock(); defer b.mu.Unlock(); b.Lines = append(b.Lines, line) }
func (b *Buffer) Len() int                     { b.mu.Lock(); defer b.mu.Unlock(); return len(b.Lines) }
func FieldError(name string, err error) string { return fmt.Sprintf("%s: %v", name, err) }
