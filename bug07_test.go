package probeorch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug07_RunOnceContextHonorsCancel(t *testing.T) {
	o := probeorch.New()
	defer o.Close()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "ctx", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = o.RunOnceContext(ctx, id)
	if err == nil {
		t.Fatal("want cancel error")
	}
	if !errors.Is(err, probeorch.ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
}
