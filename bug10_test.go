package probeorch_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug10_CloseFlushesBeforeDropRings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snap.json")
	o := probeorch.New(probeorch.WithPersistPath(path), probeorch.WithHistorySize(8))
	id, err := o.RegisterCustom(probeorch.Spec{
		Name: "flush", Address: "x", Interval: time.Second, Timeout: time.Second,
	}, func(ctx context.Context, address string, timeout time.Duration) (bool, []byte, error) {
		return true, []byte("kept"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := o.RunOnce(id); err != nil {
		t.Fatal(err)
	}
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	o2 := probeorch.New(probeorch.WithPersistPath(path))
	defer o2.Close()
	if err := o2.LoadPersist(); err != nil {
		t.Fatal(err)
	}
	hs, err := o2.History(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(hs) < 1 {
		t.Fatal("Close dropped result ring before Flush; history empty after reload")
	}
}
