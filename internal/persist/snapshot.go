package persist

import (
	"time"

	"example.com/probeorch/internal/result"
	"example.com/probeorch/internal/target"
)

// Snapshot 落盘结构。
type Snapshot struct {
	Targets   []target.Target            `json:"targets"`
	History   map[string][]result.Entry  `json:"history"`
	TotalRuns uint64                     `json:"total_runs"`
	OKRuns    uint64                     `json:"ok_runs"`
	FailRuns  uint64                     `json:"fail_runs"`
	SavedAt   time.Time                  `json:"saved_at"`
}
