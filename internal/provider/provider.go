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
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	TopP        *float64  `json:"top_p,omitempty"`
}

// normalizeBase trims surrounding whitespace and a trailing slash from a base
// URL.
func normalizeBase(baseURL string) string {
	return strings.TrimRight(strings.TrimSpace(baseURL), "/")
}