package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"cdnmanager/pkg/database"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed data/schema.sql
var schemaFS embed.FS

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/assets/images/appicon.png
var icon []byte

const appFolderName = "cdnmanager"

func initializeDatabase(dbPath string) (*database.Database, error) {
	if _, err := os.Stat(dbPath); err == nil {
		return database.NewDatabase(dbPath), nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat database: %w", err)
	}

	fmt.Println("Database not found. Creating a new one...")

	schema, err := schemaFS.ReadFile("data/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded schema: %w", err)
	}

	db := database.NewDatabaseFromSchema(dbPath, schema)
	if db == nil {
		return nil, fmt.Errorf("database creation returned nil")
	}

	return db, nil
}

func initializeConfig(configPath string) error {
	if _, err := os.Stat(configPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat config: %w", err)
	}

	fmt.Println("Config not found. Creating a new one...")

	return SaveConfig(configPath, Config{})
}

func appPaths() (appDir string, dbPath string, configPath string, err error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to determine user config directory: %w", err)
	}

	appDir = filepath.Join(userConfigDir, appFolderName)
	dbPath = filepath.Join(appDir, "cdnmanager.sqlite3")
	configPath = filepath.Join(appDir, "config.json")

	return appDir, dbPath, configPath, nil
}

func main() {
	appDir, dbPath, configPath, err := appPaths()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll(appDir, 0755); err != nil {
		fmt.Println("Failed to create app directory:", err)
		return
	}

	if err := initializeConfig(configPath); err != nil {
		fmt.Println("Failed to initialize config:", err)
		return
	}

	cdnDB, err := initializeDatabase(dbPath)
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}

	app := NewApp(cdnDB, configPath)

	err = wails.Run(&options.App{
		Title:         "Content Delivery Network Manager",
		DisableResize: false,
		MinWidth:      1400,
		MinHeight:     800,
		AlwaysOnTop:   false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 100, G: 38, B: 54, A: 50},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			cdnDB,
		},
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisablePinchZoom:                  false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 true,
				HideToolbarSeparator:       false,
			},
			Preferences: &mac.Preferences{
				FullscreenEnabled: mac.Enabled,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Content Delivery Network Manager",
				Message: "Created By William Veith",
				Icon:    icon,
			},
		},
		Linux: &linux.Options{
			Icon:                icon,
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "wails",
		},
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
