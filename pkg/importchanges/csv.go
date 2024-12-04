package importchanges

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"cdnmanager/pkg/models"
)

// CSVEntry represents a row in the CSV file.
type CSVEntry struct {
	Name     string
	Value    string
	Metadata models.Metadata
}

// ReadCSV reads a CSV file and converts each row into a CSVEntry.
func ReadCSV(filePath string) ([]CSVEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read the header row
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	var entries []CSVEntry

	// Read rows and convert to CSVEntry
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		// Convert row to map[string]string for flexibility
		data := make(map[string]string)
		for i, value := range row {
			data[headers[i]] = value
		}

		// Convert map to CSVEntry
		entry, err := mapToCSVEntry(data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row to CSVEntry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func mapToCSVEntry(data map[string]string) (CSVEntry, error) {
	external, err := strconv.ParseBool(data["metadata_external"])
	if err != nil {
		return CSVEntry{}, fmt.Errorf("invalid value for 'metadata_external': %w", err)
	}

	meta := models.Metadata{
		Name:        data["metadata_name"],
		External:    external,
		MimeType:    data["metadata_mimetype"],
		Location:    data["metadata_location"],
		Description: data["metadata_description"],
	}

	if !external {
		meta.CloudStorageID = data["metadata_cloud_storage_id"]
		meta.MD5Checksum = data["metadata_md5Checksum"]
	}

	return CSVEntry{
		Name:     data["name"],
		Value:    data["value"],
		Metadata: meta,
	}, nil
}
