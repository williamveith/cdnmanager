package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"cdnmanager/pkg/config"
	"cdnmanager/pkg/database"
	"cdnmanager/pkg/models"
	"cdnmanager/pkg/reconcile"
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

	session, err := session.NewCloudflareSession(*cfg)
	if err != nil {
		return err
	}

	a.cloudflareSession = session
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

	cloudflareEntries, err := a.cloudflareSession.GetAllEntriesBulk()
	if err != nil {
		return fmt.Errorf("fetch cloudflare entries: %w", err)
	}

	databaseEntries, err := a.db.GetAllEntries()
	if err != nil {
		return fmt.Errorf("fetch database entries: %w", err)
	}

	plan, err := reconcile.Reconcile(cloudflareEntries, databaseEntries)
	if err != nil {
		return fmt.Errorf("reconcile entries: %w", err)
	}

	if err := a.db.DeleteNames(plan.ToDelete); err != nil {
		return fmt.Errorf("delete stale database entries: %w", err)
	}

	toWrite := make([]models.Entry, 0, len(plan.ToInsert)+len(plan.ToUpdate))
	toWrite = append(toWrite, plan.ToInsert...)
	toWrite = append(toWrite, plan.ToUpdate...)

	if err := a.db.UpsertEntries(toWrite); err != nil {
		return fmt.Errorf("upsert database entries: %w", err)
	}

	fmt.Printf(
		"Sync complete. Inserted: %d, Updated: %d, Deleted: %d\n",
		len(plan.ToInsert),
		len(plan.ToUpdate),
		len(plan.ToDelete),
	)

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
// data operation primitives
// -----------------------------------------------------------------------------

func (a *App) Insert(name string, value string, metadata string) error {
	if err := a.ensureSession(); err != nil {
		return err
	}

	meta, err := models.MetadataFromJSONString(metadata)
	if err != nil {
		return err
	}

	newEntry := models.Entry{
		Name:     name,
		Metadata: meta,
		Value:    value,
	}

	newEntry.Metadata.Modified = time.Now().Unix()

	if err := a.cloudflareSession.WriteEntry(newEntry); err != nil {
		return err
	}
	a.db.InsertEntry(newEntry)

	return nil
}

func (a *App) Delete(key string) error {
	if err := a.ensureSession(); err != nil {
		return err
	}
	if err := a.cloudflareSession.DeleteKeyValue(key); err != nil {
		return err
	}
	if err := a.db.DeleteName(key); err != nil {
		return err
	}
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

func (a *App) GenerateCSV() (string, error) {
	path, err := SaveTemplateFile()
	if err != nil {
		return "", err
	}

	openInFinder(path)
	return path, nil
}

func openInFinder(path string) {
	cmd := exec.Command("open", "-R", path)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error opening Finder:", err)
	}
}
