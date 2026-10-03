package server

import (
	"encoding/json"
	"fmt"
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

// chatResponse is the non-streaming response body.
type chatResponse struct {
	Role    string         `json:"role"`
	Content string         `json:"content"`
	Usage   *provider.Usage `json:"usage,omitempty"`
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

	content, usage, err := provider.Chat(r.Context(), p.BaseURL, p.APIKey, chatReq)
	if err != nil {
		s.writeChatError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, chatResponse{Role: "assistant", Content: content, Usage: usage})
}

func (s *Server) streamChat(w http.ResponseWriter, r *http.Request, baseURL, apiKey string, req provider.ChatRequest) {
	resp, err := provider.Stream(r.Context(), baseURL, apiKey, req)
	if err != nil {
		s.writeChatError(w, err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	flusher, _ := w.(http.Flusher)

	writeEvent := func(delta string, usage *provider.Usage) {
		var b []byte
		if usage != nil {
			b, _ = json.Marshal(map[string]any{"usage": usage})
		} else {
			b, _ = json.Marshal(map[string]string{"delta": delta})
		}
		_, _ = fmt.Fprintf(w, "data: %s\n\n", b)
		if flusher != nil {
			flusher.Flush()
		}
	}

	provider.ParseSSE(resp.Body, writeEvent)
	_, _ = fmt.Fprint(w, "data: [DONE]\n\n")
	if flusher != nil {
		flusher.Flush()
	}
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