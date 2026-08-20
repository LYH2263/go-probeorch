package validate

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyAddress = errors.New("empty address")
	ErrBadKind      = errors.New("bad kind")
	ErrBadInterval  = errors.New("interval too small")
	ErrBadTimeout   = errors.New("timeout too small")
)

// Spec 校验注册参数。
func Spec(kind, address string, interval, timeout time.Duration) error {
	if strings.TrimSpace(address) == "" {
		return ErrEmptyAddress
	}
	k := strings.ToLower(strings.TrimSpace(kind))
	if k != "" && k != "http" && k != "tcp" && k != "custom" {
		return ErrBadKind
	}
	if interval < 0 {
		return ErrBadInterval
	}
	if timeout < 0 {
		return ErrBadTimeout
	}
	return nil
}

// MaxNameLen 名称长度上限。
func MaxNameLen() int { return 128 }

// ClampTimeout 限制超时范围。
func ClampTimeout(d, def, max time.Duration) time.Duration {
	if d <= 0 {
		return def
	}
	if d > max {
		return max
	}
	return d
}
