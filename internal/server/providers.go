package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"easytalk/internal/config"
)

// maskKey hides all but the last few characters of an API key so it can be
// safely shown in responses and logs.
func maskKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 4 {
		return "****"
	}
	return "****" + key[len(key)-4:]
}

// isMasked reports whether the key is a masked placeholder sent back from the
// frontend when the user left the field untouched.
func isMasked(key string) bool {
	return strings.Contains(key, "*")
}

func (s *Server) handleListProviders(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	providers := append([]config.Provider(nil), s.cfg.Providers...)
	s.mu.Unlock()

	out := make([]config.Provider, len(providers))
	for i, p := range providers {
		p.APIKey = maskKey(p.APIKey)
		out[i] = p
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleCreateProvider(w http.ResponseWriter, r *http.Request) {
	var p config.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p.Name = strings.TrimSpace(p.Name)
	p.BaseURL = strings.TrimSpace(p.BaseURL)
	if p.Name == "" || p.BaseURL == "" {
		writeError(w, http.StatusBadRequest, "name and base_url are required")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.cfg.Providers {
		if existing.Name == p.Name {
			writeError(w, http.StatusConflict, "provider already exists")
			return
		}
	}

	s.cfg.Providers = append(s.cfg.Providers, p)
	if err := s.saveLocked(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save config")
		return
	}

	p.APIKey = maskKey(p.APIKey)
	writeJSON(w, http.StatusCreated, p)
}

func (s *Server) handleUpdateProvider(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	var p config.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p.Name = strings.TrimSpace(p.Name)
	p.BaseURL = strings.TrimSpace(p.BaseURL)
	if p.Name == "" || p.BaseURL == "" {
		writeError(w, http.StatusBadRequest, "name and base_url are required")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.cfg.Providers {
		if existing.Name == name {
			if isMasked(p.APIKey) {
				p.APIKey = existing.APIKey
			}
			s.cfg.Providers[i] = p
			if err := s.saveLocked(); err != nil {
				writeError(w, http.StatusInternalServerError, "failed to save config")
				return
			}
			p.APIKey = maskKey(p.APIKey)
			writeJSON(w, http.StatusOK, p)
			return
		}
	}

	writeError(w, http.StatusNotFound, "provider not found")
}

func (s *Server) handleDeleteProvider(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.cfg.Providers {
		if existing.Name == name {
			s.cfg.Providers = append(s.cfg.Providers[:i], s.cfg.Providers[i+1:]...)
			if err := s.saveLocked(); err != nil {
				writeError(w, http.StatusInternalServerError, "failed to save config")
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	writeError(w, http.StatusNotFound, "provider not found")
}

// saveLocked persists the current config. The caller must hold s.mu.
func (s *Server) saveLocked() error {
	return config.Save(s.configPath, s.cfg)
}