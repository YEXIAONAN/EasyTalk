package provider

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ErrInvalidKey is returned when the provider rejects the supplied API key.
var ErrInvalidKey = errors.New("invalid api key")

const testTimeout = 10 * time.Second

// Test verifies that a provider is reachable and its API key is accepted by
// requesting the OpenAI-compatible /models endpoint.
func Test(baseURL, apiKey string) error {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		return errors.New("base url is empty")
	}

	req, err := http.NewRequest(http.MethodGet, base+"/models", nil)
	if err != nil {
		return err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := (&http.Client{Timeout: testTimeout}).Do(req)
	if err != nil {
		return fmt.Errorf("unable to connect: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrInvalidKey
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
}