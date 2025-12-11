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

	// Global metadata
	globalTables   map[string]*models.GlobalTable
	globalColumns  map[string]map[string]*models.GlobalColumn        // globalTable -> columnName -> column
	tableMappings  map[string][]*models.TableMapping                 // globalTable -> mappings
	columnMappings map[string]map[string][]*models.ColumnMapping     // globalTable -> columnName -> mappings
}

func NewMemoryMetadataStorage() *MemoryMetadataStorage {
	return &MemoryMetadataStorage{
		catalogs: make(map[string]*models.Catalog),
		schemas:  make(map[string]map[string]*models.Schema),

		globalTables:   make(map[string]*models.GlobalTable),
		globalColumns:  make(map[string]map[string]*models.GlobalColumn),
		tableMappings:  make(map[string][]*models.TableMapping),
		columnMappings: make(map[string]map[string][]*models.ColumnMapping),
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

func (m *MemoryMetadataStorage) UpdateCatalog(catalog *models.Catalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if catalog.Name == "" {
		return fmt.Errorf("catalog name cannot be empty")
	}

	if _, exists := m.catalogs[catalog.Name]; !exists {
		return fmt.Errorf("catalog '%s' not found", catalog.Name)
	}

	m.catalogs[catalog.Name] = catalog
	return nil
}

func (m *MemoryMetadataStorage) UpsertCatalog(catalog *models.Catalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if catalog.Name == "" {
		return fmt.Errorf("catalog name cannot be empty")
	}

	m.catalogs[catalog.Name] = catalog
	return nil
}

func (m *MemoryMetadataStorage) UpdateSchema(schema *models.Schema) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if schema.CatalogName == "" || schema.Name == "" {
		return fmt.Errorf("catalog name and schema name cannot be empty")
	}

	catalogSchemas, exists := m.schemas[schema.CatalogName]
	if !exists {
		return fmt.Errorf("no schemas found for catalog '%s'", schema.CatalogName)
	}

	if _, exists := catalogSchemas[schema.Name]; !exists {
		return fmt.Errorf("schema '%s' not found in catalog '%s'", schema.Name, schema.CatalogName)
	}

	m.schemas[schema.CatalogName][schema.Name] = schema
	return nil
}

func (m *MemoryMetadataStorage) UpsertSchema(schema *models.Schema) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if schema.CatalogName == "" || schema.Name == "" {
		return fmt.Errorf("catalog name and schema name cannot be empty")
	}

	if _, exists := m.catalogs[schema.CatalogName]; !exists {
		return fmt.Errorf("catalog '%s' not found", schema.CatalogName)
	}

	// initialize nested map if needed
	if m.schemas[schema.CatalogName] == nil {
		m.schemas[schema.CatalogName] = make(map[string]*models.Schema)
	}

	m.schemas[schema.CatalogName][schema.Name] = schema
	return nil
}

// ============================================================================
// Global Table Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateGlobalTable(table *models.GlobalTable) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if table.Name == "" {
		return fmt.Errorf("global table name cannot be empty")
	}

	if _, exists := m.globalTables[table.Name]; exists {
		return fmt.Errorf("global table '%s' already exists", table.Name)
	}

	m.globalTables[table.Name] = table
	return nil
}

func (m *MemoryMetadataStorage) GetGlobalTable(name string) (*models.GlobalTable, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	table, exists := m.globalTables[name]
	if !exists {
		return nil, fmt.Errorf("global table '%s' not found", name)
	}

	return table, nil
}

func (m *MemoryMetadataStorage) ListGlobalTables() ([]*models.GlobalTable, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tables := make([]*models.GlobalTable, 0, len(m.globalTables))
	for _, table := range m.globalTables {
		tables = append(tables, table)
	}

	return tables, nil
}

func (m *MemoryMetadataStorage) DeleteGlobalTable(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.globalTables[name]; !exists {
		return fmt.Errorf("global table '%s' not found", name)
	}

	// Delete the table and all related data
	delete(m.globalTables, name)
	delete(m.globalColumns, name)
	delete(m.tableMappings, name)
	delete(m.columnMappings, name)

	return nil
}

// ============================================================================
// Global Column Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateGlobalColumn(column *models.GlobalColumn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if column.GlobalTableName == "" || column.Name == "" {
		return fmt.Errorf("global table name and column name cannot be empty")
	}

	// Check if global table exists
	if _, exists := m.globalTables[column.GlobalTableName]; !exists {
		return fmt.Errorf("global table '%s' not found", column.GlobalTableName)
	}

	// Initialize nested map if needed
	if m.globalColumns[column.GlobalTableName] == nil {
		m.globalColumns[column.GlobalTableName] = make(map[string]*models.GlobalColumn)
	}

	// Check if column already exists
	if _, exists := m.globalColumns[column.GlobalTableName][column.Name]; exists {
		return fmt.Errorf("global column '%s' already exists in table '%s'", column.Name, column.GlobalTableName)
	}

	m.globalColumns[column.GlobalTableName][column.Name] = column
	return nil
}

func (m *MemoryMetadataStorage) ListGlobalColumns(globalTableName string) ([]*models.GlobalColumn, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	columns, exists := m.globalColumns[globalTableName]
	if !exists {
		return []*models.GlobalColumn{}, nil
	}

	result := make([]*models.GlobalColumn, 0, len(columns))
	for _, column := range columns {
		result = append(result, column)
	}

	return result, nil
}

