package probeorch_test

import (
	"context"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug02_SnapshotTagsSliceAlias(t *testing.T) {
	o := probeorch.New()
	defer o.Close()
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "snap", Address: "local", Interval: time.Second, Timeout: time.Second,
		Tags: []string{"prod", "api"},
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = id
	snap := o.Snapshot()
	if len(snap.Targets) != 1 || len(snap.Targets[0].Tags) < 1 {
		t.Fatalf("bad snap: %+v", snap.Targets)
	}
	snap.Targets[0].Tags[0] = "HACKED"
	snap2 := o.Snapshot()
	if snap2.Targets[0].Tags[0] != "prod" {
		t.Fatalf("internal tags polluted: %v", snap2.Targets[0].Tags)
	}
}
