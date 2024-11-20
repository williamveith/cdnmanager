package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	cdb := &Database{
		dbName: dbName,
		db:     db,
	}

	cdb.CreateTable()

	return cdb
}

func (cdb *Database) CreateTable() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	_, err := cdb.db.Exec(`
		CREATE TABLE IF NOT EXISTS records (
			name TEXT PRIMARY KEY,
			value TEXT,
			metadata TEXT
		)
	`)
	return err
}

func (cdb *Database) GetRowCount(tableName string) (int, error) {
	var rowCount int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	err := cdb.db.QueryRow(query).Scan(&rowCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count: %w", err)
	}

	return rowCount, nil
}

func (cdb *Database) DropTable(tableName string) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := cdb.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop table %s: %w", tableName, err)
	}

	return nil
}

func (cdb *Database) InsertEntry(datavalues session.Entry) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, _ := tx.Prepare(`
		INSERT INTO records (name, value, metadata)
		VALUES (?, ?, ?)
	`)
	defer stmt.Close()
	metadata := convertMetadataToString(datavalues.Metadata)
	stmt.Exec(datavalues.Name, datavalues.Value, metadata)
	return tx.Commit()
}

func (cdb *Database) Close() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	return cdb.db.Close()
}

func (cdb *Database) InsertEntries(datavalues []session.Entry) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, _ := tx.Prepare(`
		INSERT INTO records (name, value, metadata)
		VALUES (?, ?, ?)
	`)
	defer stmt.Close()

	for _, datavalue := range datavalues {
		metadata := convertMetadataToString(datavalue.Metadata)
		stmt.Exec(datavalue.Name, datavalue.Value, metadata)
	}

	return tx.Commit()
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
