package storage

import (
	"fmt"
	"sync"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
)

type MemoryMetadataStorage struct {
	mu       sync.RWMutex
	catalogs map[string]*models.Catalog
	schemas  map[string]map[string]*models.Schema
}

func NewMemoryMetadataStorage() *MemoryMetadataStorage {
	return &MemoryMetadataStorage{
		catalogs: make(map[string]*models.Catalog),
		schemas:  make(map[string]map[string]*models.Schema),
	}
}

func (m *MemoryMetadataStorage) CreateCatalog(catalog *models.Catalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if catalog.Name == "" {
		return fmt.Errorf("catalog name cannot be empty")
	}

	if _, exists := m.catalogs[catalog.Name]; exists {
		return fmt.Errorf("catalog '%s' already exists", catalog.Name)
	}

	m.catalogs[catalog.Name] = catalog
	return nil
}

func (m *MemoryMetadataStorage) GetCatalog(name string) (*models.Catalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalog, exists := m.catalogs[name]
	if !exists {
		return nil, fmt.Errorf("catalog '%s' not found", name)
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

	if schema.CatalogName == "" || schema.Name == "" {
		return fmt.Errorf("catalog name and schema name cannot be empty")
	}

	// Check if catalog exists
	if _, exists := m.catalogs[schema.CatalogName]; !exists {
		return fmt.Errorf("catalog '%s' not found", schema.CatalogName)
	}

	// Initialize nested map if needed
	if m.schemas[schema.CatalogName] == nil {
		m.schemas[schema.CatalogName] = make(map[string]*models.Schema)
	}

	// Check if schema already exists
	if _, exists := m.schemas[schema.CatalogName][schema.Name]; exists {
		return fmt.Errorf("schema '%s' already exists in catalog '%s'", schema.Name, schema.CatalogName)
	}

	m.schemas[schema.CatalogName][schema.Name] = schema
	return nil
}

func (m *MemoryMetadataStorage) GetSchema(catalogName, schemaName string) (*models.Schema, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogSchemas, exists := m.schemas[catalogName]
	if !exists {
		return nil, fmt.Errorf("no schemas found for catalog '%s'", catalogName)
	}

	schema, exists := catalogSchemas[schemaName]
	if !exists {
		return nil, fmt.Errorf("schema '%s' not found in catalog '%s'", schemaName, catalogName)
	}

	return schema, nil
}

func (m *MemoryMetadataStorage) ListSchemas(catalogName string) ([]*models.Schema, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogSchemas, exists := m.schemas[catalogName]
	if !exists {
		return []*models.Schema{}, nil
	}

	schemas := make([]*models.Schema, 0, len(catalogSchemas))
	for _, schema := range catalogSchemas {
		schemas = append(schemas, schema)
	}

	return schemas, nil
}
