package main

import (
	"embed"
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
	initializeEnvironment()
	initializeSession()
	initializeDatabase()
}

func initializeEnvironment() {
	loadEmbeddedEnv()
}

func initializeSession() {
	cloudflareSession = session.NewCloudflareSession()
}

func initializeDatabase() {
	cdnDB = database.NewDatabase("data/cdn.sqlite3")
	synchronizeDatabaseWithCloudflare()
}

func synchronizeDatabaseWithCloudflare() {
	keys := session.GetAllKeys(cloudflareSession)
	numberOfKeysWeb := len(keys)
	numberOfKeysLocal, _ := cdnDB.GetRowCount("records")

	if numberOfKeysWeb > numberOfKeysLocal {
		response := session.GetAllEntries(cloudflareSession)
		cdnDB.DropTable("records")
		cdnDB.CreateTable()
		cdnDB.InsertEntries(response)
	}
}

func loadEmbeddedEnv() {
	envBytes, _ := embeddedEnvFile.ReadFile(".env")
	envString := string(envBytes)
	envMap, _ := godotenv.Unmarshal(envString)

	for key, value := range envMap {
		_ = os.Setenv(key, value)
	}
}

func main() {
	session.GetAllKeys(cloudflareSession)
	response := session.GetAllEntries(cloudflareSession)
	cdnDB.InsertEntries(response)
}
