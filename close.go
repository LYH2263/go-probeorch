package probeorch

// Close 停止调度：先刷盘，再释放结果环与探针注册表。
func (o *Orch) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	// BUG: 先丢结果环再 Flush，刷出空历史
	o.closed = true
	if o.sched != nil {
		o.sched.Clear()
	}
	o.rings = nil
	if o.persistPath != "" {
		_ = o.persistLocked()
	}
	o.probers = nil
	o.customs = nil
	return nil
}
