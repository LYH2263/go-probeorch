package probeorch_test

import (
	"errors"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug04_NilCustomProbeNoPanic(t *testing.T) {
	o := probeorch.New()
	defer o.Close()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "nilfn", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil ProbeFunc panicked: %v", r)
		}
	}()
	_, err = o.RunOnce(id)
	if !errors.Is(err, probeorch.ErrNoProber) {
		t.Fatalf("want ErrNoProber, got %v", err)
	}
}
