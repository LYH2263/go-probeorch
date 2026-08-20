package api

import (
	"context"
	"net"
	"net/http"
	"time"
)

// ListenAndServe 启动 HTTP 服务。
func ListenAndServe(addr string, h http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}

// ServeListener 在已有 listener 上服务。
func ServeListener(ln net.Listener, h http.Handler) error {
	srv := &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second}
	return srv.Serve(ln)
}

// Shutdown 优雅关闭。
func Shutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}
