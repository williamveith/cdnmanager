package models

import (
	"encoding/json"
	"fmt"
)

type Entry struct {
	Name     string
	Metadata Metadata
	Value    string
}

type Metadata struct {
	Name           string `json:"name"`
	External       bool   `json:"external"`
	MimeType       string `json:"mimetype"`
	Location       string `json:"location"`
	CloudStorageID string `json:"cloud_storage_id,omitempty"`
	MD5Checksum    string `json:"md5Checksum,omitempty"`
	Description    string `json:"description,omitempty"`

	Modified int64 `json:"modified,omitempty"`
}

func (m Metadata) ToJSONString() (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal metadata: %w", err)
	}
	return string(jsonBytes), nil
}

func MetadataFromJSONString(jsonStr string) (Metadata, error) {
	var meta Metadata
	if err := json.Unmarshal([]byte(jsonStr), &meta); err != nil {
		return Metadata{}, fmt.Errorf("unmarshal metadata: %w", err)
	}
	return meta, nil
}

func (e Entry) ToJSONString() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("marshal entry: %w", err)
	}
	return string(jsonBytes), nil
}

func EntryFromJSONString(jsonStr string) (Entry, error) {
	var entry Entry
	if err := json.Unmarshal([]byte(jsonStr), &entry); err != nil {
		return Entry{}, fmt.Errorf("unmarshal entry: %w", err)
	}
	return entry, nil
}
