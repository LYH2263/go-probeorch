package schedule

import (
	"container/heap"
	"time"
)

type item struct {
	id   string
	due  time.Time
	idx  int
}

type dueHeap []*item

func (h dueHeap) Len() int           { return len(h) }
func (h dueHeap) Less(i, j int) bool { return h[i].due.Before(h[j].due) }
func (h dueHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}
func (h *dueHeap) Push(x any) {
	it := x.(*item)
	it.idx = len(*h)
	*h = append(*h, it)
}
func (h *dueHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	it.idx = -1
	return it
}

// Table 到期调度表。
type Table struct {
	h    dueHeap
	byID map[string]*item
}

func NewTable() *Table {
	t := &Table{byID: make(map[string]*item)}
	heap.Init(&t.h)
	return t
}

func (t *Table) Upsert(id string, due time.Time) {
	if it, ok := t.byID[id]; ok {
		it.due = due
		heap.Fix(&t.h, it.idx)
		return
	}
	it := &item{id: id, due: due}
	heap.Push(&t.h, it)
	t.byID[id] = it
}

func (t *Table) Remove(id string) {
	it, ok := t.byID[id]
	if !ok {
		return
	}
	heap.Remove(&t.h, it.idx)
	delete(t.byID, id)
}

func (t *Table) PopDue(at time.Time) []string {
	var out []string
	for t.h.Len() > 0 {
		it := t.h[0]
		if it.due.After(at) {
			break
		}
		heap.Pop(&t.h)
		delete(t.byID, it.id)
		out = append(out, it.id)
	}
	return out
}

func (t *Table) Peek() (string, time.Time, bool) {
	if t.h.Len() == 0 {
		return "", time.Time{}, false
	}
	it := t.h[0]
	return it.id, it.due, true
}

func (t *Table) Len() int { return t.h.Len() }

func (t *Table) Clear() {
	t.h = nil
	t.byID = make(map[string]*item)
	heap.Init(&t.h)
}

// NextDue 查询某 ID 的到期时间。
func (t *Table) NextDue(id string) (time.Time, bool) {
	it, ok := t.byID[id]
	if !ok {
		return time.Time{}, false
	}
	return it.due, true
}
