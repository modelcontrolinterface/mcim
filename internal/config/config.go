package config

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

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

func (c *Config) WriteFile(path string) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0o644)
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

var sha256Regex = regexp.MustCompile(`^sha256-[0-9a-fA-F]{64}$`)

func validate(cfg *Config) error {
	if cfg.Sandboxes != nil {
		for name, sandbox := range cfg.Sandboxes {
			if sandbox.Hash == "" {
				return fmt.Errorf("sandbox '%s' requires a hash", name)
			}
			if !sha256Regex.MatchString(sandbox.Hash) {
				return fmt.Errorf("sandbox '%s' has an invalid hash format: %s", name, sandbox.Hash)
			}
		}
	}

	if cfg.Interceptors != nil {
		for name, interceptor := range cfg.Interceptors {
			if interceptor.Hash == "" {
				return fmt.Errorf("interceptor '%s' requires a hash", name)
			}
			if !sha256Regex.MatchString(interceptor.Hash) {
				return fmt.Errorf("interceptor '%s' has an invalid hash format: %s", name, interceptor.Hash)
			}
		}
	}

	if cfg.Services != nil {
		for name, service := range cfg.Services {
			if service.Hash == "" {
				return fmt.Errorf("service '%s' requires a hash", name)
			}
			if !sha256Regex.MatchString(service.Hash) {
				return fmt.Errorf("service '%s' has an invalid hash format: %s", name, service.Hash)
			}
		}
	}

	return nil
}
