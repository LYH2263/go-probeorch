package probe

import "fmt"

// FormatOK 成功消息。
func FormatOK(kind Kind, status int) string {
	return fmt.Sprintf("%s ok status=%d", kind, status)
}

// FormatFail 失败消息。
func FormatFail(kind Kind, err error) string {
	if err == nil {
		return fmt.Sprintf("%s fail", kind)
	}
	return fmt.Sprintf("%s fail: %v", kind, err)
}
