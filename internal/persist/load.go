package persist

import (
	"encoding/json"
	"os"

	"example.com/probeorch/internal/result"
)

// Load 读取快照。
func Load(path string) (Snapshot, error) {
	var snap Snapshot
	data, err := os.ReadFile(path)
	if err != nil {
		return snap, err
	}
	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, err
	}
	if snap.History == nil {
		snap.History = make(map[string][]result.Entry)
	}
	return snap, nil
}
