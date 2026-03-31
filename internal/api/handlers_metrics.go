package api

import (
	"net/http"
	"strconv"
	"time"
)

func (s *Server) handleEnvMetrics(w http.ResponseWriter, r *http.Request) {
	if s.metrics == nil {
		http.Error(w, `{"error":"metrics not enabled"}`, http.StatusServiceUnavailable)
		return
	}
	env := r.PathValue("env")
	from, to, step := parseTimeRange(r)

	points, err := s.metrics.QueryEnvMetrics(env, from, to, step)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"points": points})
}

func (s *Server) handleServerMetrics(w http.ResponseWriter, r *http.Request) {
	if s.metrics == nil {
		http.Error(w, `{"error":"metrics not enabled"}`, http.StatusServiceUnavailable)
		return
	}
	env := r.PathValue("env")
	serverID := r.URL.Query().Get("server_id")
	from, to, step := parseTimeRange(r)

	points, err := s.metrics.QueryServerMetrics(env, serverID, from, to, step)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"points": points})
}

func (s *Server) handleMQTTBridgeMetrics(w http.ResponseWriter, r *http.Request) {
	if s.metrics == nil {
		http.Error(w, `{"error":"metrics not enabled"}`, http.StatusServiceUnavailable)
		return
	}
	env := r.PathValue("env")
	bridgeID := r.URL.Query().Get("bridge_id")
	from, to, step := parseTimeRange(r)

	points, err := s.metrics.QueryMQTTMetrics(env, bridgeID, from, to, step)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"points": points})
}

func parseTimeRange(r *http.Request) (from, to, step int64) {
	now := time.Now().Unix()
	to = now
	from = now - 3600 // default 1 hour

	if v := r.URL.Query().Get("from"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			from = n
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			to = n
		}
	}
	if v := r.URL.Query().Get("step"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			step = n
		}
	}
	return
}
