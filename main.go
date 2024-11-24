package main

import (
	"embed"
	"fmt"
	"os"

	"cdnmanager/internal/database"
	"cdnmanager/internal/session"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

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

func SyncFromCloudflare() {
	cloudflareSize, storageKeys := cloudflareSession.Size()
	if cdnDB.Size() != cloudflareSize {
		fmt.Println("Initializing Table With Cloudflare Values...")
		entries := cloudflareSession.GetAllEntriesFromKeys(storageKeys)
		cdnDB.DropTable()
		cdnDB.CreateTable()
		cdnDB.InsertEntries(entries)
	} else {
		fmt.Println("Existing Database Up To Date")
	}
}

func main() {
	loadEmbeddedEnv()
	cdnDB = database.NewDatabase("data/cdn.sqlite3")
	cloudflareSession = session.NewCloudflareSession()
	SyncFromCloudflare()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "Content Delivery Network Manager",
		Width:         1024,
		Height:        768,
		DisableResize: false,
		AlwaysOnTop:   false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 100, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			cdnDB,
			cloudflareSession,
		},
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
