package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	"cdnmanager/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dbName string
	db     *sql.DB
	lock   sync.Mutex
}

func (cdb *Database) GetFileName() string {
	return cdb.dbName
}

func NewDatabase(dbName string) *Database {
	db, _ := sql.Open("sqlite3", dbName)

	return &Database{
		dbName: dbName,
		db:     db,
	}
}

func NewDatabaseFromSchema(dbName string, schema []byte) *Database {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Execute the schema
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	return &Database{
		dbName: dbName,
		db:     db,
	}
}

func (cdb *Database) CreateTable() {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	cdb.db.Exec(`
		CREATE TABLE IF NOT EXISTS records (
			name TEXT PRIMARY KEY,
			value TEXT,
			metadata TEXT
		)
	`)
}

func (cdb *Database) DropTable() {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	cdb.db.Exec(`DROP TABLE IF EXISTS records`)
}

func (cdb *Database) GetEntryByName(name string) models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	var value, metadataStr string
	_ = cdb.db.QueryRow(`SELECT value, metadata FROM records WHERE name = ?`, name).Scan(&value, &metadataStr)
	metadata, _ := models.MetadataFromJSONString(metadataStr)
	entry := models.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	return entry
}

func (cdb *Database) GetEntryByValue(value string) models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	var name, metadataStr string
	_ = cdb.db.QueryRow(`SELECT name, metadata FROM records WHERE value = ?`, value).Scan(&name, &metadataStr)
	metadata, _ := models.MetadataFromJSONString(metadataStr)
	entry := models.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	return entry
}

func (cdb *Database) GetEntriesByValue(value string) []models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, _ := cdb.db.Query(`SELECT name, value, metadata FROM records WHERE value = ?`, value)
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var name, valueStr, metadataStr string
		rows.Scan(&name, &valueStr, &metadataStr)

		metadata, _ := models.MetadataFromJSONString(metadataStr)
		entries = append(entries, models.Entry{
			Name:     name,
			Metadata: metadata,
			Value:    value,
		})
	}
	return entries
}

func (cdb *Database) GetAllEntries() []models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, _ := cdb.db.Query(`SELECT name, value, metadata FROM records`)
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var name, valueStr, metadataStr string
		rows.Scan(&name, &valueStr, &metadataStr)

		metadata, _ := models.MetadataFromJSONString(metadataStr)
		entries = append(entries, models.Entry{
			Name:     name,
			Value:    valueStr,
			Metadata: metadata,
		})
	}
	return entries
}

func (cdb *Database) InsertEntry(datavalues models.Entry) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		fmt.Print("Failed To Begin Transaction:", err)
	}
	defer tx.Rollback()

	stmt, _ := tx.Prepare(`
		INSERT OR REPLACE INTO records (name, value, metadata)
		VALUES (?, ?, ?)
	`)
	defer stmt.Close()
	metadata, _ := datavalues.Metadata.ToJSONString()
	stmt.Exec(datavalues.Name, datavalues.Value, metadata)
	err = tx.Commit()
	if err != nil {
		fmt.Print("Failed To Commit Transaction:", err)
	}
}

func (cdb *Database) InsertKVEntryIntoDatabase(name string, value string, metadata string) {
	Metadata, _ := models.MetadataFromJSONString(metadata)
	newEntry := models.Entry{
		Name:     name,
		Metadata: Metadata,
		Value:    value,
	}
	cdb.InsertEntry(newEntry)
}

func (cdb *Database) InsertEntries(datavalues []models.Entry) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		fmt.Print("Failed To Begin Transaction:", err)
	}
	defer tx.Rollback()

	stmt, _ := tx.Prepare(`
		INSERT OR REPLACE INTO records (name, value, metadata)
		VALUES (?, ?, ?)
	`)
	defer stmt.Close()

	for _, datavalue := range datavalues {
		metadata, _ := datavalue.Metadata.ToJSONString()
		stmt.Exec(datavalue.Name, datavalue.Value, metadata)
	}

	tx.Commit()
	if err != nil {
		fmt.Print("Failed To Commit Transaction:", err)
	}
}

func (cdb *Database) DeleteName(key string) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	_, _ = cdb.db.Exec(`DELETE FROM records WHERE name = ?`, key)
}

func (cdb *Database) DeleteNames(names []string) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	query := `DELETE FROM records WHERE name IN (?` + strings.Repeat(",?", len(names)-1) + `)`
	args := make([]interface{}, len(names))
	for i, name := range names {
		args[i] = name
	}
	cdb.db.Exec(query, args...)
}

func (cdb *Database) DeleteEntry(entry models.Entry) {
	cdb.DeleteName(entry.Name)
}

func (cdb *Database) DeleteEntries(entries []models.Entry) {
	for _, entry := range entries {
		cdb.DeleteEntry(entry)
	}
}

func (cdb *Database) Size() int {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	query := `SELECT COUNT(*) FROM records`

	var rowCount int
	cdb.db.QueryRow(query).Scan(&rowCount)

	return rowCount
}
