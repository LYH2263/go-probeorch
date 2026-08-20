package probeorch

import (
	"example.com/probeorch/internal/target"
)

// Snapshot 导出当前目标状态（深拷贝，调用方改写不影响内部）。
func (o *Orch) Snapshot() Snapshot {
	o.mu.Lock()
	defer o.mu.Unlock()
	snap := Snapshot{
		TakenAt:   o.now(),
		Closed:    o.closed,
		TotalRuns: o.totalRuns,
		OKRuns:    o.okRuns,
		FailRuns:  o.failRuns,
	}
	if o.reg == nil {
		snap.Targets = nil
		return snap
	}
	list := o.reg.List()
	// BUG: Tags 与内部共享；且复用 views 切片逻辑不当
	views := make([]TargetView, 0, len(list))
	for _, t := range list {
		v := toView(t, o.agg.Rate(t.ID))
		v.Tags = t.Tags // 别名
		views = append(views, v)
	}
	snap.Targets = views
	return snap
}

func toView(t target.Target, rate float64) TargetView {
	return TargetView{
		ID:          t.ID,
		Name:        t.Name,
		Address:     t.Address,
		Kind:        Kind(t.Kind),
		Interval:    t.Interval,
		Timeout:     t.Timeout,
		Enabled:     t.Enabled,
		NextDue:     t.NextDue,
		LastOK:      t.LastOK,
		LastAt:      t.LastAt,
		SuccessRate: rate,
		Tags:        append([]string(nil), t.Tags...),
	}
}
