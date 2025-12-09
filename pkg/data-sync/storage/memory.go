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
	tables   map[string]map[string]map[string]*models.Table
	columns  map[string]map[string]map[string]map[string]*models.Column
}

func NewMemoryMetadataStorage() *MemoryMetadataStorage {
	return &MemoryMetadataStorage{
		catalogs: make(map[string]*models.Catalog),
		schemas:  make(map[string]map[string]*models.Schema),
		tables:   make(map[string]map[string]map[string]*models.Table),
		columns:  make(map[string]map[string]map[string]map[string]*models.Column),
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
// Table Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateTable(table *models.Table) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if table.CatalogName == "" || table.SchemaName == "" || table.Name == "" {
		return fmt.Errorf("catalog name, schema name, and table name cannot be empty")
	}

	// Check if schema exists
	if _, exists := m.schemas[table.CatalogName]; !exists {
		return fmt.Errorf("catalog '%s' not found", table.CatalogName)
	}
	if _, exists := m.schemas[table.CatalogName][table.SchemaName]; !exists {
		return fmt.Errorf("schema '%s' not found in catalog '%s'", table.SchemaName, table.CatalogName)
	}

	// Initialize nested maps if needed
	if m.tables[table.CatalogName] == nil {
		m.tables[table.CatalogName] = make(map[string]map[string]*models.Table)
	}
	if m.tables[table.CatalogName][table.SchemaName] == nil {
		m.tables[table.CatalogName][table.SchemaName] = make(map[string]*models.Table)
	}

	// Check if table already exists
	if _, exists := m.tables[table.CatalogName][table.SchemaName][table.Name]; exists {
		return fmt.Errorf("table '%s' already exists in schema '%s.%s'", table.Name, table.CatalogName, table.SchemaName)
	}

	m.tables[table.CatalogName][table.SchemaName][table.Name] = table
	return nil
}

func (m *MemoryMetadataStorage) GetTable(catalogName, schemaName, tableName string) (*models.Table, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogTables, exists := m.tables[catalogName]
	if !exists {
		return nil, fmt.Errorf("no tables found for catalog '%s'", catalogName)
	}

	schemaTables, exists := catalogTables[schemaName]
	if !exists {
		return nil, fmt.Errorf("no tables found for schema '%s.%s'", catalogName, schemaName)
	}

	table, exists := schemaTables[tableName]
	if !exists {
		return nil, fmt.Errorf("table '%s' not found in schema '%s.%s'", tableName, catalogName, schemaName)
	}

	return table, nil
}

func (m *MemoryMetadataStorage) ListTables(catalogName, schemaName string) ([]*models.Table, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogTables, exists := m.tables[catalogName]
	if !exists {
		return []*models.Table{}, nil
	}

	schemaTables, exists := catalogTables[schemaName]
	if !exists {
		return []*models.Table{}, nil
	}

	tables := make([]*models.Table, 0, len(schemaTables))
	for _, table := range schemaTables {
		tables = append(tables, table)
	}

	return tables, nil
}

func (m *MemoryMetadataStorage) UpdateTable(table *models.Table) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if table.CatalogName == "" || table.SchemaName == "" || table.Name == "" {
		return fmt.Errorf("catalog name, schema name, and table name cannot be empty")
	}

	catalogTables, exists := m.tables[table.CatalogName]
	if !exists {
		return fmt.Errorf("no tables found for catalog '%s'", table.CatalogName)
	}

	schemaTables, exists := catalogTables[table.SchemaName]
	if !exists {
		return fmt.Errorf("no tables found for schema '%s.%s'", table.CatalogName, table.SchemaName)
	}

	if _, exists := schemaTables[table.Name]; !exists {
		return fmt.Errorf("table '%s' not found in schema '%s.%s'", table.Name, table.CatalogName, table.SchemaName)
	}

	m.tables[table.CatalogName][table.SchemaName][table.Name] = table
	return nil
}

func (m *MemoryMetadataStorage) UpsertTable(table *models.Table) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if table.CatalogName == "" || table.SchemaName == "" || table.Name == "" {
		return fmt.Errorf("catalog name, schema name, and table name cannot be empty")
	}

	// Check if schema exists
	if _, exists := m.schemas[table.CatalogName]; !exists {
		return fmt.Errorf("catalog '%s' not found", table.CatalogName)
	}
	if _, exists := m.schemas[table.CatalogName][table.SchemaName]; !exists {
		return fmt.Errorf("schema '%s' not found in catalog '%s'", table.SchemaName, table.CatalogName)
	}

	// Initialize nested maps if needed
	if m.tables[table.CatalogName] == nil {
		m.tables[table.CatalogName] = make(map[string]map[string]*models.Table)
	}
	if m.tables[table.CatalogName][table.SchemaName] == nil {
		m.tables[table.CatalogName][table.SchemaName] = make(map[string]*models.Table)
	}

	m.tables[table.CatalogName][table.SchemaName][table.Name] = table
	return nil
}

