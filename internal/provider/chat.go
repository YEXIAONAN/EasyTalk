package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// rawUsage mirrors the standard OpenAI-compatible usage object.
type rawUsage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens *int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
}

func normalizeUsage(u rawUsage) Usage {
	usage := Usage{
		InputTokens:  u.PromptTokens,
		OutputTokens: u.CompletionTokens,
		TotalTokens:  u.TotalTokens,
	}
	if usage.TotalTokens == 0 {
		usage.TotalTokens = u.PromptTokens + u.CompletionTokens
	}
	if u.PromptTokensDetails.CachedTokens != nil {
		usage.CachedTokens = u.PromptTokensDetails.CachedTokens
	}
	return usage
}

// Chat sends a non-streaming completion request and returns the assistant reply
// along with its normalized token usage (nil when the provider omitted it).
func Chat(ctx context.Context, baseURL, apiKey string, req ChatRequest) (string, *Usage, error) {
	httpReq, err := newChatRequest(ctx, baseURL, apiKey, req)
	if err != nil {
		return "", nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", nil, &Error{Status: http.StatusBadGateway, Message: fmt.Sprintf("unable to connect: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, providerError(resp)
	}

	var out struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage *rawUsage `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", nil, &Error{Status: http.StatusBadGateway, Message: "invalid response from provider"}
	}
	if len(out.Choices) == 0 {
		return "", nil, &Error{Status: http.StatusBadGateway, Message: "empty response from provider"}
	}

	var usage *Usage
	if out.Usage != nil {
		u := normalizeUsage(*out.Usage)
		usage = &u
	}
	return out.Choices[0].Message.Content, usage, nil
}

// Stream sends a streaming completion request and returns the upstream response
// whose body carries an SSE stream. The caller is responsible for closing the
// body.
func Stream(ctx context.Context, baseURL, apiKey string, req ChatRequest) (*http.Response, error) {
	req.Stream = true
	req.StreamOptions = &StreamOptions{IncludeUsage: true}
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

// ParseSSE reads an OpenAI-compatible SSE stream and invokes fn for each content
// delta and, once received, the normalized usage.
func ParseSSE(r io.Reader, fn func(delta string, usage *Usage)) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
		}

		var evt struct {
			Choices []struct {
				Delta   struct{ Content string `json:"content"` } `json:"delta"`
				Message struct{ Content string `json:"content"` } `json:"message"`
			} `json:"choices"`
			Usage *rawUsage `json:"usage"`
		}
		if err := json.Unmarshal([]byte(payload), &evt); err != nil {
			continue
		}

		if evt.Usage != nil {
			u := normalizeUsage(*evt.Usage)
			fn("", &u)
		}
		for _, c := range evt.Choices {
			d := c.Delta.Content
			if d == "" {
				d = c.Message.Content
			}
			if d != "" {
				fn(d, nil)
			}
		}
	}
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