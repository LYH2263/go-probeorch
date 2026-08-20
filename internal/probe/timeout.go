package probe

import (
	"context"
	"time"
)

// WithTimeout 派生带超时的 context。
func WithTimeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		d = 3 * time.Second
	}
	return context.WithTimeout(parent, d)
}

// Remaining 剩余超时。
func Remaining(ctx context.Context) time.Duration {
	dl, ok := ctx.Deadline()
	if !ok {
		return 0
	}
	return time.Until(dl)
}
