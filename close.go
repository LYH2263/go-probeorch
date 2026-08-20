package probeorch

// Close 停止调度：先刷盘（此时结果环仍持有历史），再释放结果环与探针注册表。
func (o *Orch) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	o.closed = true
	if o.sched != nil {
		o.sched.Clear()
	}
	// 先刷盘：结果环尚未释放，历史可正确写入快照。
	if o.persistPath != "" {
		_ = o.persistLocked()
	}
	// 刷盘完成后再丢弃结果环与探针注册表。
	o.rings = nil
	o.probers = nil
	o.customs = nil
	return nil
}
