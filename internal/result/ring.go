package result

// Ring 固定容量环形缓冲（新在前）。
type Ring struct {
	cap  int
	buf  []Entry
	head int
	len  int
}

func NewRing(capacity int) *Ring {
	if capacity < 1 {
		capacity = 1
	}
	return &Ring{cap: capacity, buf: make([]Entry, capacity)}
}

func (r *Ring) Push(e Entry) {
	r.buf[r.head] = CloneEntry(e)
	r.head = (r.head + 1) % r.cap
	if r.len < r.cap {
		r.len++
	}
}

func (r *Ring) List() []Entry {
	out := make([]Entry, 0, r.len)
	for i := 0; i < r.len; i++ {
		idx := (r.head - 1 - i + r.cap*2) % r.cap
		out = append(out, CloneEntry(r.buf[idx]))
	}
	return out
}

func (r *Ring) Len() int { return r.len }
