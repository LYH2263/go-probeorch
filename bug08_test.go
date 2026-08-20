package probeorch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug08_ProbeWaitHonorsContext(t *testing.T) {
	o := probeorch.New(
		probeorch.WithProbeWaitDelay(2 * time.Second),
	)
	defer o.Close()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "wait", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err = o.RunOnceContext(ctx, id)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("want cancel during wait")
	}
	if !errors.Is(err, probeorch.ErrCanceled) && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("wait ignored ctx, elapsed=%v", elapsed)
	}
}
