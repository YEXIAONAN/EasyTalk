package server

import (
	"errors"
	"net/http"

	"easytalk/internal/config"
)

// providerSummary is the minimal provider info exposed to the frontend. The
// API key is deliberately excluded.
type providerSummary struct {
	Name   string   `json:"name"`
	Models []string `json:"models"`
}

func (s *Server) handleListProviders(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	providers := append([]config.Provider(nil), s.cfg.Providers...)
	s.mu.Unlock()

	out := make([]providerSummary, len(providers))
	for i, p := range providers {
		out[i] = providerSummary{Name: p.Name, Models: p.Models}
	}
	writeJSON(w, http.StatusOK, out)
}

// findProvider returns the provider with the given name (with its real API
// key), if it exists. Only used internally by the chat handler.
func (s *Server) findProvider(name string) (config.Provider, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.cfg.Providers {
		if p.Name == name {
			return p, true
		}
	}
	return config.Provider{}, false
}

func (s *Server) handleReloadConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load(s.configPath)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			writeError(w, http.StatusNotFound, "config file not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.mu.Lock()
	s.cfg = cfg
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}