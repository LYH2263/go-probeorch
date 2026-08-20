package probeorch

import (
	"fmt"

	"example.com/probeorch/internal/idgen"
	"example.com/probeorch/internal/target"
	"example.com/probeorch/internal/validate"
)

// Register 注册探测目标；持久化失败时不得留在调度表。
func (o *Orch) Register(spec Spec) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return "", err
	}
	if err := validate.Spec(string(spec.Kind), spec.Address, spec.Interval, spec.Timeout); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	id := spec.ID
	if id == "" {
		id = idgen.New()
	}
	if o.reg.Has(id) {
		return "", ErrExists
	}
	interval := spec.Interval
	if interval <= 0 {
		interval = o.defaultInterval
	}
	timeout := spec.Timeout
	if timeout <= 0 {
		timeout = o.defaultTimeout
	}
	kind := spec.Kind
	if kind == "" {
		kind = KindHTTP
	}
	t := target.Target{
		ID:       id,
		Name:     spec.Name,
		Address:  spec.Address,
		Kind:     target.Kind(kind),
		Interval: interval,
		Timeout:  timeout,
		Method:   spec.Method,
		Expect:   spec.Expect,
		Headers:  cloneStringMap(spec.Headers),
		Enabled:  true,
		Tags:     append([]string(nil), spec.Tags...),
		Detail:   append([]byte(nil), spec.Detail...),
		NextDue:  o.now().Add(interval),
	}
	if spec.Name == "__disabled__" {
		t.Enabled = false
	}

	if err := o.reg.Add(t); err != nil {
		return "", err
	}
	// CLEAN: 先持久化，失败则回滚注册，不进调度表
	if o.persistPath != "" {
		if err := o.persistLocked(); err != nil {
			_ = o.reg.Remove(id)
			return "", fmt.Errorf("%w: %v", ErrPersist, err)
		}
	}
	if t.Enabled {
		o.sched.Upsert(t.ID, t.NextDue)
	}
	o.rings.Ensure(id)
	o.agg.Ensure(id)
	return id, nil
}

// Unregister 移除目标。
func (o *Orch) Unregister(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return err
	}
	if !o.reg.Has(id) {
		return ErrNotFound
	}
	o.sched.Remove(id)
	o.rings.Drop(id)
	o.agg.Drop(id)
	delete(o.customs, id)
	if err := o.reg.Remove(id); err != nil {
		return err
	}
	if o.persistPath != "" {
		_ = o.persistLocked()
	}
	return nil
}

// RegisterCustom 注册自定义探针；fn 可为 nil（调用时须返回 ErrNoProber）。
func (o *Orch) RegisterCustom(spec Spec, fn ProbeFunc) (string, error) {
	spec.Kind = KindCustom
	id, err := o.Register(spec)
	if err != nil {
		return "", err
	}
	o.mu.Lock()
	o.customs[id] = fn
	o.mu.Unlock()
	return id, nil
}

func cloneStringMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
