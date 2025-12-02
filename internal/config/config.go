package config

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

func FromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) ToFile(path string) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0o644)
}
