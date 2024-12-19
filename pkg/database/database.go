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
	schema string
	db     *sql.DB
	lock   sync.Mutex
}

func (cdb *Database) GetFileName() string {
	return cdb.dbName
}

func NewDatabase(dbName string, schema []byte) *Database {
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
		schema: string(schema),
		db:     db,
	}
}

func (cdb *Database) RecreateTable(tableName string) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	createStmt := ""
	statements := strings.Split(cdb.schema, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if strings.HasPrefix(stmt, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", tableName)) {
			createStmt = stmt
			break
		}
	}

	if createStmt == "" {
		return fmt.Errorf("table %s not found in schema file", tableName)
	}

	// Drop the table
	_, err := cdb.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName))
	if err != nil {
		return fmt.Errorf("failed to drop table %s: %w", tableName, err)
	}

	// Recreate the table
	_, err = cdb.db.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	return nil
}

func (cdb *Database) GetEntryByName(tableName string, name string) models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	query := fmt.Sprintf(`SELECT value, metadata FROM %s WHERE name = ?`, tableName)
	var value, metadataStr string
	_ = cdb.db.QueryRow(query, name).Scan(&value, &metadataStr)
	metadata, _ := models.MetadataFromJSONString(metadataStr)
	entry := models.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	return entry
}

func (cdb *Database) GetEntryByValue(tableName string, value string) models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	query := fmt.Sprintf(`SELECT name, metadata FROM %s WHERE value = ?`, tableName)
	var name, metadataStr string
	_ = cdb.db.QueryRow(query, value).Scan(&name, &metadataStr)
	metadata, _ := models.MetadataFromJSONString(metadataStr)
	entry := models.Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	return entry
}

func (cdb *Database) GetEntriesByValue(tableName string, value string) []models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	query := fmt.Sprintf(`SELECT name, value, metadata FROM %s WHERE value = ?`, tableName)
	rows, _ := cdb.db.Query(query, value)
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

func (cdb *Database) GetAllEntries(tableName string) []models.Entry {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	query := fmt.Sprintf("SELECT name, value, metadata FROM %s", tableName)
	rows, _ := cdb.db.Query(query)
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

func (cdb *Database) InsertEntry(tableName string, datavalues models.Entry) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		fmt.Print("Failed To Begin Transaction:", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
		INSERT OR REPLACE INTO %s (name, value, metadata)
		VALUES (?, ?, ?)
	`, tableName)

	stmt, _ := tx.Prepare(query)
	defer stmt.Close()
	metadata, _ := datavalues.Metadata.ToJSONString()
	stmt.Exec(datavalues.Name, datavalues.Value, metadata)
	err = tx.Commit()
	if err != nil {
		fmt.Print("Failed To Commit Transaction:", err)
	}
}

func (cdb *Database) InsertKVEntryIntoDatabase(tableName string, name string, value string, metadata string) {
	Metadata, _ := models.MetadataFromJSONString(metadata)
	newEntry := models.Entry{
		Name:     name,
		Metadata: Metadata,
		Value:    value,
	}
	cdb.InsertEntry(tableName, newEntry)
}

func (cdb *Database) InsertEntries(tableName string, datavalues []models.Entry) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		fmt.Print("Failed To Begin Transaction:", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
		INSERT OR REPLACE INTO %s (name, value, metadata)
		VALUES (?, ?, ?)
	`, tableName)

	stmt, _ := tx.Prepare(query)
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

func (cdb *Database) DeleteName(tableName string, key string) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	query := fmt.Sprintf(`DELETE FROM %s WHERE name = ?`, tableName)
	_, _ = cdb.db.Exec(query, key)
}

func (cdb *Database) DeleteNames(tableName string, names []string) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	query := fmt.Sprintf(`DELETE FROM %s WHERE name IN (?`+strings.Repeat(",?", len(names)-1)+`)`, tableName)
	args := make([]interface{}, len(names))
	for i, name := range names {
		args[i] = name
	}
	cdb.db.Exec(query, args...)
}

func (cdb *Database) DeleteEntry(tableName string, entry models.Entry) {
	cdb.DeleteName(tableName, entry.Name)
}

func (cdb *Database) DeleteEntries(tableName string, entries []models.Entry) {
	for _, entry := range entries {
		cdb.DeleteEntry(tableName, entry)
	}
}

func (cdb *Database) Size(tableName string) int {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tableName)

	var rowCount int
	cdb.db.QueryRow(query).Scan(&rowCount)

	return rowCount
}
