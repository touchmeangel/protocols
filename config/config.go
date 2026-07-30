package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/touchmeangel/protocols/defillama"
)

const defaultProxyTimeout = 60 * time.Second

type jsonDuration time.Duration

func (d *jsonDuration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	*d = jsonDuration(parsed)
	return nil
}

type proxyEntry struct {
	Address string       `json:"address"`
	Timeout jsonDuration `json:"timeout"`
}

type fileConfig struct {
	Proxies []proxyEntry `json:"proxies"`
}

type Config struct {
	Proxies []defillama.ProxyConfig
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}

	cfg := &Config{}

	for _, p := range fc.Proxies {
		timeout := time.Duration(p.Timeout)
		if timeout == 0 {
			timeout = defaultProxyTimeout
		}
		cfg.Proxies = append(cfg.Proxies, defillama.ProxyConfig{
			Address: p.Address,
			Timeout: timeout,
		})
	}
	if len(cfg.Proxies) == 0 {
		cfg.Proxies = []defillama.ProxyConfig{{Address: "", Timeout: defaultProxyTimeout}}
	}

	return cfg, nil
}
