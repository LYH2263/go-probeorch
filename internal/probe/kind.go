package probe

import "time"

// Kind 探针种类。
type Kind string

const (
	KindHTTP   Kind = "http"
	KindTCP    Kind = "tcp"
	KindCustom Kind = "custom"
)

// Request 探针请求。
type Request struct {
	Kind    Kind
	Address string
	Timeout time.Duration
	Method  string
	Expect  int
	Headers map[string]string
}

// Response 探针响应。
type Response struct {
	OK      bool
	Status  int
	Message string
	Detail  []byte
	Err     error
}
