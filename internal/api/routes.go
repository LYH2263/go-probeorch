package api

import "net/http"

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/snapshot", s.handleSnapshot)
	s.mux.HandleFunc("/api/targets", s.handleTargets)
	s.mux.HandleFunc("/api/targets/register", s.handleRegister)
	s.mux.HandleFunc("/api/targets/run", s.handleRun)
	s.mux.HandleFunc("/api/targets/enable", s.handleEnable)
	s.mux.HandleFunc("/api/targets/disable", s.handleDisable)
	s.mux.HandleFunc("/api/targets/history", s.handleHistory)
	s.mux.HandleFunc("/api/tick", s.handleTick)
	if s.opts.WebDir != "" {
		s.mux.Handle("/", s.staticHandler())
	} else {
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "no web", http.StatusNotFound)
		})
	}
}
