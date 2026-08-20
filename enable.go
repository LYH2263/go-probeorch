package probeorch

// Enable 启用目标并加入调度。
func (o *Orch) Enable(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return err
	}
	t, ok := o.reg.Get(id)
	if !ok {
		return ErrNotFound
	}
	t.Enabled = true
	t.NextDue = o.now().Add(t.Interval)
	if err := o.reg.Update(t); err != nil {
		return err
	}
	o.sched.Upsert(id, t.NextDue)
	return nil
}

// Disable 禁用目标并移出调度。
func (o *Orch) Disable(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return err
	}
	t, ok := o.reg.Get(id)
	if !ok {
		return ErrNotFound
	}
	t.Enabled = false
	if err := o.reg.Update(t); err != nil {
		return err
	}
	o.sched.Remove(id)
	return nil
}
