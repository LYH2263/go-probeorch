package probeorch

import "time"

// Tick 推进到期目标并依次 RunOnce。
func (o *Orch) Tick() (int, error) {
	return o.TickAt(o.now())
}

// TickAt 在指定时刻弹出到期 ID 并探测。
func (o *Orch) TickAt(at time.Time) (int, error) {
	o.mu.Lock()
	if err := o.checkOpenLocked(); err != nil {
		o.mu.Unlock()
		return 0, err
	}
	due := o.sched.PopDue(at)
	o.mu.Unlock()
	n := 0
	for _, id := range due {
		if _, err := o.RunOnce(id); err == nil {
			n++
		} else if err == ErrClosed {
			return n, err
		} else {
			n++ // 失败也算执行过
		}
	}
	return n, nil
}
