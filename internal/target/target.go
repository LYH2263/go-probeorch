package target

import "time"

// Kind 探针类型。
type Kind string

const (
	KindHTTP   Kind = "http"
	KindTCP    Kind = "tcp"
	KindCustom Kind = "custom"
)

// Target 内部目标模型。
type Target struct {
	ID       string
	Name     string
	Address  string
	Kind     Kind
	Interval time.Duration
	Timeout  time.Duration
	Method   string
	Expect   int
	Headers  map[string]string
	Enabled  bool
	Tags     []string
	Detail   []byte
	NextDue  time.Time
	LastOK   bool
	LastAt   time.Time
}
