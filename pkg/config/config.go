package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const CurrentVersion = 1

type Config struct {
	Version int `json:"version"`

	CloudflareAPIToken string `json:"cloudflare_api_token"`
	AccountID          string `json:"account_id"`
	NamespaceID        string `json:"namespace_id"`
	Domain             string `json:"domain"`
}

func (c Config) IsComplete() bool {
	return c.CloudflareAPIToken != "" &&
		c.AccountID != "" &&
		c.NamespaceID != "" &&
		c.Domain != ""
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// handle missing version (older files or manual edits)
	if cfg.Version == 0 {
		cfg.Version = CurrentVersion
	}

	return &cfg, nil
}

func SaveConfig(configPath string, cfg Config) error {
	cfg.Version = CurrentVersion

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o600)
}
