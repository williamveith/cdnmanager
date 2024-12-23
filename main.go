package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"cdnmanager/pkg/database"
	"cdnmanager/pkg/session"

	"github.com/joho/godotenv"
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

//go:embed .env
var embeddedEnvFile embed.FS

var cdnDB *database.Database
var cloudflareSession *session.CloudflareSession

func loadEmbeddedEnv() {
	envBytes, _ := embeddedEnvFile.ReadFile(".env")
	envString := string(envBytes)
	envMap, _ := godotenv.Unmarshal(envString)

	for key, value := range envMap {
		_ = os.Setenv(key, value)
	}
}

func initializeDatabase(dbPath string, schema []byte) *database.Database {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database not found. Creating a new one...")
	}
	return database.NewDatabase(dbPath, schema)
}

func SyncFromCloudflare() {
	cloudflareSize, storageKeys := cloudflareSession.Size()
	if cdnDB.Size("records") != cloudflareSize {
		fmt.Println("Initializing Table With Cloudflare Values...")
		entries := cloudflareSession.GetAllEntriesFromKeys(storageKeys)
		cdnDB.RecreateTable("records")
		cdnDB.InsertEntries("records", entries)
	} else {
		fmt.Println("Existing Database Up To Date")
	}
}

func main() {
	loadEmbeddedEnv()

	// Read the schema from the embedded filesystem
	schema, schemaerror := schemaFS.ReadFile("data/schema.sql")
	if schemaerror != nil {
		fmt.Println("Failed to read embedded schema:", schemaerror)
		return
	}

	// Determine the path for the persistent database file
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to determine user config directory:", err)
		return
	}
	dbPath := filepath.Join(appDataDir, "cdnmanager", "cdnmanager.sqlite3")

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		fmt.Println("Failed to create database directory:", err)
		return
	}

	// Create a new database from the schema
	cdnDB = initializeDatabase(dbPath, schema)
	if cdnDB == nil {
		fmt.Println("Failed to initialize database")
		return
	}

	cloudflareSession = session.NewCloudflareSession()
	SyncFromCloudflare()

	app := NewApp()

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
			cloudflareSession,
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
		println("Error:", err.Error())
	}
}
