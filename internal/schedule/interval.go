package schedule

import "time"

// NextFrom 按固定间隔计算下一次到期。
func NextFrom(last time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		interval = time.Second
	}
	if last.IsZero() {
		return time.Now().Add(interval)
	}
	return last.Add(interval)
}

// Align 将时刻对齐到间隔网格（简化版）。
func Align(at time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		return at
	}
	ns := at.UnixNano()
	step := int64(interval)
	aligned := ns - (ns % step)
	return time.Unix(0, aligned).UTC()
}

// DueWithin 判断 due 是否在窗口内。
func DueWithin(due, now time.Time, window time.Duration) bool {
	if due.After(now) {
		return false
	}
	if window <= 0 {
		return true
	}
	return !due.Before(now.Add(-window))
}
