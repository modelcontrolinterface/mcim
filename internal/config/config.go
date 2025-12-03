package config

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	applyDefaults(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	t := true

	if cfg.Services != nil {
		for name, service := range cfg.Services {
			if service.Enable == nil {
				service.Enable = &t
				cfg.Services[name] = service
			}
		}
	}

	if cfg.Sandboxes != nil {
		for name, sandbox := range cfg.Sandboxes {
			if sandbox.Enable == nil {
				sandbox.Enable = &t
				cfg.Sandboxes[name] = sandbox
			}
		}
	}

	if cfg.Interceptors != nil {
		for name, interceptor := range cfg.Interceptors {
			if interceptor.Enable == nil {
				interceptor.Enable = &t
				cfg.Interceptors[name] = interceptor
			}
		}
	}
}

func validateHash(s string) error {
	validHashRegexes := map[string]*regexp.Regexp{
		"sha256-": regexp.MustCompile(`^[a-fA-F0-9]{64}$`),
	}

	for prefix, re := range validHashRegexes {
		if strings.HasPrefix(s, prefix) {
			hashPart := strings.TrimPrefix(s, prefix)
			if !re.MatchString(hashPart) {
				return fmt.Errorf("hash does not match expected format for prefix %q", prefix)
			}
			return nil
		}
	}

	return fmt.Errorf("unsupported hash prefix")
}

func validate(cfg *Config) error {
	if cfg.Sandboxes != nil {
		for name, sandbox := range cfg.Sandboxes {
			err := validateHash(sandbox.Hash)
			if err != nil {
				return fmt.Errorf("sandbox '%s' configuration error: %s", name, err)
			}
		}
	}

	if cfg.Interceptors != nil {
		for name, interceptor := range cfg.Interceptors {
			err := validateHash(interceptor.Hash)
			if err != nil {
				return fmt.Errorf("interceptor '%s' configuration error: %s", name, err)
			}
		}
	}

	if cfg.Services != nil {
		for name, service := range cfg.Services {
			err := validateHash(service.Hash)
			if err != nil {
				return fmt.Errorf("service '%s' configuration error: %s", name, err)
			}
		}
	}

	return nil
}

func (c *Config) WriteFile(path string) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0o644)
}
