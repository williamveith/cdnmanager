package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/williamveith/cdnmanager/internal/session"
)

type Database struct {
	dbName string
	db     *sql.DB
	lock   sync.Mutex
}

func NewDatabase(dbName string) *Database {
	db, _ := sql.Open("sqlite3", dbName)

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

func (cdb *Database) InsertEntry(datavalues session.Entry) {
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
	metadata := convertMetadataToString(datavalues.Metadata)
	stmt.Exec(datavalues.Name, datavalues.Value, metadata)
	err = tx.Commit()
	if err != nil {
		fmt.Print("Failed To Commit Transaction:", err)
	}
}

func (cdb *Database) InsertEntries(datavalues []session.Entry) {
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
		stmt.Exec(datavalue.Name, datavalue.Value, convertMetadataToString(datavalue.Metadata))
	}

	tx.Commit()
	if err != nil {
		fmt.Print("Failed To Commit Transaction:", err)
	}
}

func (cdb *Database) GetEntryByName(name string) session.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	var value, metadata string
	_ = cdb.db.QueryRow(`SELECT value, metadata FROM records WHERE name = ?`, name).Scan(&value, &metadata)
	entry := session.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	return entry
}

func (cdb *Database) GetEntryByValue(value string) session.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	var name, metadata string
	_ = cdb.db.QueryRow(`SELECT name, metadata FROM records WHERE value = ?`, value).Scan(&name, &metadata)
	entry := session.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	return entry
}

func (cdb *Database) GetEntriesByValue(value string) []session.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, _ := cdb.db.Query(`SELECT name, value, metadata FROM records WHERE value = ?`, value)
	defer rows.Close()

	var entries []session.Entry
	for rows.Next() {
		var name, metadata string
		rows.Scan(&name, &metadata)
		entries = append(entries, session.Entry{
			Name:     name,
			Metadata: metadata,
			Value:    value,
		})
	}
	return entries
}

func (cdb *Database) GetAllEntries() []session.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, _ := cdb.db.Query(`SELECT name, value, metadata FROM records`)
	defer rows.Close()

	var entries []session.Entry
	for rows.Next() {
		var name, value, metadata string
		rows.Scan(&name, &value, &metadata)
		entries = append(entries, session.Entry{
			Name:     name,
			Value:    value,
			Metadata: metadata,
		})
	}
	return entries
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

func (cdb *Database) DeleteEntry(entry session.Entry) {
	cdb.DeleteName(entry.Name)
}

func (cdb *Database) DeleteEntries(entries []session.Entry) {
	for _, entry := range entries {
		cdb.DeleteEntry(entry)
	}
}

func convertMetadataToString(metadata interface{}) string {
	switch v := metadata.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "Error converting metadata to string"
		}
		return string(jsonBytes)
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
