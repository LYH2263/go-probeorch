package probeorch

// Close 停止调度：先刷盘，再释放结果环与探针注册表。
func (o *Orch) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	// CLEAN: 先 Flush 再丢弃 rings / probers
	if o.persistPath != "" && o.rings != nil {
		_ = o.persistLocked()
	}
	o.closed = true
	if o.sched != nil {
		o.sched.Clear()
	}
	o.rings = nil
	o.probers = nil // 置空后 RunOnce 若未检查会 panic
	o.customs = nil
	return nil
}
