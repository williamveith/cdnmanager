package models

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// Replicate the old "ToJSONString" for *Metadata
func (m *Metadata) ToJSONString() (string, error) {
	bytes, err := protojson.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Replicate the old "MetadataFromJSONString"
func MetadataFromJSONString(jsonStr string) (*Metadata, error) {
	var meta Metadata
	if err := protojson.Unmarshal([]byte(jsonStr), &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// Replicate the old "ToJSONString" for *Entry
func (e *Entry) ToJSONString() (string, error) {
	bytes, err := protojson.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Replicate the old "EntryFromJSONString"
func EntryFromJSONString(jsonStr string) (*Entry, error) {
	var entry Entry
	if err := protojson.Unmarshal([]byte(jsonStr), &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}
