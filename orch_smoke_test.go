package probeorch_test

import (
	"context"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestSmoke_RegisterAndHistory(t *testing.T) {
	o := probeorch.New(probeorch.WithDefaultInterval(time.Second), probeorch.WithHistorySize(8))
	defer o.Close()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "ok", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, []byte("detail"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	res, err := o.RunOnce(id)
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("want ok: %+v", res)
	}
	hs, err := o.History(id)
	if err != nil || len(hs) < 1 {
		t.Fatalf("history: %v %#v", err, hs)
	}
	snap := o.Snapshot()
	if len(snap.Targets) != 1 {
		t.Fatalf("snap targets=%d", len(snap.Targets))
	}
}
