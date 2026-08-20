package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var seq uint64

// New 生成短 ID。
func New() string {
	n := atomic.AddUint64(&seq, 1)
	var b [4]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("t%x-%d-%s", time.Now().UnixNano()&0xffffff, n, hex.EncodeToString(b[:]))
}
