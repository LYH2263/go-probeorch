package aggregate

import "time"

// Sample 时间窗样本。
type Sample struct {
	At time.Time
	OK bool
}

// Window 滑动时间窗成功率。
type Window struct {
	keep time.Duration
	buf  []Sample
}

func NewWindow(keep time.Duration) *Window {
	if keep <= 0 {
		keep = time.Hour
	}
	return &Window{keep: keep}
}

func (w *Window) Add(at time.Time, ok bool) {
	w.buf = append(w.buf, Sample{At: at, OK: ok})
	w.gc(at)
}

func (w *Window) gc(now time.Time) {
	cut := now.Add(-w.keep)
	i := 0
	for i < len(w.buf) && w.buf[i].At.Before(cut) {
		i++
	}
	if i > 0 {
		w.buf = append([]Sample(nil), w.buf[i:]...)
	}
}

func (w *Window) Rate(now time.Time) float64 {
	w.gc(now)
	if len(w.buf) == 0 {
		return 0
	}
	var okn int
	for _, s := range w.buf {
		if s.OK {
			okn++
		}
	}
	return float64(okn) / float64(len(w.buf))
}

func (w *Window) Len() int { return len(w.buf) }
