package server

import (
	"encoding/json"
	"io"
	"net/http"

	"easytalk/internal/provider"
)

// chatRequest is the body the frontend sends to start a chat.
type chatRequest struct {
	Provider    string             `json:"provider"`
	Model       string             `json:"model"`
	Messages    []provider.Message `json:"messages"`
	Stream      bool               `json:"stream"`
	Temperature *float64           `json:"temperature"`
	MaxTokens   *int               `json:"max_tokens"`
	TopP        *float64           `json:"top_p"`
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, ok := s.findProvider(req.Provider)
	if !ok {
		writeError(w, http.StatusNotFound, "provider not found")
		return
	}

	chatReq := provider.ChatRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Stream:      req.Stream,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	}

	if req.Stream {
		s.streamChat(w, r, p.BaseURL, p.APIKey, chatReq)
		return
	}

	content, err := provider.Chat(r.Context(), p.BaseURL, p.APIKey, chatReq)
	if err != nil {
		s.writeChatError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"role":    "assistant",
		"content": content,
	})
}

func (s *Server) streamChat(w http.ResponseWriter, r *http.Request, baseURL, apiKey string, req provider.ChatRequest) {
	upstream, err := provider.Stream(r.Context(), baseURL, apiKey, req)
	if err != nil {
		s.writeChatError(w, err)
		return
	}
	defer upstream.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	_, _ = io.Copy(w, upstream.Body)
}

func (s *Server) writeChatError(w http.ResponseWriter, err error) {
	status := http.StatusBadGateway
	msg := err.Error()
	if pe, ok := err.(*provider.Error); ok {
		status = pe.Status
		msg = pe.Message
	}
	writeError(w, status, msg)
}