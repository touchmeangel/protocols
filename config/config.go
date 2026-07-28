package config

import (
	"os"
	"strings"
)

type Config struct {
	Proxies []string
}

func Load() *Config {
	cfg := &Config{}

	if proxies := os.Getenv("PROXIES"); proxies != "" {
		cfg.Proxies = strings.Split(proxies, ",")
	}

	for i := range cfg.Proxies {
		cfg.Proxies[i] = strings.TrimSpace(cfg.Proxies[i])
	}

	return cfg
}
