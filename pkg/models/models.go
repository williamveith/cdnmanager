package models

import (
	"encoding/json"
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
}

// Serialization methods for Metadata
func (m *Metadata) ToJSONString() (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func MetadataFromJSONString(jsonStr string) (Metadata, error) {
	var meta Metadata
	err := json.Unmarshal([]byte(jsonStr), &meta)
	if err != nil {
		return Metadata{}, err
	}
	return meta, nil
}

// Serialization methods for Entry
func (e *Entry) ToJSONString() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func EntryFromJSONString(jsonStr string) (Entry, error) {
	var entry Entry
	err := json.Unmarshal([]byte(jsonStr), &entry)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}