// ============================================================================
// Column Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateColumn(column *models.Column) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if column.CatalogName == "" || column.SchemaName == "" || column.TableName == "" || column.Name == "" {
		return fmt.Errorf("catalog name, schema name, table name, and column name cannot be empty")
	}

	// Check if table exists
	if _, exists := m.tables[column.CatalogName]; !exists {
		return fmt.Errorf("catalog '%s' not found", column.CatalogName)
	}
	if _, exists := m.tables[column.CatalogName][column.SchemaName]; !exists {
		return fmt.Errorf("schema '%s' not found in catalog '%s'", column.SchemaName, column.CatalogName)
	}
	if _, exists := m.tables[column.CatalogName][column.SchemaName][column.TableName]; !exists {
		return fmt.Errorf("table '%s' not found in schema '%s.%s'", column.TableName, column.CatalogName, column.SchemaName)
	}

	// Initialize nested maps if needed
	if m.columns[column.CatalogName] == nil {
		m.columns[column.CatalogName] = make(map[string]map[string]map[string]*models.Column)
	}
	if m.columns[column.CatalogName][column.SchemaName] == nil {
		m.columns[column.CatalogName][column.SchemaName] = make(map[string]map[string]*models.Column)
	}
	if m.columns[column.CatalogName][column.SchemaName][column.TableName] == nil {
		m.columns[column.CatalogName][column.SchemaName][column.TableName] = make(map[string]*models.Column)
	}

	// Check if column already exists
	if _, exists := m.columns[column.CatalogName][column.SchemaName][column.TableName][column.Name]; exists {
		return fmt.Errorf("column '%s' already exists in table '%s.%s.%s'", column.Name, column.CatalogName, column.SchemaName, column.TableName)
	}

	m.columns[column.CatalogName][column.SchemaName][column.TableName][column.Name] = column
	return nil
}

func (m *MemoryMetadataStorage) GetColumn(catalogName, schemaName, tableName, columnName string) (*models.Column, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogColumns, exists := m.columns[catalogName]
	if !exists {
		return nil, fmt.Errorf("no columns found for catalog '%s'", catalogName)
	}

	schemaColumns, exists := catalogColumns[schemaName]
	if !exists {
		return nil, fmt.Errorf("no columns found for schema '%s.%s'", catalogName, schemaName)
	}

	tableColumns, exists := schemaColumns[tableName]
	if !exists {
		return nil, fmt.Errorf("no columns found for table '%s.%s.%s'", catalogName, schemaName, tableName)
	}

	column, exists := tableColumns[columnName]
	if !exists {
		return nil, fmt.Errorf("column '%s' not found in table '%s.%s.%s'", columnName, catalogName, schemaName, tableName)
	}

	return column, nil
}

func (m *MemoryMetadataStorage) ListColumns(catalogName, schemaName, tableName string) ([]*models.Column, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalogColumns, exists := m.columns[catalogName]
	if !exists {
		return []*models.Column{}, nil
	}

	schemaColumns, exists := catalogColumns[schemaName]
	if !exists {
		return []*models.Column{}, nil
	}

	tableColumns, exists := schemaColumns[tableName]
	if !exists {
		return []*models.Column{}, nil
	}

	columns := make([]*models.Column, 0, len(tableColumns))
	for _, column := range tableColumns {
		columns = append(columns, column)
	}

	return columns, nil
}

func (m *MemoryMetadataStorage) UpdateColumn(column *models.Column) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if column.CatalogName == "" || column.SchemaName == "" || column.TableName == "" || column.Name == "" {
		return fmt.Errorf("catalog name, schema name, table name, and column name cannot be empty")
	}

	catalogColumns, exists := m.columns[column.CatalogName]
	if !exists {
		return fmt.Errorf("no columns found for catalog '%s'", column.CatalogName)
	}

	schemaColumns, exists := catalogColumns[column.SchemaName]
	if !exists {
		return fmt.Errorf("no columns found for schema '%s.%s'", column.CatalogName, column.SchemaName)
	}

	tableColumns, exists := schemaColumns[column.TableName]
	if !exists {
		return fmt.Errorf("no columns found for table '%s.%s.%s'", column.CatalogName, column.SchemaName, column.TableName)
	}

	if _, exists := tableColumns[column.Name]; !exists {
		return fmt.Errorf("column '%s' not found in table '%s.%s.%s'", column.Name, column.CatalogName, column.SchemaName, column.TableName)
	}

	m.columns[column.CatalogName][column.SchemaName][column.TableName][column.Name] = column
	return nil
}

func (m *MemoryMetadataStorage) UpsertColumn(column *models.Column) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if column.CatalogName == "" || column.SchemaName == "" || column.TableName == "" || column.Name == "" {
		return fmt.Errorf("catalog name, schema name, table name, and column name cannot be empty")
	}

	// Check if table exists
	if _, exists := m.tables[column.CatalogName]; !exists {
		return fmt.Errorf("catalog '%s' not found", column.CatalogName)
	}
	if _, exists := m.tables[column.CatalogName][column.SchemaName]; !exists {
		return fmt.Errorf("schema '%s' not found in catalog '%s'", column.SchemaName, column.CatalogName)
	}
	if _, exists := m.tables[column.CatalogName][column.SchemaName][column.TableName]; !exists {
		return fmt.Errorf("table '%s' not found in schema '%s.%s'", column.TableName, column.CatalogName, column.SchemaName)
	}

	// Initialize nested maps if needed
	if m.columns[column.CatalogName] == nil {
		m.columns[column.CatalogName] = make(map[string]map[string]map[string]*models.Column)
	}
	if m.columns[column.CatalogName][column.SchemaName] == nil {
		m.columns[column.CatalogName][column.SchemaName] = make(map[string]map[string]*models.Column)
	}
	if m.columns[column.CatalogName][column.SchemaName][column.TableName] == nil {
		m.columns[column.CatalogName][column.SchemaName][column.TableName] = make(map[string]*models.Column)
	}

	m.columns[column.CatalogName][column.SchemaName][column.TableName][column.Name] = column
	return nil
}
