package probeorch

// ListIDs 返回已注册目标 ID 列表（拷贝）。
func (o *Orch) ListIDs() []string {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.reg == nil {
		return nil
	}
	list := o.reg.List()
	out := make([]string, 0, len(list))
	for _, t := range list {
		out = append(out, t.ID)
	}
	return out
}

// GetTarget 返回单个目标视图。
func (o *Orch) GetTarget(id string) (TargetView, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return TargetView{}, err
	}
	t, ok := o.reg.Get(id)
	if !ok {
		return TargetView{}, ErrNotFound
	}
	return toView(t, o.agg.Rate(id)), nil
}
