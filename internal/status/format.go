package status

import (
	"fmt"
	"strings"
)

// FormatLine 单行摘要。
func FormatLine(r Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "targets=%d healthy=%d at=%s", r.Targets, r.Healthy, r.GeneratedAt.Format(timeRFC3339))
	return b.String()
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

// FormatRates 格式化成功率表。
func FormatRates(rates map[string]float64) string {
	var b strings.Builder
	for id, rate := range rates {
		fmt.Fprintf(&b, "%s=%.2f ", id, rate)
	}
	return strings.TrimSpace(b.String())
}
