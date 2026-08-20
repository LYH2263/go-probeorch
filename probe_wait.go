package probeorch

import (
	"context"
	"time"
)

// ProbeWaiter 探测前可中断等待（测试注入长延时）。
type ProbeWaiter interface {
	Wait(ctx context.Context, d time.Duration) error
}

type sleepWaiter struct{}

// Wait 监听 ctx，可被取消打断。
func (sleepWaiter) Wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
