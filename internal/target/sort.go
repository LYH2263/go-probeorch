package target

import "sort"

// SortByNextDue 按下次到期排序。
func SortByNextDue(in []Target) []Target {
	out := CloneSlice(in)
	sort.Slice(out, func(i, j int) bool {
		return out[i].NextDue.Before(out[j].NextDue)
	})
	return out
}

// SortByName 按名称排序。
func SortByName(in []Target) []Target {
	out := CloneSlice(in)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}
