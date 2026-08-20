package probeorch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/probeorch"
)

func TestBug06_PersistFailureNotScheduled(t *testing.T) {
	dir := t.TempDir()
	block := filepath.Join(dir, "notadir")
	if err := os.WriteFile(block, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	persist := filepath.Join(block, "snap.json")
	o := probeorch.New(probeorch.WithPersistPath(persist), probeorch.WithDefaultInterval(time.Millisecond))
	defer o.Close()
	_, err := o.Register(probeorch.Spec{
		Name: "p", Address: "http://127.0.0.1:9/", Kind: probeorch.KindHTTP,
		Interval: time.Millisecond, Timeout: time.Millisecond,
	})
	if err == nil {
		t.Fatal("want persist failure")
	}
	snap := o.Snapshot()
	if len(snap.Targets) != 0 {
		t.Fatalf("persist fail left target in registry: %+v", snap.Targets)
	}
	n, _ := o.TickAt(time.Now().Add(time.Hour))
	if n != 0 {
		t.Fatalf("persist fail still scheduled runs=%d", n)
	}
}
