package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Server holds network listen settings.
type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Provider is a single OpenAI-compatible provider configuration.
type Provider struct {
	Name    string   `json:"name"`
	BaseURL string   `json:"base_url"`
	APIKey  string   `json:"api_key"`
	Models  []string `json:"models"`
}

// Config is the root configuration stored in config.json.
type Config struct {
	Server    Server     `json:"server"`
	Providers []Provider `json:"providers"`
}

// Default returns the built-in default configuration.
func Default() *Config {
	return &Config{
		Server:    Server{Host: "0.0.0.0", Port: 8080},
		Providers: []Provider{},
	}
}

// Load reads the configuration from path. If the file does not exist, it is
// created with default values so the user has a starting point to edit.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := Default()
		if err := Save(path, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the configuration to path with restricted permissions, since it
// may contain API keys.
func Save(path string, cfg *Config) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}