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

func Test() []Entry {
	var entries []Entry
	entry := &Entry{
		Name:     "Name",
		Metadata: "JSON",
		Value:    "Value",
	}
	entries = append(entries, *entry)
	return entries
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

// Works
func GetValue(cloudflareSession *CloudflareSession, key string) string {
	resp, _ := cloudflareSession.api.GetWorkersKV(context.Background(), cloudflareSession.account_id, cloudflare.GetWorkersKVParams{
		NamespaceID: cloudflareSession.namespace_id,
		Key:         key,
	})

	return string(resp)
}

func GetAllValues(cloudflareSession *CloudflareSession) []string {
	storageKeys := GetAllKeys(cloudflareSession)
	var values []string
	for _, entry := range storageKeys {
		values = append(values, GetValue(cloudflareSession, entry.Name))
	}
	return values
}

func GetAllKeys(cloudflareSession *CloudflareSession) []cloudflare.StorageKey {
	resp, _ := cloudflareSession.api.ListWorkersKVKeys(context.Background(), cloudflareSession.account_id, cloudflare.ListWorkersKVsParams{
		NamespaceID: cloudflareSession.namespace_id,
	})
	return resp.Result
}

func GetAllEntries(cloudflareSession *CloudflareSession) []Entry {
	storageKeys := GetAllKeys(cloudflareSession)
	var entries []Entry
	for _, entry := range storageKeys {
		entry := Entry{
			Name:     entry.Name,
			Metadata: entry.Metadata,
			Value:    GetValue(cloudflareSession, entry.Name),
		}
		entries = append(entries, entry)
	}
	return entries
}

//	entries := []*cloudflare.WorkersKVPair{
//		{
//			Key:   "key1",
//			Value: "value1",
//			Metadata: "metadata1",
//		},
//		{
//			Key:   "key2",
//			Value: "value2",
//			Metadata: "metadata2",
//		},
//	}
func WriteEntries(cloudflareSession *CloudflareSession, entries []*cloudflare.WorkersKVPair) {
	resp, err := cloudflareSession.api.WriteWorkersKVEntries(context.Background(), cloudflareSession.account_id, cloudflare.WriteWorkersKVEntriesParams{
		NamespaceID: cloudflareSession.namespace_id,
		KVs:         entries,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

func DeleteKeyValue(cloudflareSession *CloudflareSession, key string) {
	resp, err := cloudflareSession.api.DeleteWorkersKVEntry(context.Background(), cloudflareSession.account_id, cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: cloudflareSession.namespace_id,
		Key:         key,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}

func DeleteKeyValues(cloudflareSession *CloudflareSession, keys []string) {
	resp, err := cloudflareSession.api.DeleteWorkersKVEntries(context.Background(), cloudflareSession.account_id, cloudflare.DeleteWorkersKVEntriesParams{
		NamespaceID: cloudflareSession.namespace_id,
		Keys:        keys,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
