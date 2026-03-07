package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"cdnmanager/pkg/database"
	"cdnmanager/pkg/session"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx               context.Context
	db                *database.Database
	configPath        string
	cloudflareSession *session.CloudflareSession
}

func NewApp(db *database.Database, configPath string) *App {
	return &App{
		db:         db,
		configPath: configPath,
	}
}

func (a *App) HasConfig() bool {
	cfg, err := LoadConfig(a.configPath)
	if err != nil {
		return false
	}
	return cfg.IsComplete()
}

func (a *App) GetConfig() (*Config, error) {
	return LoadConfig(a.configPath)
}

func (a *App) SaveConfig(cfg Config) error {
	return SaveConfig(a.configPath, cfg)
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

func (a *App) GenerateCSV() (string, error) {
	path, err := SaveTemplateFile()
	openInFinder(path)
	return path, err
}

func (a *App) ShowAlert(message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Alert",
		Message: message,
		Type:    runtime.InfoDialog,
	})
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) InitializeSession() error {
	cfg, err := LoadConfig(a.configPath)
	if err != nil {
		return err
	}
	if !cfg.IsComplete() {
		return fmt.Errorf("config is incomplete")
	}

	sess, err := session.NewCloudflareSession(session.Config{
		CloudflareEmail:  cfg.CloudflareEmail,
		CloudflareAPIKey: cfg.CloudflareAPIKey,
		AccountID:        cfg.AccountID,
		NamespaceID:      cfg.NamespaceID,
		Domain:           cfg.Domain,
	})
	if err != nil {
		return err
	}

	a.cloudflareSession = sess
	return nil
}

func (a *App) SyncFromCloudflare() error {
	if a.cloudflareSession == nil {
		if err := a.InitializeSession(); err != nil {
			return err
		}
	}

	cloudflareSize, storageKeys := a.cloudflareSession.Size()
	if a.db.Size() != cloudflareSize {
		fmt.Println("Initializing Table With Cloudflare Values...")
		entries := a.cloudflareSession.GetAllEntriesFromKeys(storageKeys)
		a.db.DropTable()
		a.db.CreateTable()
		a.db.InsertEntries(entries)
		fmt.Println("Local Database Updated...")
	} else {
		fmt.Println("Existing Database Up To Date")
	}

	return nil
}

func (a *App) SetupAndSync(cfg Config) error {
	if !cfg.IsComplete() {
		return fmt.Errorf("config is incomplete")
	}

	if err := SaveConfig(a.configPath, cfg); err != nil {
		return err
	}

	if err := a.InitializeSession(); err != nil {
		return err
	}

	return a.SyncFromCloudflare()
}

func (a *App) IsConfigured() bool {
	return a.HasConfig()
}

func (a *App) InsertKVEntry(name string, value string, metadata string) error {
	if a.cloudflareSession == nil {
		if err := a.InitializeSession(); err != nil {
			return err
		}
	}

	a.cloudflareSession.InsertKVEntry(name, value, metadata)
	return nil
}

func (a *App) DeleteKeyValue(key string) error {
	if a.cloudflareSession == nil {
		if err := a.InitializeSession(); err != nil {
			return err
		}
	}

	a.cloudflareSession.DeleteKeyValue(key)
	return nil
}
