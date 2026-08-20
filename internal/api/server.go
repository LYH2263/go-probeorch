package api

import (
	"net/http"

	"example.com/probeorch"
)

// Server HTTP 管理面。
type Server struct {
	orch *probeorch.Orch
	mux  *http.ServeMux
	opts Options
}

// Options 服务选项。
type Options struct {
	WebDir    string
	AllowCORS bool
}

func New(o *probeorch.Orch, opts Options) *Server {
	s := &Server{orch: o, mux: http.NewServeMux(), opts: opts}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.opts.AllowCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}
