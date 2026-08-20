package validate

import "strings"

// NormalizeName 去掉首尾空白。
func NormalizeName(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeTags 过滤空标签。
func NormalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// HasTag 判断是否包含标签。
func HasTag(tags []string, want string) bool {
	for _, t := range tags {
		if t == want {
			return true
		}
	}
	return false
}
