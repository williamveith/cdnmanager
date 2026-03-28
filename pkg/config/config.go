package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const CurrentVersion = 1

type Config struct {
	Version int `json:"version"`

	CloudflareAPIToken string `json:"cloudflare_api_token"`
	AccountID          string `json:"account_id"`
	NamespaceID        string `json:"namespace_id"`
	Domain             string `json:"domain"`
}

func (c *Config) normalize() {
	if c == nil {
		return
	}
	c.CloudflareAPIToken = strings.TrimSpace(c.CloudflareAPIToken)
	c.AccountID = strings.TrimSpace(c.AccountID)
	c.NamespaceID = strings.TrimSpace(c.NamespaceID)
	c.Domain = strings.TrimSpace(c.Domain)
}

func (c Config) IsComplete() bool {
	return strings.TrimSpace(c.CloudflareAPIToken) != "" &&
		strings.TrimSpace(c.AccountID) != "" &&
		strings.TrimSpace(c.NamespaceID) != "" &&
		strings.TrimSpace(c.Domain) != ""
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", configPath, err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config %q: %w", configPath, err)
	}
	if dec.More() {
		return nil, fmt.Errorf("decode config %q: multiple JSON objects", configPath)
	}

	// handle missing version (older files or manual edits)
	if cfg.Version == 0 {
		cfg.Version = CurrentVersion
	}
	if cfg.Version > CurrentVersion {
		return nil, fmt.Errorf("unsupported config version %d", cfg.Version)
	}

	cfg.normalize()

	return &cfg, nil
}

func SaveConfig(configPath string, cfg Config) error {
	cfg.Version = CurrentVersion
	cfg.normalize()

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return fmt.Errorf("create config directory for %q: %w", configPath, err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config %q: %w", configPath, err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("write config %q: %w", configPath, err)
	}

	return nil
}
