package session

import (
	"context"
	"fmt"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Entry struct {
	Name     string
	Metadata interface{}
	Value    string
}

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

func (cloudflareSession *CloudflareSession) GetAllEntries() []Entry {
	storageKeys := cloudflareSession.GetAllKeys()
	var entries []Entry
	for _, entry := range storageKeys {
		entry := Entry{
			Name:     entry.Name,
			Metadata: entry.Metadata,
			Value:    cloudflareSession.GetValue(entry.Name),
		}
		entries = append(entries, entry)
	}
	return entries
}

func (cloudflareSession *CloudflareSession) GetAllEntriesFromKeys(storageKeys []cloudflare.StorageKey) []Entry {
	var entries []Entry
	for _, entry := range storageKeys {
		entry := Entry{
			Name:     entry.Name,
			Metadata: entry.Metadata,
			Value:    cloudflareSession.GetValue(entry.Name),
		}
		entries = append(entries, entry)
	}
	return entries
}

func (cloudflareSession *CloudflareSession) Size() (int, []cloudflare.StorageKey) {
	entries := cloudflareSession.GetAllKeys()
	return len(entries), entries
}

func (cloudflareSession *CloudflareSession) WriteEntry(entry Entry) {
	workersKVPairs := entryToWorkersKVPairs(entry)
	resp, err := cloudflareSession.api.WriteWorkersKVEntries(context.Background(), cloudflareSession.account_id, cloudflare.WriteWorkersKVEntriesParams{
		NamespaceID: cloudflareSession.namespace_id,
		KVs:         workersKVPairs,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

func (cloudflareSession *CloudflareSession) InsertEntry(name string, value string, metadata string) {
	newEntry := Entry{
		Name:     name,
		Metadata: metadata,
		Value:    value,
	}
	cloudflareSession.WriteEntry(newEntry)
}

func (cloudflareSession *CloudflareSession) WriteEntries(entries []Entry) {
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

func entryToWorkersKVPairs(entry Entry) []*cloudflare.WorkersKVPair {
	return entriesToWorkersKVPairs([]Entry{entry})
}

func entriesToWorkersKVPairs(entries []Entry) []*cloudflare.WorkersKVPair {
	var kvPairs []*cloudflare.WorkersKVPair
	for _, entry := range entries {
		kvPairs = append(kvPairs, &cloudflare.WorkersKVPair{
			Key:      entry.Name,
			Value:    entry.Value,
			Metadata: entry.Metadata,
		})
	}
	return kvPairs
}
