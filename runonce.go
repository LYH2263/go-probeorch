package probeorch

import (
	"context"
	"fmt"
	"time"

	"example.com/probeorch/internal/probe"
	"example.com/probeorch/internal/result"
	"example.com/probeorch/internal/target"
)

// RunOnce 对指定目标执行一次探测。
func (o *Orch) RunOnce(id string) (Result, error) {
	return o.RunOnceContext(context.Background(), id)
}

// RunOnceContext 带取消的单次探测。
func (o *Orch) RunOnceContext(ctx context.Context, id string) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// CLEAN: 入口处响应取消
	select {
	case <-ctx.Done():
		return Result{}, fmt.Errorf("%w: %v", ErrCanceled, ctx.Err())
	default:
	}

	o.mu.Lock()
	if err := o.checkOpenLocked(); err != nil {
		o.mu.Unlock()
		return Result{}, err
	}
	if o.probers == nil {
		o.mu.Unlock()
		return Result{}, ErrClosed
	}
	t, ok := o.reg.Get(id)
	if !ok {
		o.mu.Unlock()
		return Result{}, ErrNotFound
	}
	if !t.Enabled {
		o.mu.Unlock()
		return Result{}, ErrDisabled
	}
	wait := o.wait
	delay := o.probeWaitDelay
	custom := o.customs[id]
	clientProbers := o.probers
	timeout := t.Timeout
	kind := t.Kind
	addr := t.Address
	method := t.Method
	expect := t.Expect
	headers := cloneStringMap(t.Headers)
	o.mu.Unlock()

	if delay > 0 {
		if err := wait.Wait(ctx, delay); err != nil {
			return Result{}, fmt.Errorf("%w: %v", ErrCanceled, err)
		}
	}
	select {
	case <-ctx.Done():
		return Result{}, fmt.Errorf("%w: %v", ErrCanceled, ctx.Err())
	default:
	}

	start := time.Now()
	var (
		okProbe bool
		status  int
		msg     string
		detail  []byte
		perr    error
	)

	switch kind {
	case target.KindCustom:
		if custom == nil {
			return Result{}, ErrNoProber
		}
		okProbe, detail, perr = custom(ctx, addr, timeout)
		msg = "custom"
	default:
		req := probe.Request{
			Kind:    mapProbeKind(Kind(kind)),
			Address: addr,
			Timeout: timeout,
			Method:  method,
			Expect:  expect,
			Headers: headers,
		}
		pr, err := clientProbers.Run(ctx, req)
		if err != nil {
			perr = err
		} else {
			okProbe = pr.OK
			status = pr.Status
			msg = pr.Message
			detail = pr.Detail
			perr = pr.Err
		}
	}
	lat := time.Since(start)

	res := result.Entry{
		TargetID: id,
		OK:       okProbe && perr == nil,
		At:       o.now(),
		Latency:  lat,
		Status:   status,
		Message:  msg,
		Detail:   result.CloneBytes(detail),
	}
	if perr != nil {
		res.OK = false
		res.ErrText = perr.Error()
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed || o.rings == nil {
		return Result{}, ErrClosed
	}
	o.rings.Push(id, res)
	o.agg.Observe(id, res.OK)
	if tt, ok := o.reg.Get(id); ok {
		tt.LastOK = res.OK
		tt.LastAt = res.At
		tt.NextDue = res.At.Add(tt.Interval)
		_ = o.reg.Update(tt)
		if tt.Enabled {
			o.sched.Upsert(id, tt.NextDue)
		}
	}
	o.totalRuns++
	if res.OK {
		o.okRuns++
	} else {
		o.failRuns++
	}
	o.lastID = id
	o.lastAt = res.At

	out := toPublicResult(res)
	if perr != nil {
		return out, wrapProbeErr(perr)
	}
	return out, nil
}

func toPublicResult(e result.Entry) Result {
	return Result{
		TargetID: e.TargetID,
		OK:       e.OK,
		At:       e.At,
		Latency:  e.Latency,
		Status:   e.Status,
		Message:  e.Message,
		Detail:   result.CloneBytes(e.Detail),
		ErrText:  e.ErrText,
	}
}
