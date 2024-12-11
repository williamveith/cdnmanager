package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
