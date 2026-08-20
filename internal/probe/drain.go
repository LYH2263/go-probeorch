package probe

import "io"

// DrainAndClose 排空并关闭 HTTP 响应体，避免连接池句柄泄漏。
func DrainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 64<<10))
	_ = body.Close()
}
