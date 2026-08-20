package probeorch

import (
	"net/http"
	"sync"
	"time"

	"example.com/probeorch/internal/aggregate"
	"example.com/probeorch/internal/clock"
	"example.com/probeorch/internal/probe"
	"example.com/probeorch/internal/result"
	"example.com/probeorch/internal/schedule"
	"example.com/probeorch/internal/target"
)

const (
	defaultHistorySize  = 64
	defaultTimeout      = 3 * time.Second
	defaultInterval     = 30 * time.Second
	defaultMaxTargets   = 1024
)

// Orch 健康探测编排门面。零值不可用，须 New。
type Orch struct {
	mu sync.Mutex

	closed  bool
	clk     clock.Clock
	reg     *target.Registry
	sched   *schedule.Table
	rings   *result.Store
	agg     *aggregate.Book
	probers *probe.Registry

	httpClient      *http.Client
	wait            ProbeWaiter
	persistPath     string
	historySize     int
	defaultTimeout  time.Duration
	defaultInterval time.Duration
	maxTargets      int
	probeWaitDelay  time.Duration

	customs map[string]ProbeFunc

	totalRuns uint64
	okRuns    uint64
	failRuns  uint64
	lastID    string
	lastAt    time.Time
}

// New 构造 Orch。
func New(opts ...Option) *Orch {
	o := &Orch{
		clk:             clock.Real{},
		historySize:     defaultHistorySize,
		defaultTimeout:  defaultTimeout,
		defaultInterval: defaultInterval,
		maxTargets:      defaultMaxTargets,
		customs:         make(map[string]ProbeFunc),
		wait:            sleepWaiter{},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	if o.clk == nil {
		o.clk = clock.Real{}
	}
	if o.wait == nil {
		o.wait = sleepWaiter{}
	}
	if o.historySize < 4 {
		o.historySize = 4
	}
	if o.maxTargets < 8 {
		o.maxTargets = 8
	}
	if o.httpClient == nil {
		o.httpClient = &http.Client{Timeout: o.defaultTimeout}
	}
	o.reg = target.NewRegistry(o.maxTargets)
	o.sched = schedule.NewTable()
	o.rings = result.NewStore(o.historySize)
	o.agg = aggregate.NewBook()
	o.probers = probe.NewRegistry(o.httpClient)
	return o
}

func (o *Orch) now() time.Time { return o.clk.Now() }

func (o *Orch) checkOpenLocked() error {
	if o.closed {
		return ErrClosed
	}
	return nil
}
