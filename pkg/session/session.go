package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cdnmanager/pkg/models"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type CloudflareSession struct {
	api          *cloudflare.API
	account_id   *cloudflare.ResourceContainer
	namespace_id string
	domain       string
}

func NewCloudflareSession() *CloudflareSession {
	api, _ := cloudflare.New(os.Getenv("cloudflare_api_key"), os.Getenv("cloudflare_email"))

	cloudflareSession := &CloudflareSession{
		api:          api,
		account_id:   cloudflare.AccountIdentifier(os.Getenv("account_id")),
		namespace_id: os.Getenv("namespace_id"),
		domain:       os.Getenv("domain"),
	}

	return cloudflareSession
}

func (cloudflareSession *CloudflareSession) GetValue(key string) string {
	resp, _ := cloudflareSession.api.GetWorkersKV(context.Background(), cloudflareSession.account_id, cloudflare.GetWorkersKVParams{
		NamespaceID: cloudflareSession.namespace_id,
		Key:         key,
	})

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
	resp, _ := cloudflareSession.api.ListWorkersKVKeys(context.Background(), cloudflareSession.account_id, cloudflare.ListWorkersKVsParams{
		NamespaceID: cloudflareSession.namespace_id,
	})
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

func (cloudflareSession *CloudflareSession) Size() (int, []cloudflare.StorageKey) {
	entries := cloudflareSession.GetAllKeys()
	return len(entries), entries
}

func (cloudflareSession *CloudflareSession) WriteEntry(entry models.Entry) (resp cloudflare.Response) {
	workersKVPairs := entryToWorkersKVPairs(entry)
	resp, err := cloudflareSession.api.WriteWorkersKVEntries(context.Background(), cloudflareSession.account_id, cloudflare.WriteWorkersKVEntriesParams{
		NamespaceID: cloudflareSession.namespace_id,
		KVs:         workersKVPairs,
	})
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
	resp, err := cloudflareSession.api.WriteWorkersKVEntries(context.Background(), cloudflareSession.account_id, cloudflare.WriteWorkersKVEntriesParams{
		NamespaceID: cloudflareSession.namespace_id,
		KVs:         workersKVPairs,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

func (cloudflareSession *CloudflareSession) DeleteKeyValue(key string) {
	resp, err := cloudflareSession.api.DeleteWorkersKVEntry(context.Background(), cloudflareSession.account_id, cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: cloudflareSession.namespace_id,
		Key:         key,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}

func (cloudflareSession *CloudflareSession) DeleteKeyValues(keys []string) {
	resp, err := cloudflareSession.api.DeleteWorkersKVEntries(context.Background(), cloudflareSession.account_id, cloudflare.DeleteWorkersKVEntriesParams{
		NamespaceID: cloudflareSession.namespace_id,
		Keys:        keys,
	})
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
		json.Unmarshal(metadataJSON, &metadataMap)

		kvPairs = append(kvPairs, &cloudflare.WorkersKVPair{
			Key:      entry.Name,
			Value:    entry.Value,
			Metadata: metadataMap,
		})
	}
	return kvPairs
}
