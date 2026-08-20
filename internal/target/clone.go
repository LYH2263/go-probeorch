package target

// Clone 深拷贝单个目标。
func Clone(t Target) Target {
	out := t
	out.Tags = append([]string(nil), t.Tags...)
	out.Detail = append([]byte(nil), t.Detail...)
	if t.Headers != nil {
		out.Headers = make(map[string]string, len(t.Headers))
		for k, v := range t.Headers {
			out.Headers[k] = v
		}
	}
	return out
}

// CloneSlice 深拷贝目标切片（新底层数组）。
func CloneSlice(in []Target) []Target {
	if in == nil {
		return nil
	}
	out := make([]Target, len(in))
	for i := range in {
		out[i] = Clone(in[i])
	}
	return out
}
