package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"easytalk/internal/config"
	"easytalk/internal/version"
)

// Server wires up the HTTP routes and serves the frontend assets.
type Server struct {
	cfg        *config.Config
	configPath string
	static     fs.FS
	mu         sync.Mutex
	mux        *http.ServeMux
}

// New creates a Server. static is the frontend asset filesystem (dist content
// at its root); it may be nil while the frontend is not yet available.
func New(cfg *config.Config, configPath string, static fs.FS) *Server {
	s := &Server{
		cfg:        cfg,
		configPath: configPath,
		static:     static,
		mux:        http.NewServeMux(),
	}
	s.routes()
	return s
}

// Handler returns the root HTTP handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.handleHealth)

	s.mux.HandleFunc("GET /api/providers", s.handleListProviders)
	s.mux.HandleFunc("POST /api/providers", s.handleCreateProvider)
	s.mux.HandleFunc("PUT /api/providers/{name}", s.handleUpdateProvider)
	s.mux.HandleFunc("DELETE /api/providers/{name}", s.handleDeleteProvider)
	s.mux.HandleFunc("POST /api/providers/{name}/test", s.handleTestProvider)

	s.mux.HandleFunc("POST /api/chat", s.handleChat)

	if s.static != nil {
		fileServer := http.FileServerFS(s.static)
		s.mux.Handle("/", fileServer)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"version": version.Version,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}