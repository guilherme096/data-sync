package storage

import (
	"fmt"
	"sync"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
)

type MemoryMetadataStorage struct {
	mu       sync.RWMutex
	catalogs map[int]*models.Catalog
	schemas  map[int]*models.Schema
	nextCatalogID int
	nextSchemaID  int
}

func NewMemoryMetadataStorage() *MemoryMetadataStorage {
	return &MemoryMetadataStorage{
		catalogs: make(map[int]*models.Catalog),
		schemas:  make(map[int]*models.Schema),
		nextCatalogID: 1,
		nextSchemaID:  1,
	}
}

func (m *MemoryMetadataStorage) CreateCatalog(catalog *models.Catalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	catalog.ID = m.nextCatalogID
	m.catalogs[catalog.ID] = catalog
	m.nextCatalogID++

	return nil
}

func (m *MemoryMetadataStorage) GetCatalog(id int) (*models.Catalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalog, exists := m.catalogs[id]
	if !exists {
		return nil, fmt.Errorf("catalog with id %d not found", id)
	}

	return catalog, nil
}

func (m *MemoryMetadataStorage) ListCatalogs() ([]*models.Catalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogs := make([]*models.Catalog, 0, len(m.catalogs))
	for _, catalog := range m.catalogs {
		catalogs = append(catalogs, catalog)
	}

	return catalogs, nil
}

func (m *MemoryMetadataStorage) CreateSchema(schema *models.Schema) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if catalog exists
	if _, exists := m.catalogs[schema.CatalogID]; !exists {
		return fmt.Errorf("catalog with id %d not found", schema.CatalogID)
	}

	schema.ID = m.nextSchemaID
	m.schemas[schema.ID] = schema
	m.nextSchemaID++

	return nil
}

func (m *MemoryMetadataStorage) GetSchema(id int) (*models.Schema, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	schema, exists := m.schemas[id]
	if !exists {
		return nil, fmt.Errorf("schema with id %d not found", id)
	}

	return schema, nil
}

func (m *MemoryMetadataStorage) ListSchemas(catalogID int) ([]*models.Schema, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	schemas := make([]*models.Schema, 0)
	for _, schema := range m.schemas {
		if schema.CatalogID == catalogID {
			schemas = append(schemas, schema)
		}
	}

	return schemas, nil
}
