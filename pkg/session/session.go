package session

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"cdnmanager/pkg/config"
	"cdnmanager/pkg/models"

	cloudflare "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/kv"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

const bulkGetChunkSize = 100
const workerCount = 18

type CloudflareSession struct {
	client      *cloudflare.Client
	accountID   string
	namespaceID string
	domain      string
}

type bulkGetRawEnvelope struct {
	Values map[string]bulkGetRawItem `json:"values"`
}

type bulkGetRawItem struct {
	Value    string         `json:"value"`
	Metadata map[string]any `json:"metadata"`
}

func NewCloudflareSession(cfg config.Config) (*CloudflareSession, error) {
	if !cfg.IsComplete() {
		return nil, fmt.Errorf("cloudflare session config is incomplete")
	}

	client := cloudflare.NewClient(
		option.WithAPIToken(cfg.CloudflareAPIToken),
	)

	return &CloudflareSession{
		client:      client,
		accountID:   cfg.AccountID,
		namespaceID: cfg.NamespaceID,
		domain:      cfg.Domain,
	}, nil
}

func (s *CloudflareSession) GetValue(key string) string {
	resp, err := s.client.KV.Namespaces.Values.Get(
		context.Background(),
		s.namespaceID,
		key,
		kv.NamespaceValueGetParams{
			AccountID: cloudflare.F(s.accountID),
		},
	)
	if err != nil {
		log.Printf("failed to get KV value for key %q: %v", key, err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read KV response body for key %q: %v", key, err)
		return ""
	}

	return string(body)
}

func (s *CloudflareSession) GetAllValues() []string {
	storageKeys := s.GetAllKeys()
	values := make([]string, 0, len(storageKeys))
	for _, entry := range storageKeys {
		values = append(values, s.GetValue(entry.Name))
	}
	return values
}

func (s *CloudflareSession) GetAllKeys() []kv.Key {
	pager := s.client.KV.Namespaces.Keys.ListAutoPaging(
		context.Background(),
		s.namespaceID,
		kv.NamespaceKeyListParams{
			AccountID: cloudflare.F(s.accountID),
		},
	)

	var keys []kv.Key
	for pager.Next() {
		keys = append(keys, pager.Current())
	}

	if err := pager.Err(); err != nil {
		log.Printf("failed to list KV keys: %v", err)
		return nil
	}

	return keys
}

func (s *CloudflareSession) GetAllEntries() []models.Entry {
	storageKeys := s.GetAllKeys()
	return s.GetAllEntriesFromKeys(storageKeys)
}

func (s *CloudflareSession) GetAllEntriesFromKeys(storageKeys []kv.Key) []models.Entry {
	if len(storageKeys) == 0 {
		return []models.Entry{}
	}

	entries := make([]models.Entry, len(storageKeys))

	type job struct {
		index int
		key   kv.Key
	}

	jobs := make(chan job, len(storageKeys))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for j := range jobs {
			metadata := metadataFromAny(j.key.Metadata)

			entries[j.index] = models.Entry{
				Name:     j.key.Name,
				Metadata: metadata,
				Value:    s.GetValue(j.key.Name),
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

func (s *CloudflareSession) GetAllEntriesBulk() []models.Entry {
	keys := s.GetAllKeys()
	if len(keys) == 0 {
		return []models.Entry{}
	}

	keyNames := make([]string, 0, len(keys))
	for _, k := range keys {
		keyNames = append(keyNames, k.Name)
	}

	var entries []models.Entry

	for start := 0; start < len(keyNames); start += bulkGetChunkSize {
		end := start + bulkGetChunkSize
		if end > len(keyNames) {
			end = len(keyNames)
		}

		chunk := keyNames[start:end]

		resp, err := s.client.KV.Namespaces.BulkGet(
			context.Background(),
			s.namespaceID,
			kv.NamespaceBulkGetParams{
				AccountID:    cloudflare.F(s.accountID),
				Keys:         cloudflare.F(chunk),
				WithMetadata: cloudflare.F(true),
			},
		)
		if err != nil {
			log.Printf("failed bulk get for keys %d:%d: %v", start, end, err)
			continue
		}

		chunkEntries, err := normalizeBulkGetResponse(resp)
		if err != nil {
			log.Printf("failed to normalize bulk get for keys %d:%d: %v", start, end, err)
			continue
		}

		entries = append(entries, chunkEntries...)
	}

	return entries
}

func (s *CloudflareSession) Size() (int, []kv.Key) {
	entries := s.GetAllKeys()
	return len(entries), entries
}

func (s *CloudflareSession) WriteEntry(entry models.Entry) {
	s.WriteEntries([]models.Entry{entry})
}

func (s *CloudflareSession) InsertKVEntry(name string, value string, metadata string) {
	meta, _ := models.MetadataFromJSONString(metadata)
	newEntry := models.Entry{
		Name:     name,
		Metadata: meta,
		Value:    value,
	}
	s.WriteEntry(newEntry)
}

func (s *CloudflareSession) WriteEntries(entries []models.Entry) {
	if len(entries) == 0 {
		return
	}

	kvs := entriesToBulkUpdateBodies(entries)

	_, err := s.client.KV.Namespaces.BulkUpdate(
		context.Background(),
		s.namespaceID,
		kv.NamespaceBulkUpdateParams{
			AccountID: cloudflare.F(s.accountID),
			Body:      kvs,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *CloudflareSession) DeleteKeyValue(key string) {
	_, err := s.client.KV.Namespaces.Values.Delete(
		context.Background(),
		s.namespaceID,
		key,
		kv.NamespaceValueDeleteParams{
			AccountID: cloudflare.F(s.accountID),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *CloudflareSession) DeleteKeyValues(keys []string) {
	if len(keys) == 0 {
		return
	}

	_, err := s.client.KV.Namespaces.BulkDelete(
		context.Background(),
		s.namespaceID,
		kv.NamespaceBulkDeleteParams{
			AccountID: cloudflare.F(s.accountID),
			Body:      keys,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func entriesToBulkUpdateBodies(entries []models.Entry) []kv.NamespaceBulkUpdateParamsBody {
	bodies := make([]kv.NamespaceBulkUpdateParamsBody, 0, len(entries))

	for _, entry := range entries {
		metadataMap := metadataToMap(entry.Metadata)

		bodies = append(bodies, kv.NamespaceBulkUpdateParamsBody{
			Key:      cloudflare.F(entry.Name),
			Value:    cloudflare.F(entry.Value),
			Metadata: cloudflare.F(any(metadataMap)),
		})
	}

	return bodies
}

func metadataToMap(metadata models.Metadata) map[string]any {
	var metadataMap map[string]any
	metadataJSON, _ := json.Marshal(metadata)
	_ = json.Unmarshal(metadataJSON, &metadataMap)
	return metadataMap
}

func metadataFromAny(raw any) models.Metadata {
	if raw == nil {
		return models.Metadata{}
	}

	metadataJSON, err := json.Marshal(raw)
	if err != nil {
		return models.Metadata{}
	}

	var metadata models.Metadata
	if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
		return models.Metadata{}
	}

	return metadata
}

func normalizeBulkGetResponse(resp *kv.NamespaceBulkGetResponse) ([]models.Entry, error) {
	if resp == nil {
		return nil, fmt.Errorf("nil bulk get response")
	}

	raw := resp.JSON.RawJSON()
	if raw == "" {
		return nil, fmt.Errorf("empty raw JSON in bulk get response")
	}

	var envelope bulkGetRawEnvelope
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal raw bulk get JSON: %w", err)
	}

	entries := make([]models.Entry, 0, len(envelope.Values))
	for key, item := range envelope.Values {
		entries = append(entries, models.Entry{
			Name:     key,
			Value:    item.Value,
			Metadata: metadataFromAny(item.Metadata),
		})
	}

	return entries, nil
}
