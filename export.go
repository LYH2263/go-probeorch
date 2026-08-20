package probeorch

import (
	"example.com/probeorch/internal/status"
)

// ExportStatus 导出状态报告 JSON。
func (o *Orch) ExportStatus() ([]byte, error) {
	snap := o.Snapshot()
	rates := make(map[string]float64, len(snap.Targets))
	healthy := 0
	for _, t := range snap.Targets {
		rates[t.ID] = t.SuccessRate
		if t.LastOK {
			healthy++
		}
	}
	rep := status.Build(len(snap.Targets), healthy, rates, "probeorch export")
	return rep.JSON()
}
