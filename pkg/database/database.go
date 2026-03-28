package database

import (
	"database/sql"
	"fmt"
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

func NewDatabase(dbName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	return &Database{
		dbName: dbName,
		db:     db,
	}, nil
}

func NewDatabaseFromSchema(dbName string, schema []byte) (*Database, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	if _, err = db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize schema: %w", err)
	}

	return &Database{
		dbName: dbName,
		db:     db,
	}, nil
}

func (cdb *Database) CreateTable() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	if _, err := cdb.db.Exec(`
		CREATE TABLE IF NOT EXISTS records (
			name TEXT PRIMARY KEY,
			value TEXT,
			metadata TEXT
		)
	`); err != nil {
		return fmt.Errorf("create table records: %w", err)
	}
	return nil
}

func (cdb *Database) DropTable() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	if _, err := cdb.db.Exec(`DROP TABLE IF EXISTS records`); err != nil {
		return fmt.Errorf("drop table records: %w", err)
	}
	return nil
}

func (cdb *Database) GetEntryByName(name string) (models.Entry, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	var value, metadataStr string
	err := cdb.db.QueryRow(
		`SELECT value, metadata FROM records WHERE name = ?`,
		name,
	).Scan(&value, &metadataStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Entry{}, nil
		}
		return models.Entry{}, fmt.Errorf("get entry by name %q: %w", name, err)
	}

	metadata, err := models.MetadataFromJSONString(metadataStr)
	if err != nil {
		return models.Entry{}, fmt.Errorf("parse metadata for %q: %w", name, err)
	}

	return models.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}, nil
}

func (cdb *Database) GetEntryByValue(value string) (models.Entry, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	var name, metadataStr string
	err := cdb.db.QueryRow(
		`SELECT name, metadata FROM records WHERE value = ?`,
		value,
	).Scan(&name, &metadataStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Entry{}, nil
		}
		return models.Entry{}, fmt.Errorf("get entry by value %q: %w", value, err)
	}

	metadata, err := models.MetadataFromJSONString(metadataStr)
	if err != nil {
		return models.Entry{}, fmt.Errorf("parse metadata for %q: %w", name, err)
	}

	return models.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}, nil
}

func (cdb *Database) GetEntriesByValue(value string) ([]models.Entry, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, err := cdb.db.Query(`SELECT name, value, metadata FROM records WHERE value = ?`, value)
	if err != nil {
		return nil, fmt.Errorf("query entries by value %q: %w", value, err)
	}
	defer rows.Close()

	entries := make([]models.Entry, 0)
	for rows.Next() {
		var name, valueStr, metadataStr string
		if err := rows.Scan(&name, &valueStr, &metadataStr); err != nil {
			return nil, fmt.Errorf("scan entry by value %q: %w", value, err)
		}

		metadata, err := models.MetadataFromJSONString(metadataStr)
		if err != nil {
			return nil, fmt.Errorf("parse metadata for %q: %w", name, err)
		}

		entries = append(entries, models.Entry{
			Name:     name,
			Metadata: metadata,
			Value:    valueStr,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate entries by value %q: %w", value, err)
	}

	return entries, nil
}

func (cdb *Database) GetAllEntries() ([]models.Entry, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, err := cdb.db.Query(`SELECT name, value, metadata FROM records`)
	if err != nil {
		return nil, fmt.Errorf("query all entries: %w", err)
	}
	defer rows.Close()

	entries := make([]models.Entry, 0)
	for rows.Next() {
		var name, valueStr, metadataStr string
		if err := rows.Scan(&name, &valueStr, &metadataStr); err != nil {
			return nil, fmt.Errorf("scan entry: %w", err)
		}

		metadata, err := models.MetadataFromJSONString(metadataStr)
		if err != nil {
			return nil, fmt.Errorf("parse metadata for %q: %w", name, err)
		}

		entries = append(entries, models.Entry{
			Name:     name,
			Value:    valueStr,
			Metadata: metadata,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate all entries: %w", err)
	}

	return entries, nil
}

func (cdb *Database) UpsertEntry(entry models.Entry) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		return fmt.Errorf("begin upsert transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO records (name, value, metadata)
		VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			value = excluded.value,
			metadata = excluded.metadata
	`)
	if err != nil {
		return fmt.Errorf("prepare upsert statement: %w", err)
	}
	defer stmt.Close()

	metadata, err := entry.Metadata.ToJSONString()
	if err != nil {
		return fmt.Errorf("serialize metadata for %q: %w", entry.Name, err)
	}

	if _, err := stmt.Exec(entry.Name, entry.Value, metadata); err != nil {
		return fmt.Errorf("upsert entry %q: %w", entry.Name, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit upsert transaction: %w", err)
	}
	return nil
}

func (cdb *Database) UpsertEntries(entries []models.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO records (name, value, metadata)
		VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			value = excluded.value,
			metadata = excluded.metadata
	`)
	if err != nil {
		return fmt.Errorf("prepare upsert statement: %w", err)
	}
	defer stmt.Close()

	for _, entry := range entries {
		metadataJSON, err := entry.Metadata.ToJSONString()
		if err != nil {
			return fmt.Errorf("serialize metadata for %q: %w", entry.Name, err)
		}

		if _, err := stmt.Exec(entry.Name, entry.Value, metadataJSON); err != nil {
			return fmt.Errorf("upsert entry %q: %w", entry.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit upsert transaction: %w", err)
	}

	return nil
}

func (cdb *Database) DeleteName(key string) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	if _, err := cdb.db.Exec(`DELETE FROM records WHERE name = ?`, key); err != nil {
		return fmt.Errorf("delete name %q: %w", key, err)
	}

	return nil
}

func (cdb *Database) DeleteNames(names []string) error {
	if len(names) == 0 {
		return nil
	}

	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	query := `DELETE FROM records WHERE name IN (?` + strings.Repeat(",?", len(names)-1) + `)`
	args := make([]interface{}, len(names))
	for i, name := range names {
		args[i] = name
	}

	if _, err := cdb.db.Exec(query, args...); err != nil {
		return fmt.Errorf("delete names: %w", err)
	}

	return nil
}

func (cdb *Database) DeleteEntry(entry models.Entry) error {
	return cdb.DeleteName(entry.Name)
}

func (cdb *Database) DeleteEntries(entries []models.Entry) error {
	for _, entry := range entries {
		if err := cdb.DeleteEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

func (cdb *Database) Size() (int, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	var rowCount int
	if err := cdb.db.QueryRow(`SELECT COUNT(*) FROM records`).Scan(&rowCount); err != nil {
		return 0, fmt.Errorf("count records: %w", err)
	}

	return rowCount, nil
}
