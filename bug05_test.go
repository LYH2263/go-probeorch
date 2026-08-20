package probeorch_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug05_ProbeErrorWrapsSentinel(t *testing.T) {
	o := probeorch.New()
	defer o.Close()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "fail", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return false, nil, fmt.Errorf("upstream down")
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = o.RunOnce(id)
	if err == nil {
		t.Fatal("want error")
	}
	if !errors.Is(err, probeorch.ErrProbeFailed) {
		t.Fatalf("want errors.Is ErrProbeFailed, got %v", err)
	}
}
