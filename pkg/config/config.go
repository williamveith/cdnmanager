package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	CloudflareEmail  string `json:"cloudflare_email"`
	CloudflareAPIKey string `json:"cloudflare_api_key"`
	AccountID        string `json:"account_id"`
	NamespaceID      string `json:"namespace_id"`
	Domain           string `json:"domain"`
}

func (c Config) IsComplete() bool {
	return c.CloudflareEmail != "" &&
		c.CloudflareAPIKey != "" &&
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

	return &cfg, nil
}

func SaveConfig(configPath string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o600)
}
