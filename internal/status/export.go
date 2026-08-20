package status

import (
	"encoding/json"
	"time"
)

// Report 导出用状态报告。
type Report struct {
	GeneratedAt time.Time         `json:"generated_at"`
	Targets     int               `json:"targets"`
	Healthy     int               `json:"healthy"`
	Rates       map[string]float64 `json:"rates"`
	Notes       []string          `json:"notes,omitempty"`
}

func (r Report) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Build 构造报告。
func Build(targets, healthy int, rates map[string]float64, notes ...string) Report {
	if rates == nil {
		rates = map[string]float64{}
	}
	return Report{
		GeneratedAt: time.Now().UTC(),
		Targets:     targets,
		Healthy:     healthy,
		Rates:       rates,
		Notes:       append([]string(nil), notes...),
	}
}
