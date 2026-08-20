package probeorch

import (
	"fmt"

	"example.com/probeorch/internal/probe"
)

// wrapProbeErr 将探针失败包装为可 errors.Is 的哨兵链。
func wrapProbeErr(err error) error {
	if err == nil {
		return nil
	}
	// 用 %w 挂上哨兵 ErrProbeFailed 与上游错误，errors.Is 均可命中。
	return fmt.Errorf("%w: %w", ErrProbeFailed, err)
}

// mapProbeKind 映射对外 Kind 到 internal。
func mapProbeKind(k Kind) probe.Kind {
	switch k {
	case KindTCP:
		return probe.KindTCP
	case KindCustom:
		return probe.KindCustom
	default:
		return probe.KindHTTP
	}
}
