package probeorch_test

import (
	"context"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug01_ResultDetailSliceAlias(t *testing.T) {
	o := probeorch.New(probeorch.WithHistorySize(8))
	defer o.Close()
	payload := []byte("probe-body-stable")
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "alias", Address: "local", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, payload, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := o.RunOnce(id); err != nil {
		t.Fatal(err)
	}
	hs, err := o.History(id)
	if err != nil || len(hs) < 1 {
		t.Fatalf("history: %v", err)
	}
	if len(hs[0].Detail) == 0 {
		t.Fatal("empty detail")
	}
	hs[0].Detail[0] = 'X'
	payload[0] = 'Y'
	hs2, err := o.History(id)
	if err != nil {
		t.Fatal(err)
	}
	if string(hs2[0].Detail) != "probe-body-stable" {
		t.Fatalf("history detail polluted: %q", hs2[0].Detail)
	}
}
