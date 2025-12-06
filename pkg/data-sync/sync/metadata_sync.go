package sync

import (
	"fmt"
	"log"

	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type MetadataSync interface {
	SyncCatalogs() error
	SyncSchemas(catalogName string) error
	SyncAll() error
}

type metadataSync struct {
	discovery discovery.MetadataDiscovery
	storage   storage.MetadataStorage
}

func NewMetadataSync(discovery discovery.MetadataDiscovery, storage storage.MetadataStorage) MetadataSync {
	return &metadataSync{
		discovery: discovery,
		storage:   storage,
	}
}

func (s *metadataSync) SyncCatalogs() error {
	catalogs, err := s.discovery.DiscoverCatalogs()
	if err != nil {
		return fmt.Errorf("failed to discover catalogs: %w", err)
	}

	for _, catalog := range catalogs {
		if err := s.storage.UpsertCatalog(catalog); err != nil {
			log.Printf("Warning: failed to upsert catalog '%s': %v", catalog.Name, err)
			continue
		}
		log.Printf("Synced catalog: %s", catalog.Name)
	}

	return nil
}

func (s *metadataSync) SyncSchemas(catalogName string) error {
	schemas, err := s.discovery.DiscoverSchemas(catalogName)
	if err != nil {
		return fmt.Errorf("failed to discover schemas for catalog '%s': %w", catalogName, err)
	}

	for _, schema := range schemas {
		if err := s.storage.UpsertSchema(schema); err != nil {
			log.Printf("Warning: failed to upsert schema '%s.%s': %v", catalogName, schema.Name, err)
			continue
		}
		log.Printf("Synced schema: %s.%s", catalogName, schema.Name)
	}

	return nil
}

func (s *metadataSync) SyncAll() error {
	if err := s.SyncCatalogs(); err != nil {
		return fmt.Errorf("failed to sync catalogs: %w", err)
	}

	catalogs, err := s.storage.ListCatalogs()
	if err != nil {
		return fmt.Errorf("failed to list catalogs: %w", err)
	}

	for _, catalog := range catalogs {
		if err := s.SyncSchemas(catalog.Name); err != nil {
			log.Printf("Warning: failed to sync schemas for catalog '%s': %v", catalog.Name, err)
			continue
		}
	}

	log.Printf("Full sync completed: %d catalogs synced", len(catalogs))
	return nil
}
