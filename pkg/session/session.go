package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"cdnmanager/pkg/models"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Config struct {
	CloudflareEmail  string
	CloudflareAPIKey string
	AccountID        string
	NamespaceID      string
	Domain           string
}

func (c Config) IsComplete() bool {
	return c.CloudflareEmail != "" &&
		c.CloudflareAPIKey != "" &&
		c.AccountID != "" &&
		c.NamespaceID != "" &&
		c.Domain != ""
}

type CloudflareSession struct {
	api         *cloudflare.API
	accountID   *cloudflare.ResourceContainer
	namespaceID string
	domain      string
}

func NewCloudflareSession(cfg Config) (*CloudflareSession, error) {
	if !cfg.IsComplete() {
		return nil, fmt.Errorf("cloudflare session config is incomplete")
	}

	api, err := cloudflare.New(cfg.CloudflareAPIKey, cfg.CloudflareEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudflare api client: %w", err)
	}

	cloudflareSession := &CloudflareSession{
		api:         api,
		accountID:   cloudflare.AccountIdentifier(cfg.AccountID),
		namespaceID: cfg.NamespaceID,
		domain:      cfg.Domain,
	}

	return cloudflareSession, nil
}

func (cloudflareSession *CloudflareSession) GetValue(key string) string {
	resp, _ := cloudflareSession.api.GetWorkersKV(
		context.Background(),
		cloudflareSession.accountID,
		cloudflare.GetWorkersKVParams{
			NamespaceID: cloudflareSession.namespaceID,
			Key:         key,
		},
	)

	return string(resp)
}

func (cloudflareSession *CloudflareSession) GetAllValues() []string {
	storageKeys := cloudflareSession.GetAllKeys()
	var values []string
	for _, entry := range storageKeys {
		values = append(values, cloudflareSession.GetValue(entry.Name))
	}
	return values
}

func (cloudflareSession *CloudflareSession) GetAllKeys() []cloudflare.StorageKey {
	resp, _ := cloudflareSession.api.ListWorkersKVKeys(
		context.Background(),
		cloudflareSession.accountID,
		cloudflare.ListWorkersKVsParams{
			NamespaceID: cloudflareSession.namespaceID,
		},
	)
	return resp.Result
}

func (cloudflareSession *CloudflareSession) GetAllEntries() []models.Entry {
	storageKeys := cloudflareSession.GetAllKeys()
	var entries []models.Entry
	for _, sk := range storageKeys {
		var metadata models.Metadata

		if sk.Metadata != nil {
			metadataJSON, _ := json.Marshal(sk.Metadata)
			metadata, _ = models.MetadataFromJSONString(string(metadataJSON))
		} else {
			metadata = models.Metadata{}
		}

		entry := models.Entry{
			Name:     sk.Name,
			Metadata: metadata,
			Value:    cloudflareSession.GetValue(sk.Name),
		}
		entries = append(entries, entry)
	}
	return entries
}

func (cloudflareSession *CloudflareSession) GetAllEntriesFromKeys(storageKeys []cloudflare.StorageKey) []models.Entry {
	if len(storageKeys) == 0 {
		return []models.Entry{}
	}

	entries := make([]models.Entry, len(storageKeys))

	const workerCount = 18
	type job struct {
		index int
		key   cloudflare.StorageKey
	}

	jobs := make(chan job, len(storageKeys))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for j := range jobs {
			var metadata models.Metadata

			if j.key.Metadata != nil {
				metadataJSON, err := json.Marshal(j.key.Metadata)
				if err == nil {
					_ = json.Unmarshal(metadataJSON, &metadata)
				}
			}

			entries[j.index] = models.Entry{
				Name:     j.key.Name,
				Metadata: metadata,
				Value:    cloudflareSession.GetValue(j.key.Name),
			}
		}
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}

	for i, sk := range storageKeys {
		jobs <- job{
			index: i,
			key:   sk,
		}
	}

	close(jobs)
	wg.Wait()

	return entries
}

func (cloudflareSession *CloudflareSession) Size() (int, []cloudflare.StorageKey) {
	entries := cloudflareSession.GetAllKeys()
	return len(entries), entries
}

func (cloudflareSession *CloudflareSession) WriteEntry(entry models.Entry) (resp cloudflare.Response) {
	workersKVPairs := entryToWorkersKVPairs(entry)
	resp, err := cloudflareSession.api.WriteWorkersKVEntries(
		context.Background(),
		cloudflareSession.accountID,
		cloudflare.WriteWorkersKVEntriesParams{
			NamespaceID: cloudflareSession.namespaceID,
			KVs:         workersKVPairs,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func (cloudflareSession *CloudflareSession) InsertKVEntry(name string, value string, metadata string) (resp cloudflare.Response) {
	Metadata, _ := models.MetadataFromJSONString(metadata)
	newEntry := models.Entry{
		Name:     name,
		Metadata: Metadata,
		Value:    value,
	}
	return cloudflareSession.WriteEntry(newEntry)
}

func (cloudflareSession *CloudflareSession) WriteEntries(entries []models.Entry) {
	workersKVPairs := entriesToWorkersKVPairs(entries)
	resp, err := cloudflareSession.api.WriteWorkersKVEntries(
		context.Background(),
		cloudflareSession.accountID,
		cloudflare.WriteWorkersKVEntriesParams{
			NamespaceID: cloudflareSession.namespaceID,
			KVs:         workersKVPairs,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

func (cloudflareSession *CloudflareSession) DeleteKeyValue(key string) {
	resp, err := cloudflareSession.api.DeleteWorkersKVEntry(
		context.Background(),
		cloudflareSession.accountID,
		cloudflare.DeleteWorkersKVEntryParams{
			NamespaceID: cloudflareSession.namespaceID,
			Key:         key,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}

func (cloudflareSession *CloudflareSession) DeleteKeyValues(keys []string) {
	resp, err := cloudflareSession.api.DeleteWorkersKVEntries(
		context.Background(),
		cloudflareSession.accountID,
		cloudflare.DeleteWorkersKVEntriesParams{
			NamespaceID: cloudflareSession.namespaceID,
			Keys:        keys,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

func entryToWorkersKVPairs(entry models.Entry) []*cloudflare.WorkersKVPair {
	return entriesToWorkersKVPairs([]models.Entry{entry})
}

func entriesToWorkersKVPairs(entries []models.Entry) []*cloudflare.WorkersKVPair {
	var kvPairs []*cloudflare.WorkersKVPair
	for _, entry := range entries {
		var metadataMap map[string]interface{}
		metadataJSON, _ := json.Marshal(entry.Metadata)
		_ = json.Unmarshal(metadataJSON, &metadataMap)

		kvPairs = append(kvPairs, &cloudflare.WorkersKVPair{
			Key:      entry.Name,
			Value:    entry.Value,
			Metadata: metadataMap,
		})
	}
	return kvPairs
}
