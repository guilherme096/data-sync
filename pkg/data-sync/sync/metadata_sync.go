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
	SyncTables(catalogName, schemaName string) error
	SyncColumns(catalogName, schemaName, tableName string) error
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

func (s *metadataSync) SyncTables(catalogName, schemaName string) error {
	tables, err := s.discovery.DiscoverTables(catalogName, schemaName)
	if err != nil {
		return fmt.Errorf("failed to discover tables for schema '%s.%s': %w", catalogName, schemaName, err)
	}

	for _, table := range tables {
		if err := s.storage.UpsertTable(table); err != nil {
			log.Printf("Warning: failed to upsert table '%s.%s.%s': %v", catalogName, schemaName, table.Name, err)
			continue
		}
		log.Printf("Synced table: %s.%s.%s", catalogName, schemaName, table.Name)
	}

	return nil
}

func (s *metadataSync) SyncColumns(catalogName, schemaName, tableName string) error {
	columns, err := s.discovery.DiscoverColumns(catalogName, schemaName, tableName)
	if err != nil {
		return fmt.Errorf("failed to discover columns for table '%s.%s.%s': %w", catalogName, schemaName, tableName, err)
	}

	for _, column := range columns {
		if err := s.storage.UpsertColumn(column); err != nil {
			log.Printf("Warning: failed to upsert column '%s.%s.%s.%s': %v", catalogName, schemaName, tableName, column.Name, err)
			continue
		}
		log.Printf("Synced column: %s.%s.%s.%s", catalogName, schemaName, tableName, column.Name)
	}

	return nil
}

func (s *metadataSync) SyncAll() error {
	// Sync catalogs
	if err := s.SyncCatalogs(); err != nil {
		return fmt.Errorf("failed to sync catalogs: %w", err)
	}

	catalogs, err := s.storage.ListCatalogs()
	if err != nil {
		return fmt.Errorf("failed to list catalogs: %w", err)
	}

	totalSchemas := 0
	totalTables := 0
	totalColumns := 0

	// Sync schemas for each catalog
	for _, catalog := range catalogs {
		if err := s.SyncSchemas(catalog.Name); err != nil {
			log.Printf("Warning: failed to sync schemas for catalog '%s': %v", catalog.Name, err)
			continue
		}

		schemas, err := s.storage.ListSchemas(catalog.Name)
		if err != nil {
			log.Printf("Warning: failed to list schemas for catalog '%s': %v", catalog.Name, err)
			continue
		}
		totalSchemas += len(schemas)

		// Sync tables for each schema
		for _, schema := range schemas {
			if err := s.SyncTables(catalog.Name, schema.Name); err != nil {
				log.Printf("Warning: failed to sync tables for schema '%s.%s': %v", catalog.Name, schema.Name, err)
				continue
			}

			tables, err := s.storage.ListTables(catalog.Name, schema.Name)
			if err != nil {
				log.Printf("Warning: failed to list tables for schema '%s.%s': %v", catalog.Name, schema.Name, err)
				continue
			}
			totalTables += len(tables)

			// Sync columns for each table
			for _, table := range tables {
				if err := s.SyncColumns(catalog.Name, schema.Name, table.Name); err != nil {
					log.Printf("Warning: failed to sync columns for table '%s.%s.%s': %v", catalog.Name, schema.Name, table.Name, err)
					continue
				}

				columns, err := s.storage.ListColumns(catalog.Name, schema.Name, table.Name)
				if err != nil {
					log.Printf("Warning: failed to list columns for table '%s.%s.%s': %v", catalog.Name, schema.Name, table.Name, err)
					continue
				}
				totalColumns += len(columns)
			}
		}
	}

	log.Printf("Full sync completed: %d catalogs, %d schemas, %d tables, %d columns",
		len(catalogs), totalSchemas, totalTables, totalColumns)
	return nil
}
