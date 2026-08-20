package schedule

import "time"

// SelectDue 从候选 map 中选出已到期 ID（非堆路径，供测试/导出）。
func SelectDue(next map[string]time.Time, at time.Time) []string {
	var out []string
	for id, due := range next {
		if !due.After(at) {
			out = append(out, id)
		}
	}
	return out
}

// Soonest 返回最近到期。
func Soonest(next map[string]time.Time) (string, time.Time, bool) {
	var (
		bestID string
		best   time.Time
		ok     bool
	)
	for id, due := range next {
		if !ok || due.Before(best) {
			bestID, best, ok = id, due, true
		}
	}
	return bestID, best, ok
}
