package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/probeorch"
	"example.com/probeorch/internal/api"
)

func main() {
	addr := flag.String("addr", ":8100", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	persist := flag.String("persist", "", "快照 JSON 路径（可选）")
	flag.Parse()

	opts := []probeorch.Option{}
	if *persist != "" {
		opts = append(opts, probeorch.WithPersistPath(*persist))
	}
	o := probeorch.New(opts...)
	defer o.Close()
	if *persist != "" {
		_ = o.LoadPersist()
	}

	srv := api.New(o, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("probed 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
