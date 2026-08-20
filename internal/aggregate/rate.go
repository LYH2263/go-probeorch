package aggregate

// Book 成功率账本。
type Book struct {
	m map[string]*counter
}

type counter struct {
	ok    uint64
	total uint64
}

func NewBook() *Book { return &Book{m: make(map[string]*counter)} }

func (b *Book) Ensure(id string) {
	if _, ok := b.m[id]; !ok {
		b.m[id] = &counter{}
	}
}

func (b *Book) Observe(id string, ok bool) {
	b.Ensure(id)
	c := b.m[id]
	c.total++
	if ok {
		c.ok++
	}
}

func (b *Book) Rate(id string) float64 {
	c, ok := b.m[id]
	if !ok || c.total == 0 {
		return 0
	}
	return float64(c.ok) / float64(c.total)
}

func (b *Book) Drop(id string) { delete(b.m, id) }

func (b *Book) Reset(id string) {
	b.m[id] = &counter{}
}
