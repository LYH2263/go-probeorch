package target

import "errors"

var (
	ErrFull     = errors.New("target registry full")
	ErrNotFound = errors.New("target not found")
	ErrExists   = errors.New("target exists")
)

// Registry 目标注册表。
type Registry struct {
	max int
	m   map[string]Target
	ord []string
}

func NewRegistry(max int) *Registry {
	if max < 1 {
		max = 1
	}
	return &Registry{max: max, m: make(map[string]Target)}
}

func (r *Registry) Has(id string) bool {
	_, ok := r.m[id]
	return ok
}

func (r *Registry) Get(id string) (Target, bool) {
	t, ok := r.m[id]
	if !ok {
		return Target{}, false
	}
	return Clone(t), true
}

func (r *Registry) Add(t Target) error {
	if !Valid(t) {
		return errors.New("invalid target")
	}
	if _, ok := r.m[t.ID]; ok {
		return ErrExists
	}
	if len(r.m) >= r.max {
		return ErrFull
	}
	r.m[t.ID] = Clone(t)
	r.ord = append(r.ord, t.ID)
	return nil
}

func (r *Registry) Update(t Target) error {
	if _, ok := r.m[t.ID]; !ok {
		return ErrNotFound
	}
	r.m[t.ID] = Clone(t)
	return nil
}

func (r *Registry) Remove(id string) error {
	if _, ok := r.m[id]; !ok {
		return ErrNotFound
	}
	delete(r.m, id)
	for i, x := range r.ord {
		if x == id {
			r.ord = append(r.ord[:i], r.ord[i+1:]...)
			break
		}
	}
	return nil
}

func (r *Registry) List() []Target {
	out := make([]Target, 0, len(r.ord))
	for _, id := range r.ord {
		if t, ok := r.m[id]; ok {
			out = append(out, Clone(t))
		}
	}
	return out
}

func (r *Registry) Len() int { return len(r.m) }
