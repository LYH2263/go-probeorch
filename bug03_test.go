package probeorch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug03_RunOnceAfterCloseNoPanic(t *testing.T) {
	o := probeorch.New()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "c", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RunOnce after Close panicked: %v", r)
		}
	}()
	_, err = o.RunOnce(id)
	if !errors.Is(err, probeorch.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
