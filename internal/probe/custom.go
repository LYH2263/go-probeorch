package probe

import (
	"context"
	"errors"
	"time"
)

// Func 自定义探针签名。
type Func func(ctx context.Context, address string, timeout time.Duration) (ok bool, detail []byte, err error)

// ErrNilFunc 空自定义函数。
var ErrNilFunc = errors.New("nil probe func")

// RunCustom 执行自定义探针；fn==nil 时返回错误而非 panic。
func RunCustom(ctx context.Context, fn Func, address string, timeout time.Duration) Response {
	if fn == nil {
		return Response{OK: false, Message: "nil", Err: ErrNilFunc}
	}
	ok, detail, err := fn(ctx, address, timeout)
	return Response{OK: ok && err == nil, Detail: append([]byte(nil), detail...), Err: err, Message: "custom"}
}
