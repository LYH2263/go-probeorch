package probe

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPProber HTTP 健康探测。
type HTTPProber struct {
	Client *http.Client
}

func (p *HTTPProber) Probe(ctx context.Context, req Request) Response {
	method := req.Method
	if method == "" {
		method = http.MethodGet
	}
	expect := req.Expect
	if expect == 0 {
		expect = http.StatusOK
	}
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(cctx, method, req.Address, nil)
	if err != nil {
		return Response{OK: false, Message: "build", Err: err}
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Response{OK: false, Message: "do", Err: err}
	}
	// BUG: 未 Close Body
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	ok := resp.StatusCode == expect || (expect == 0 && resp.StatusCode >= 200 && resp.StatusCode < 300)
	msg := strings.TrimSpace(resp.Status)
	return Response{OK: ok, Status: resp.StatusCode, Message: msg, Detail: append([]byte(nil), body...)}
}
