package probe

import (
	"context"
	"fmt"
	"net/http"
)

// Registry 内置探针集合。
type Registry struct {
	http *HTTPProber
	tcp  TCPProber
}

func NewRegistry(c *http.Client) *Registry {
	return &Registry{http: &HTTPProber{Client: c}, tcp: TCPProber{}}
}

func (r *Registry) Run(ctx context.Context, req Request) (Response, error) {
	if r == nil {
		return Response{}, fmt.Errorf("nil probe registry")
	}
	switch req.Kind {
	case KindTCP:
		return r.tcp.Probe(ctx, req), nil
	case KindCustom:
		return Response{OK: false, Message: "use custom path"}, fmt.Errorf("custom via orch")
	default:
		return r.http.Probe(ctx, req), nil
	}
}
