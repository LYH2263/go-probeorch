package probeorch

// Stats 返回运行统计。
func (o *Orch) Stats() Stats {
	o.mu.Lock()
	defer o.mu.Unlock()
	st := Stats{
		TotalRuns: o.totalRuns,
		OKRuns:    o.okRuns,
		FailRuns:  o.failRuns,
		LastID:    o.lastID,
		LastAt:    o.lastAt,
	}
	if o.reg != nil {
		list := o.reg.List()
		st.Registered = len(list)
		for _, t := range list {
			if t.Enabled {
				st.Enabled++
			}
		}
	}
	return st
}
