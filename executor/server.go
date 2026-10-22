package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	maxBody = 1 << 20

	defaultTimeSec, maxTimeSec   = 5, 10
	defaultMemoryMB, maxMemoryMB = 256, 512
)

func newServer(r Runner) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("POST /api/exec", func(w http.ResponseWriter, req *http.Request) { handleExec(r, w, req) })
	return mux
}

func clamp(raw string, def, max int) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return def
	}
	return min(n, max)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{"success": false, "message": msg})
}

func handleExec(r Runner, w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, maxBody)
	if err := req.ParseForm(); err != nil {
		fail(w, http.StatusBadRequest, "invalid or too large request body")
		return
	}
	if _, ok := languages[req.PostFormValue("language")]; !ok {
		fail(w, http.StatusBadRequest, "Invalid language")
		return
	}
	if req.PostFormValue("code") == "" {
		fail(w, http.StatusBadRequest, "code is required")
		return
	}

	res, err := r.Run(req.Context(), Request{
		Language: req.PostFormValue("language"),
		Code:     req.PostFormValue("code"),
		Stdin:    req.PostFormValue("stdin"),
		TimeSec:  clamp(req.PostFormValue("timeout"), defaultTimeSec, maxTimeSec),
		MemoryMB: clamp(req.PostFormValue("memory"), defaultMemoryMB, maxMemoryMB),
	})
	if err != nil {
		log.Printf("exec failed: %v", err)
		fail(w, http.StatusInternalServerError, "execution failed")
		return
	}
	writeJSON(w, http.StatusOK, res)
}
