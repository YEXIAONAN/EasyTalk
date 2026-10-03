package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Env holds runtime (non-provider) configuration.
type Env struct {
	Host       string
	Port       int
	ConfigPath string
}

// LoadEnv resolves runtime settings. Precedence is:
// default value < .env file < OS environment variable.
func LoadEnv() Env {
	file := map[string]string{}

	if f, err := os.Open(".env"); err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			k, v, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			file[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
		_ = f.Close()
	}

	// OS environment overrides .env.
	get := func(key string) string {
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		return file[key]
	}

	port := 8080
	if v := get("EASYTALK_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			port = n
		}
	}

	return Env{
		Host:       orDefault(get("EASYTALK_HOST"), "0.0.0.0"),
		Port:       port,
		ConfigPath: orDefault(get("EASYTALK_CONFIG"), "./config.json"),
	}
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}