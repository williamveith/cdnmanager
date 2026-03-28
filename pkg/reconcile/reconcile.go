package reconcile

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"cdnmanager/pkg/models"
)

type Plan struct {
	ToInsert []models.Entry
	ToUpdate []models.Entry
	ToDelete []string
}

func Reconcile(cloudflareEntries, databaseEntries []models.Entry) (Plan, error) {
	cloudflareMap := make(map[string]models.Entry, len(cloudflareEntries))
	databaseMap := make(map[string]models.Entry, len(databaseEntries))

	for _, entry := range cloudflareEntries {
		cloudflareMap[entry.Name] = entry
	}

	for _, entry := range databaseEntries {
		databaseMap[entry.Name] = entry
	}

	plan := Plan{
		ToInsert: make([]models.Entry, 0),
		ToUpdate: make([]models.Entry, 0),
		ToDelete: make([]string, 0),
	}

	for key, cloudflareEntry := range cloudflareMap {
		databaseEntry, exists := databaseMap[key]
		if !exists {
			plan.ToInsert = append(plan.ToInsert, cloudflareEntry)
			continue
		}

		cloudflareHash, err := HashEntry(cloudflareEntry)
		if err != nil {
			return Plan{}, fmt.Errorf("hash cloudflare entry %q: %w", key, err)
		}

		databaseHash, err := HashEntry(databaseEntry)
		if err != nil {
			return Plan{}, fmt.Errorf("hash database entry %q: %w", key, err)
		}

		if cloudflareHash != databaseHash {
			plan.ToUpdate = append(plan.ToUpdate, cloudflareEntry)
		}
	}

	for key := range databaseMap {
		if _, exists := cloudflareMap[key]; !exists {
			plan.ToDelete = append(plan.ToDelete, key)
		}
	}

	sort.Strings(plan.ToDelete)

	return plan, nil
}

func HashEntry(entry models.Entry) (string, error) {
	normalized, err := normalizeEntry(entry)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(normalized)
	return hex.EncodeToString(sum[:]), nil
}

func normalizeEntry(entry models.Entry) ([]byte, error) {
	metadataJSON, err := json.Marshal(entry.Metadata)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	payload := struct {
		Name     string          `json:"name"`
		Value    string          `json:"value"`
		Metadata json.RawMessage `json:"metadata"`
	}{
		Name:     entry.Name,
		Value:    entry.Value,
		Metadata: metadataJSON,
	}

	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal normalized entry: %w", err)
	}

	return normalized, nil
}
