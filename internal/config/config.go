package config

import (
	"encoding/json"
	"errors"
	"os"
)

// ErrNotFound is returned when the config file does not exist.
var ErrNotFound = errors.New("config file not found")

// Provider is a single OpenAI-compatible provider configuration.
type Provider struct {
	Name    string   `json:"name"`
	BaseURL string   `json:"base_url"`
	APIKey  string   `json:"api_key"`
	Models  []string `json:"models"`
}

// Config is the root configuration stored in config.json.
type Config struct {
	Providers []Provider `json:"providers"`
}

// Load reads the provider configuration from path. It returns ErrNotFound when
// the file does not exist so the caller can show a helpful message.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}