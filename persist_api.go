package probeorch

import (
	"example.com/probeorch/internal/persist"
	"example.com/probeorch/internal/result"
	"example.com/probeorch/internal/target"
)

// LoadPersist 从快照恢复目标与历史。
func (o *Orch) LoadPersist() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return err
	}
	if o.persistPath == "" {
		return nil
	}
	snap, err := persist.Load(o.persistPath)
	if err != nil {
		return err
	}
	for _, t := range snap.Targets {
		_ = o.reg.Add(t)
		if t.Enabled {
			o.sched.Upsert(t.ID, t.NextDue)
		}
		o.rings.Ensure(t.ID)
		o.agg.Ensure(t.ID)
	}
	for id, entries := range snap.History {
		for _, e := range entries {
			o.rings.Push(id, e)
			o.agg.Observe(id, e.OK)
		}
	}
	o.totalRuns = snap.TotalRuns
	o.okRuns = snap.OKRuns
	o.failRuns = snap.FailRuns
	return nil
}

// Flush 将当前状态写入持久化路径。
func (o *Orch) Flush() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return err
	}
	return o.persistLocked()
}

func (o *Orch) persistLocked() error {
	if o.persistPath == "" {
		return nil
	}
	var targets []target.Target
	if o.reg != nil {
		targets = o.reg.List()
	}
	hist := map[string][]result.Entry{}
	if o.rings != nil {
		for _, t := range targets {
			hist[t.ID] = o.rings.List(t.ID)
		}
	}
	snap := persist.Snapshot{
		Targets:   targets,
		History:   hist,
		TotalRuns: o.totalRuns,
		OKRuns:    o.okRuns,
		FailRuns:  o.failRuns,
		SavedAt:   o.now(),
	}
	return persist.Save(o.persistPath, snap)
}
