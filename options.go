package probeorch

import (
	"net/http"
	"time"

	"example.com/probeorch/internal/clock"
)

// Option 配置 Orch。
type Option func(*Orch)

// WithClock 注入可测时钟。
func WithClock(c clock.Clock) Option {
	return func(o *Orch) {
		if c != nil {
			o.clk = c
		}
	}
}

// WithPersistPath 设置目标+结果快照路径。
func WithPersistPath(path string) Option {
	return func(o *Orch) { o.persistPath = path }
}

// WithHistorySize 环形历史容量（每目标）。
func WithHistorySize(n int) Option {
	return func(o *Orch) {
		if n > 0 {
			o.historySize = n
		}
	}
}

// WithHTTPClient 注入 HTTP 客户端。
func WithHTTPClient(c *http.Client) Option {
	return func(o *Orch) {
		if c != nil {
			o.httpClient = c
		}
	}
}

// WithProbeWait 注入探测前等待。
func WithProbeWait(w ProbeWaiter) Option {
	return func(o *Orch) {
		if w != nil {
			o.wait = w
		}
	}
}

// WithDefaultTimeout 默认单次超时。
func WithDefaultTimeout(d time.Duration) Option {
	return func(o *Orch) {
		if d > 0 {
			o.defaultTimeout = d
		}
	}
}

// WithDefaultInterval 默认调度间隔。
func WithDefaultInterval(d time.Duration) Option {
	return func(o *Orch) {
		if d > 0 {
			o.defaultInterval = d
		}
	}
}

// WithMaxTargets 最大目标数。
func WithMaxTargets(n int) Option {
	return func(o *Orch) {
		if n > 0 {
			o.maxTargets = n
		}
	}
}

// WithProbeWaitDelay 默认探测前延时（配合 sleepWaiter）。
func WithProbeWaitDelay(d time.Duration) Option {
	return func(o *Orch) {
		if d >= 0 {
			o.probeWaitDelay = d
		}
	}
}
