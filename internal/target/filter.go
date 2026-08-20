package target

// FilterEnabled 过滤启用目标。
func FilterEnabled(in []Target) []Target {
	out := make([]Target, 0, len(in))
	for _, t := range in {
		if t.Enabled {
			out = append(out, Clone(t))
		}
	}
	return out
}

// FilterKind 按类型过滤。
func FilterKind(in []Target, k Kind) []Target {
	out := make([]Target, 0, len(in))
	for _, t := range in {
		if t.Kind == k {
			out = append(out, Clone(t))
		}
	}
	return out
}

// FindByName 按名称查找。
func FindByName(in []Target, name string) (Target, bool) {
	for _, t := range in {
		if t.Name == name {
			return Clone(t), true
		}
	}
	return Target{}, false
}
