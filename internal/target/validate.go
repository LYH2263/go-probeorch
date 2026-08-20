package target

import "strings"

// Valid 粗检目标字段。
func Valid(t Target) bool {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.Address) == "" {
		return false
	}
	switch t.Kind {
	case KindHTTP, KindTCP, KindCustom, "":
		return true
	default:
		return false
	}
}
