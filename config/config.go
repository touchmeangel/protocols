package config

import (
	"os"
	"strings"
	"time"

	"github.com/touchmeangel/protocols/defillama"
)

const defaultProxyTimeout = 60 * time.Second

type Config struct {
	Proxies []defillama.ProxyConfig
}

func Load() *Config {
	cfg := &Config{}

	raw := os.Getenv("PROXIES")
	if raw == "" {
		cfg.Proxies = []defillama.ProxyConfig{
			{Address: "", Timeout: defaultProxyTimeout},
		}
		return cfg
	}

	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		pc := defillama.ProxyConfig{Timeout: defaultProxyTimeout}

		addr, ts, hasTimeout := strings.Cut(entry, "|")
		addr = strings.TrimSpace(addr)

		if hasTimeout {
			if d, err := time.ParseDuration(strings.TrimSpace(ts)); err == nil {
				pc.Timeout = d
			}
		}

		if !strings.EqualFold(addr, "direct") {
			pc.Address = addr
		}

		cfg.Proxies = append(cfg.Proxies, pc)
	}

	if len(cfg.Proxies) == 0 {
		cfg.Proxies = []defillama.ProxyConfig{
			{Address: "", Timeout: defaultProxyTimeout},
		}
	}

	return cfg
}
