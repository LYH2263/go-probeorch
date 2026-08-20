package probeorch

// History 返回目标探测历史（新→旧），Detail 已深拷贝。
func (o *Orch) History(id string) ([]Result, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return nil, err
	}
	if !o.reg.Has(id) {
		return nil, ErrNotFound
	}
	entries := o.rings.List(id)
	out := make([]Result, 0, len(entries))
	for _, e := range entries {
		out = append(out, toPublicResult(e))
	}
	return out, nil
}

// SuccessRate 返回目标成功率 [0,1]。
func (o *Orch) SuccessRate(id string) (float64, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return 0, err
	}
	if !o.reg.Has(id) {
		return 0, ErrNotFound
	}
	return o.agg.Rate(id), nil
}
