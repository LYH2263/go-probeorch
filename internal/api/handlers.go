package api

import (
	"net/http"
	"strconv"
	"time"

	"example.com/probeorch"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.Stats())
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.orch.Snapshot())
}

func (s *Server) handleTargets(w http.ResponseWriter, r *http.Request) {
	snap := s.orch.Snapshot()
	writeJSON(w, http.StatusOK, snap.Targets)
}

type registerBody struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Kind     string            `json:"kind"`
	Interval string            `json:"interval"`
	Timeout  string            `json:"timeout"`
	Method   string            `json:"method"`
	Expect   int               `json:"expect"`
	Headers  map[string]string `json:"headers"`
	Tags     []string          `json:"tags"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	var body registerBody
	if err := decodeJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	interval, _ := time.ParseDuration(body.Interval)
	timeout, _ := time.ParseDuration(body.Timeout)
	id, err := s.orch.Register(probeorch.Spec{
		ID: body.ID, Name: body.Name, Address: body.Address,
		Kind: probeorch.Kind(body.Kind), Interval: interval, Timeout: timeout,
		Method: body.Method, Expect: body.Expect, Headers: body.Headers, Tags: body.Tags,
		Enabled: true,
	})
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeErr(w, http.StatusBadRequest, probeorch.ErrInvalid)
		return
	}
	res, err := s.orch.RunOnce(id)
	if err != nil && !res.OK {
		writeJSON(w, http.StatusOK, map[string]any{"result": res, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleEnable(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := s.orch.Enable(id); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "enabled"})
}

func (s *Server) handleDisable(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := s.orch.Disable(id); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "disabled"})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	hs, err := s.orch.History(id)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	limit := 20
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			limit = n
		}
	}
	if len(hs) > limit {
		hs = hs[:limit]
	}
	writeJSON(w, http.StatusOK, hs)
}

func (s *Server) handleTick(w http.ResponseWriter, r *http.Request) {
	n, err := s.orch.Tick()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"ran": n})
}
