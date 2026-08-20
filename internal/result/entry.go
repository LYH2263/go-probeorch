package result

import "time"

// Entry 环形缓冲中的一条结果。
type Entry struct {
	TargetID string
	OK       bool
	At       time.Time
	Latency  time.Duration
	Status   int
	Message  string
	Detail   []byte
	ErrText  string
}
