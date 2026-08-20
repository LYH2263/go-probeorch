package probe

import (
	"context"
	"net"
	"time"
)

// TCPProber TCP 连通探测。
type TCPProber struct{}

func (TCPProber) Probe(ctx context.Context, req Request) Response {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", req.Address)
	if err != nil {
		return Response{OK: false, Message: "dial", Err: err}
	}
	_ = conn.Close()
	return Response{OK: true, Message: "ok", Detail: []byte("tcp-ok")}
}
