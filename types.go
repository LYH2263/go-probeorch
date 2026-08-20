package probeorch

import (
	"context"
	"time"
)

// Kind 探针类型。
type Kind string

const (
	KindHTTP   Kind = "http"
	KindTCP    Kind = "tcp"
	KindCustom Kind = "custom"
)

// Spec 注册探测目标的规格。
type Spec struct {
	ID       string
	Name     string
	Address  string
	Kind     Kind
	Interval time.Duration
	Timeout  time.Duration
	Method   string
	Expect   int
	Headers  map[string]string
	Enabled  bool
	Tags     []string
	Detail   []byte
}

// ProbeFunc 自定义探针。
type ProbeFunc func(ctx context.Context, address string, timeout time.Duration) (ok bool, detail []byte, err error)

// Result 单次探测结果（对外视图）。
type Result struct {
	TargetID string
	OK       bool
	At       time.Time
	Latency  time.Duration
	Status   int
	Message  string
	Detail   []byte
	ErrText  string
}

// TargetView 快照中的目标视图。
type TargetView struct {
	ID          string
	Name        string
	Address     string
	Kind        Kind
	Interval    time.Duration
	Timeout     time.Duration
	Enabled     bool
	NextDue     time.Time
	LastOK      bool
	LastAt      time.Time
	SuccessRate float64
	Tags        []string
}

// Snapshot 编排器状态快照。
type Snapshot struct {
	Targets   []TargetView
	TakenAt   time.Time
	Closed    bool
	TotalRuns uint64
	OKRuns    uint64
	FailRuns  uint64
}

// Stats 运行计数。
type Stats struct {
	Registered int
	Enabled    int
	TotalRuns  uint64
	OKRuns     uint64
	FailRuns   uint64
	LastID     string
	LastAt     time.Time
}
