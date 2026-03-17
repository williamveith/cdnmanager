package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"cdnmanager/pkg/config"
	"cdnmanager/pkg/database"
	"cdnmanager/pkg/session"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const bulkInsertTemplateName = "CDN Manager Bulk Insert Template.csv"

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

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ShowAlert(message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Alert",
		Message: message,
		Type:    runtime.InfoDialog,
	})
}

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

func (a *App) HasConfig() bool {
	cfg, err := config.LoadConfig(a.configPath)
	if err != nil {
		return false
	}
	return cfg.IsComplete()
}

func (a *App) IsConfigured() bool {
	return a.HasConfig()
}

func (a *App) GetDomain() (string, error) {
	cfg, err := config.LoadConfig(a.configPath)
	if err != nil {
		return "", err
	}
	return cfg.Domain, nil
}

func (a *App) SaveConfig(cfg config.Config) error {
	return config.SaveConfig(a.configPath, cfg)
}

// -----------------------------------------------------------------------------
// Session
// -----------------------------------------------------------------------------

func (a *App) InitializeSession() error {
	cfg, err := config.LoadConfig(a.configPath)
	if err != nil {
		return err
	}
	if !cfg.IsComplete() {
		return fmt.Errorf("config is incomplete")
	}

	sess, err := session.NewCloudflareSession(*cfg)
	if err != nil {
		return err
	}

	a.cloudflareSession = sess
	return nil
}

func (a *App) ensureSession() error {
	if a.cloudflareSession != nil {
		return nil
	}
	return a.InitializeSession()
}

// -----------------------------------------------------------------------------
// Sync
// -----------------------------------------------------------------------------

func (a *App) SyncFromCloudflare() error {
	if err := a.ensureSession(); err != nil {
		return err
	}

	cloudflareSize, _ := a.cloudflareSession.Size()
	if a.db.Size() == cloudflareSize {
		fmt.Println("Existing Database Up To Date")
		return nil
	}

	fmt.Println("Initializing Table With Cloudflare Values...")
	entries := a.cloudflareSession.GetAllEntriesBulk()
	a.db.DropTable()
	a.db.CreateTable()
	a.db.InsertEntries(entries)
	fmt.Println("Local Database Updated...")

	return nil
}

func (a *App) SetupAndSync(cfg config.Config) error {
	if !cfg.IsComplete() {
		return fmt.Errorf("config is incomplete")
	}

	if err := a.SaveConfig(cfg); err != nil {
		return err
	}

	a.cloudflareSession = nil

	if err := a.ensureSession(); err != nil {
		return err
	}

	return a.SyncFromCloudflare()
}

// -----------------------------------------------------------------------------
// Cloudflare KV actions
// -----------------------------------------------------------------------------

func (a *App) InsertKVEntry(name string, value string, metadata string) error {
	if err := a.ensureSession(); err != nil {
		return err
	}

	a.cloudflareSession.InsertKVEntry(name, value, metadata)
	return nil
}

func (a *App) DeleteKeyValue(key string) error {
	if err := a.ensureSession(); err != nil {
		return err
	}

	a.cloudflareSession.DeleteKeyValue(key)
	return nil
}

// -----------------------------------------------------------------------------
// Files
// -----------------------------------------------------------------------------

func SaveTemplateFile() (string, error) {
	csvContent := `name,value,metadata_name,metadata_external,metadata_mimetype,metadata_location,metadata_description,metadata_cloud_storage_id,metadata_md5Checksum`

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homeDir, "Downloads", bulkInsertTemplateName)
	if err := os.WriteFile(filePath, []byte(csvContent), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

func openInFinder(path string) {
	cmd := exec.Command("open", "-R", path)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error opening Finder:", err)
	}
}

func (a *App) GenerateCSV() (string, error) {
	path, err := SaveTemplateFile()
	if err != nil {
		return "", err
	}

	openInFinder(path)
	return path, nil
}
