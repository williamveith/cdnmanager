package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/williamveith/cdnmanager/internal/database"
	"github.com/williamveith/cdnmanager/internal/session"
)

var cdnDB *database.Database
var cloudflareSession *session.CloudflareSession

//go:embed .env
var embeddedEnvFile embed.FS

func init() {
	loadEmbeddedEnv()
	cdnDB = database.NewDatabase("data/cdn.sqlite3")
	cloudflareSession = session.NewCloudflareSession()
	SyncFromCloudflare()
}

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
	if cdnDB.Size() < cloudflareSize {
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
	fmt.Print("Complete")
}
