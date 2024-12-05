package main

import (
	"context"
	"embed"
	"fmt"
	"os"

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

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
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
				Message: "© 2024 William Veith",
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

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
