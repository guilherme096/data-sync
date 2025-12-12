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

	// Local metadata (discovered from data sources)
	tables  map[string]map[string]map[string]*models.Table                    // catalog -> schema -> table
	columns map[string]map[string]map[string]map[string]*models.Column        // catalog -> schema -> table -> column

	// Global metadata (user-defined abstractions)
	globalTables   map[string]*models.GlobalTable
	globalColumns  map[string]map[string]*models.GlobalColumn        // globalTable -> columnName -> column
	tableMappings  map[string][]*models.TableMapping                 // globalTable -> mappings
	columnMappings map[string]map[string][]*models.ColumnMapping     // globalTable -> columnName -> mappings
	columnRelationships map[string][]*models.ColumnRelationship      // globalTable -> relationships
}

func NewMemoryMetadataStorage() *MemoryMetadataStorage {
	return &MemoryMetadataStorage{
		catalogs: make(map[string]*models.Catalog),
		schemas:  make(map[string]map[string]*models.Schema),

		tables:  make(map[string]map[string]map[string]*models.Table),
		columns: make(map[string]map[string]map[string]map[string]*models.Column),

		globalTables:   make(map[string]*models.GlobalTable),
		globalColumns:  make(map[string]map[string]*models.GlobalColumn),
		tableMappings:  make(map[string][]*models.TableMapping),
		columnMappings: make(map[string]map[string][]*models.ColumnMapping),
		columnRelationships: make(map[string][]*models.ColumnRelationship),
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
// Local Table Operations (from data source discovery)
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
// Local Column Operations (from data source discovery)
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

	// Delete relationships where this table is source or target
	delete(m.columnRelationships, name)
	// Also remove relationships from other tables that reference this table
	for tableName, relationships := range m.columnRelationships {
		filtered := []*models.ColumnRelationship{}
		for _, rel := range relationships {
			if rel.SourceGlobalTableName != name && rel.TargetGlobalTableName != name {
				filtered = append(filtered, rel)
			}
		}
		m.columnRelationships[tableName] = filtered
	}

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

	// Delete relationships involving this column
	for tableName, relationships := range m.columnRelationships {
		filtered := []*models.ColumnRelationship{}
		for _, rel := range relationships {
			// Keep relationships that don't involve this specific column
			if !(rel.SourceGlobalTableName == globalTableName && rel.SourceGlobalColumnName == columnName) &&
			   !(rel.TargetGlobalTableName == globalTableName && rel.TargetGlobalColumnName == columnName) {
				filtered = append(filtered, rel)
			}
		}
		m.columnRelationships[tableName] = filtered
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

// ============================================================================
// Column Relationship Operations
// ============================================================================

func (m *MemoryMetadataStorage) CreateColumnRelationship(relationship *models.ColumnRelationship) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate all fields are non-empty
	if relationship.SourceGlobalTableName == "" || relationship.SourceGlobalColumnName == "" ||
		relationship.TargetGlobalTableName == "" || relationship.TargetGlobalColumnName == "" {
		return fmt.Errorf("all relationship fields (source table, source column, target table, target column) must be non-empty")
	}

	// Verify source table exists
	if _, exists := m.globalTables[relationship.SourceGlobalTableName]; !exists {
		return fmt.Errorf("source global table '%s' not found", relationship.SourceGlobalTableName)
	}

	// Verify source column exists
	sourceColumns, exists := m.globalColumns[relationship.SourceGlobalTableName]
	if !exists || sourceColumns[relationship.SourceGlobalColumnName] == nil {
		return fmt.Errorf("source global column '%s.%s' not found",
			relationship.SourceGlobalTableName, relationship.SourceGlobalColumnName)
	}

	// Verify target table exists
	if _, exists := m.globalTables[relationship.TargetGlobalTableName]; !exists {
		return fmt.Errorf("target global table '%s' not found", relationship.TargetGlobalTableName)
	}

	// Verify target column exists
	targetColumns, exists := m.globalColumns[relationship.TargetGlobalTableName]
	if !exists || targetColumns[relationship.TargetGlobalColumnName] == nil {
		return fmt.Errorf("target global column '%s.%s' not found",
			relationship.TargetGlobalTableName, relationship.TargetGlobalColumnName)
	}

	// Check for duplicate relationship
	existingRelationships := m.columnRelationships[relationship.SourceGlobalTableName]
	for _, existing := range existingRelationships {
		if existing.SourceGlobalTableName == relationship.SourceGlobalTableName &&
			existing.SourceGlobalColumnName == relationship.SourceGlobalColumnName &&
			existing.TargetGlobalTableName == relationship.TargetGlobalTableName &&
			existing.TargetGlobalColumnName == relationship.TargetGlobalColumnName {
			return fmt.Errorf("relationship already exists between %s.%s and %s.%s",
				relationship.SourceGlobalTableName, relationship.SourceGlobalColumnName,
				relationship.TargetGlobalTableName, relationship.TargetGlobalColumnName)
		}
	}

	// Store relationship bidirectionally for efficient queries
	m.columnRelationships[relationship.SourceGlobalTableName] = append(
		m.columnRelationships[relationship.SourceGlobalTableName], relationship)

	// If source and target tables are different, also store under target table
	if relationship.SourceGlobalTableName != relationship.TargetGlobalTableName {
		m.columnRelationships[relationship.TargetGlobalTableName] = append(
			m.columnRelationships[relationship.TargetGlobalTableName], relationship)
	}

	return nil
}

func (m *MemoryMetadataStorage) ListColumnRelationships(globalTableName string) ([]*models.ColumnRelationship, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	relationships, exists := m.columnRelationships[globalTableName]
	if !exists {
		return []*models.ColumnRelationship{}, nil
	}

	// Use a map to deduplicate in case of self-referential relationships
	seen := make(map[string]bool)
	result := []*models.ColumnRelationship{}

	for _, rel := range relationships {
		// Create a unique key for deduplication
		key := fmt.Sprintf("%s.%s->%s.%s",
			rel.SourceGlobalTableName, rel.SourceGlobalColumnName,
			rel.TargetGlobalTableName, rel.TargetGlobalColumnName)

		if !seen[key] {
			seen[key] = true
			result = append(result, rel)
		}
	}

	return result, nil
}

func (m *MemoryMetadataStorage) DeleteColumnRelationship(sourceTable, sourceColumn, targetTable, targetColumn string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove from source table's relationships
	sourceRelationships, exists := m.columnRelationships[sourceTable]
	if exists {
		filtered := []*models.ColumnRelationship{}
		for _, rel := range sourceRelationships {
			if !(rel.SourceGlobalTableName == sourceTable &&
				rel.SourceGlobalColumnName == sourceColumn &&
				rel.TargetGlobalTableName == targetTable &&
				rel.TargetGlobalColumnName == targetColumn) {
				filtered = append(filtered, rel)
			}
		}
		m.columnRelationships[sourceTable] = filtered
	}

	// Remove from target table's relationships (if different from source)
	if sourceTable != targetTable {
		targetRelationships, exists := m.columnRelationships[targetTable]
		if exists {
			filtered := []*models.ColumnRelationship{}
			for _, rel := range targetRelationships {
				if !(rel.SourceGlobalTableName == sourceTable &&
					rel.SourceGlobalColumnName == sourceColumn &&
					rel.TargetGlobalTableName == targetTable &&
					rel.TargetGlobalColumnName == targetColumn) {
					filtered = append(filtered, rel)
				}
			}
			m.columnRelationships[targetTable] = filtered
		}
	}

	return nil
}
