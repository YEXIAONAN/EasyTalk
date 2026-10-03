package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatNonStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test" {
			t.Fatalf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"hi"}}],"usage":{"prompt_tokens":10,"completion_tokens":20,"total_tokens":30,"prompt_tokens_details":{"cached_tokens":5}}}`))
	}))
	defer srv.Close()

	content, usage, err := Chat(context.Background(), srv.URL, "sk-test", ChatRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if content != "hi" {
		t.Fatalf("content = %q, want %q", content, "hi")
	}
	if usage == nil || usage.InputTokens != 10 || usage.OutputTokens != 20 || usage.TotalTokens != 30 {
		t.Fatalf("unexpected usage: %+v", usage)
	}
	if usage.CachedTokens == nil || *usage.CachedTokens != 5 {
		t.Fatalf("unexpected cached tokens: %+v", usage)
	}
}

func TestChatInvalidKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer srv.Close()

	_, _, err := Chat(context.Background(), srv.URL, "sk-test", ChatRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	pe, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if pe.Status != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", pe.Status, http.StatusUnauthorized)
	}
}

func TestStreamAndParseSSE(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(
			"data: {\"choices\":[{\"delta\":{\"content\":\"hel\"}}]}\n\n" +
				"data: {\"choices\":[{\"delta\":{\"content\":\"lo\"}}]}\n\n" +
				"data: {\"choices\":[],\"usage\":{\"prompt_tokens\":1,\"completion_tokens\":2}}\n\n" +
				"data: [DONE]\n\n",
		))
	}))
	defer srv.Close()

	resp, err := Stream(context.Background(), srv.URL, "sk-test", ChatRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	defer resp.Body.Close()

	var deltas []string
	var usage *Usage
	ParseSSE(resp.Body, func(delta string, u *Usage) {
		if u != nil {
			usage = u
		} else {
			deltas = append(deltas, delta)
		}
	})

	if strings.Join(deltas, "") != "hello" {
		t.Fatalf("deltas = %v, want [hel lo]", deltas)
	}
	if usage == nil || usage.InputTokens != 1 || usage.OutputTokens != 2 {
		t.Fatalf("unexpected usage: %+v", usage)
	}
}