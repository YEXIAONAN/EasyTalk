package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Chat sends a non-streaming completion request and returns the assistant reply.
func Chat(ctx context.Context, baseURL, apiKey string, req ChatRequest) (string, error) {
	httpReq, err := newChatRequest(ctx, baseURL, apiKey, req)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", &Error{Status: http.StatusBadGateway, Message: fmt.Sprintf("unable to connect: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", providerError(resp)
	}

	var out struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", &Error{Status: http.StatusBadGateway, Message: "invalid response from provider"}
	}
	if len(out.Choices) == 0 {
		return "", &Error{Status: http.StatusBadGateway, Message: "empty response from provider"}
	}
	return out.Choices[0].Message.Content, nil
}

// Stream sends a streaming completion request and returns the upstream response
// whose body carries an SSE stream. The caller is responsible for closing the
// body.
func Stream(ctx context.Context, baseURL, apiKey string, req ChatRequest) (*http.Response, error) {
	req.Stream = true
	httpReq, err := newChatRequest(ctx, baseURL, apiKey, req)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, &Error{Status: http.StatusBadGateway, Message: fmt.Sprintf("unable to connect: %v", err)}
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, providerError(resp)
	}
	return resp, nil
}

func newChatRequest(ctx context.Context, baseURL, apiKey string, req ChatRequest) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, normalizeBase(baseURL)+"/chat/completions", bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	}
	return httpReq, nil
}

// providerError builds an Error from a non-2xx upstream response.
func providerError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

	var parsed struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	msg := strings.TrimSpace(string(body))
	if json.Unmarshal(body, &parsed) == nil && parsed.Error.Message != "" {
		msg = parsed.Error.Message
	}
	if msg == "" {
		msg = fmt.Sprintf("provider returned status %d", resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &Error{Status: http.StatusUnauthorized, Message: ErrInvalidKey.Error()}
	case http.StatusTooManyRequests:
		return &Error{Status: http.StatusTooManyRequests, Message: msg}
	default:
		return &Error{Status: http.StatusBadGateway, Message: msg}
	}
}