// Package provider implements the OpenAI-compatible API client used to talk to
// AI providers.
package provider

import (
	"errors"
	"strings"
)

// ErrInvalidKey is returned when a provider rejects the supplied API key.
var ErrInvalidKey = errors.New("invalid api key")

// Error is an upstream error carrying an HTTP status for the client.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

// Message is a single chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the subset of the OpenAI-compatible request that EasyTalk uses.
type ChatRequest struct {
	Model         string         `json:"model"`
	Messages      []Message      `json:"messages"`
	Stream        bool           `json:"stream"`
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	Temperature   *float64       `json:"temperature,omitempty"`
	MaxTokens     *int           `json:"max_tokens,omitempty"`
	TopP          *float64       `json:"top_p,omitempty"`
}

// StreamOptions requests additional data on the streamed response.
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

// Usage is EasyTalk's normalized token usage.
type Usage struct {
	InputTokens  int  `json:"input_tokens"`
	OutputTokens int  `json:"output_tokens"`
	TotalTokens  int  `json:"total_tokens"`
	CachedTokens *int `json:"cached_tokens,omitempty"`
}

// normalizeBase trims surrounding whitespace and a trailing slash from a base
// URL.
func normalizeBase(baseURL string) string {
	return strings.TrimRight(strings.TrimSpace(baseURL), "/")
}