func (m *MemoryMetadataStorage) DeleteGlobalColumn(globalTableName, columnName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	columns, exists := m.globalColumns[globalTableName]
	if !exists {
		return fmt.Errorf("no columns found for global table '%s'", globalTableName)
	}

	if _, exists := columns[columnName]; !exists {
		return fmt.Errorf("global column '%s' not found in table '%s'", columnName, globalTableName)
	}

	// Delete the column and its mappings
	delete(columns, columnName)
	if tableColumnMappings, exists := m.columnMappings[globalTableName]; exists {
		delete(tableColumnMappings, columnName)
	}

	return nil
}

// ============================================================================
// Table Mapping Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateTableMapping(mapping *models.TableMapping) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mapping.GlobalTableName == "" || mapping.CatalogName == "" || mapping.SchemaName == "" || mapping.TableName == "" {
		return fmt.Errorf("all fields in table mapping must be non-empty")
	}

	// Check if global table exists
	if _, exists := m.globalTables[mapping.GlobalTableName]; !exists {
		return fmt.Errorf("global table '%s' not found", mapping.GlobalTableName)
	}

	// Check for duplicate mapping
	existingMappings := m.tableMappings[mapping.GlobalTableName]
	for _, existing := range existingMappings {
		if existing.CatalogName == mapping.CatalogName &&
			existing.SchemaName == mapping.SchemaName &&
			existing.TableName == mapping.TableName {
			return fmt.Errorf("table mapping already exists")
		}
	}

	m.tableMappings[mapping.GlobalTableName] = append(m.tableMappings[mapping.GlobalTableName], mapping)
	return nil
}

func (m *MemoryMetadataStorage) ListTableMappings(globalTableName string) ([]*models.TableMapping, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mappings, exists := m.tableMappings[globalTableName]
	if !exists {
		return []*models.TableMapping{}, nil
	}

	return mappings, nil
}

func (m *MemoryMetadataStorage) DeleteTableMapping(globalTableName, catalog, schema, table string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mappings, exists := m.tableMappings[globalTableName]
	if !exists {
		return fmt.Errorf("no table mappings found for global table '%s'", globalTableName)
	}

	// Find and remove the mapping
	for i, mapping := range mappings {
		if mapping.CatalogName == catalog && mapping.SchemaName == schema && mapping.TableName == table {
			m.tableMappings[globalTableName] = append(mappings[:i], mappings[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("table mapping not found")
}

// ============================================================================
// Column Mapping Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateColumnMapping(mapping *models.ColumnMapping) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mapping.GlobalTableName == "" || mapping.GlobalColumnName == "" ||
		mapping.CatalogName == "" || mapping.SchemaName == "" ||
		mapping.TableName == "" || mapping.ColumnName == "" {
		return fmt.Errorf("all fields in column mapping must be non-empty")
	}

	// Check if global column exists
	if columns, exists := m.globalColumns[mapping.GlobalTableName]; !exists || columns[mapping.GlobalColumnName] == nil {
		return fmt.Errorf("global column '%s.%s' not found", mapping.GlobalTableName, mapping.GlobalColumnName)
	}

	// Initialize nested maps if needed
	if m.columnMappings[mapping.GlobalTableName] == nil {
		m.columnMappings[mapping.GlobalTableName] = make(map[string][]*models.ColumnMapping)
	}

	// Check for duplicate mapping
	existingMappings := m.columnMappings[mapping.GlobalTableName][mapping.GlobalColumnName]
	for _, existing := range existingMappings {
		if existing.CatalogName == mapping.CatalogName &&
			existing.SchemaName == mapping.SchemaName &&
			existing.TableName == mapping.TableName &&
			existing.ColumnName == mapping.ColumnName {
			return fmt.Errorf("column mapping already exists")
		}
	}

	m.columnMappings[mapping.GlobalTableName][mapping.GlobalColumnName] = append(
		m.columnMappings[mapping.GlobalTableName][mapping.GlobalColumnName],
		mapping,
	)

	return nil
}

func (m *MemoryMetadataStorage) ListColumnMappings(globalTableName, globalColumnName string) ([]*models.ColumnMapping, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tableMappings, exists := m.columnMappings[globalTableName]
	if !exists {
		return []*models.ColumnMapping{}, nil
	}

	columnMappings, exists := tableMappings[globalColumnName]
	if !exists {
		return []*models.ColumnMapping{}, nil
	}

	return columnMappings, nil
}

func (m *MemoryMetadataStorage) DeleteColumnMapping(globalTableName, globalColumnName, catalog, schema, table, column string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tableMappings, exists := m.columnMappings[globalTableName]
	if !exists {
		return fmt.Errorf("no column mappings found for global table '%s'", globalTableName)
	}

	mappings, exists := tableMappings[globalColumnName]
	if !exists {
		return fmt.Errorf("no column mappings found for global column '%s.%s'", globalTableName, globalColumnName)
	}

	// Find and remove the mapping
	for i, mapping := range mappings {
		if mapping.CatalogName == catalog &&
			mapping.SchemaName == schema &&
			mapping.TableName == table &&
			mapping.ColumnName == column {
			m.columnMappings[globalTableName][globalColumnName] = append(mappings[:i], mappings[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("column mapping not found")
}
