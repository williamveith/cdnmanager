package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func SaveTemplateFile() (string, error) {
	csvContent := `name,value,metadata_name,metadata_external,metadata_mimetype,metadata_location,metadata_description,metadata_cloud_storage_id,metadata_md5Checksum`
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(dir, "Downloads", "CDN Manager Bulk Insert Template.csv")
	err = os.WriteFile(filePath, []byte(csvContent), 0644)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func openInFinder(path string) {
	cmd := exec.Command("open", "-R", path)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error opening Finder:", err)
	}
}

func (a *App) SetupEnvFile() error {
	// Step 1: Determine if running from an app bundle
	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}

	// App bundle's Resources directory
	resourcesDir := filepath.Join(filepath.Dir(executablePath), "../Resources")
	bundledEnvPath := filepath.Join(resourcesDir, ".env")

	// Step 2: Check if the .env file exists in the bundle & loads it if it does
	if _, err := os.Stat(bundledEnvPath); err == nil {
		return godotenv.Load(bundledEnvPath)
	}

	// Step 3: Fallback to user configuration directory for writable .env file
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}
	appConfigDir := filepath.Join(configDir, "cdnmanager")
	envPath := filepath.Join(appConfigDir, ".env")

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Println(".env file not found. Creating a default one...")
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		defaultEnv := `cloudflare_email=""
		cloudflare_api_key=""
		account_id=""
		namespace_id=""
		domain=""`

		if err := os.WriteFile(envPath, []byte(defaultEnv), 0644); err != nil {
			return fmt.Errorf("failed to write default .env file: %w", err)
		}
		return fmt.Errorf("default .env created and requeires configuration %s", envPath)
	}

	// Load the writable .env file
	return godotenv.Load(envPath)
}

func (a *App) GenerateCSV() (string, error) {
	path, err := SaveTemplateFile()
	openInFinder(path)
	return path, err
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
