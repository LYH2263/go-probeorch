package probeorch

import "errors"

var (
	ErrClosed      = errors.New("probeorch: closed")
	ErrNotFound    = errors.New("probeorch: target not found")
	ErrExists      = errors.New("probeorch: target already exists")
	ErrInvalid     = errors.New("probeorch: invalid argument")
	ErrDisabled    = errors.New("probeorch: target disabled")
	ErrNoProber    = errors.New("probeorch: no probe func")
	ErrProbeFailed = errors.New("probeorch: probe failed")
	ErrPersist     = errors.New("probeorch: persist failed")
	ErrCanceled    = errors.New("probeorch: canceled")
)
