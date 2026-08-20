package probe

import "strings"

// MergeHeaders 合并默认头与请求头。
func MergeHeaders(base, over map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range base {
		out[k] = v
	}
	for k, v := range over {
		out[k] = v
	}
	return out
}

// CanonicalHeader 规范化头名。
func CanonicalHeader(k string) string {
	return strings.Title(strings.ToLower(k))
}

// HasHeader 判断头是否存在（大小写不敏感）。
func HasHeader(h map[string]string, key string) bool {
	want := strings.ToLower(key)
	for k := range h {
		if strings.ToLower(k) == want {
			return true
		}
	}
	return false
}
