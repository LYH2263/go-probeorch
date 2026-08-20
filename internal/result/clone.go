package result

// CloneBytes 深拷贝字节切片，避免与调用方共享底层数组。
func CloneBytes(src []byte) []byte {
	if src == nil {
		return nil
	}
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// CloneEntry 深拷贝结果条目。
func CloneEntry(e Entry) Entry {
	out := e
	out.Detail = CloneBytes(e.Detail)
	return out
}
