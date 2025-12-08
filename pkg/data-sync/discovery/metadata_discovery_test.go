package discovery

import (
	"fmt"
	"testing"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
)

// Mock QueryEngine for testing
type mockQueryEngine struct {
	catalogsResult datasync.QueryResult
	schemasResult  datasync.QueryResult
	tablesResult   datasync.QueryResult
	columnsResult  datasync.QueryResult
	shouldError    bool
}

func (m *mockQueryEngine) ExecuteQuery(query string, params map[string]interface{}) (datasync.QueryResult, error) {
	if m.shouldError {
		return datasync.QueryResult{}, fmt.Errorf("mock error")
	}

	// Return catalogs for SHOW CATALOGS query
	if query == "SHOW CATALOGS" {
		return m.catalogsResult, nil
	}

	// Return schemas for SHOW SCHEMAS FROM queries
	if len(query) >= 17 && query[:17] == "SHOW SCHEMAS FROM" {
		return m.schemasResult, nil
	}

	// Return tables for SHOW TABLES FROM queries
	if len(query) >= 16 && query[:16] == "SHOW TABLES FROM" {
		return m.tablesResult, nil
	}

	// Return columns for DESCRIBE queries
	if len(query) >= 8 && query[:8] == "DESCRIBE" {
		return m.columnsResult, nil
	}

	return datasync.QueryResult{}, fmt.Errorf("unexpected query: %s", query)
}

func TestDiscoverCatalogs_Success(t *testing.T) {
	mockEngine := &mockQueryEngine{
		catalogsResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Catalog": "postgresql"},
				{"Catalog": "mysql"},
				{"Catalog": "mongodb"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	catalogs, err := discovery.DiscoverCatalogs()

	if err != nil {
		t.Fatalf("DiscoverCatalogs failed: %v", err)
	}

	if len(catalogs) != 3 {
		t.Errorf("Expected 3 catalogs, got %d", len(catalogs))
	}

	expectedNames := []string{"postgresql", "mysql", "mongodb"}
	for i, catalog := range catalogs {
		if catalog.Name != expectedNames[i] {
			t.Errorf("Expected catalog name '%s', got '%s'", expectedNames[i], catalog.Name)
		}
		if catalog.Metadata == nil {
			t.Error("Expected Metadata map to be initialized")
		}
	}
}

func TestDiscoverCatalogs_Empty(t *testing.T) {
	mockEngine := &mockQueryEngine{
		catalogsResult: datasync.QueryResult{
			Rows: []map[string]interface{}{},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	catalogs, err := discovery.DiscoverCatalogs()

	if err != nil {
		t.Fatalf("DiscoverCatalogs failed: %v", err)
	}

	if len(catalogs) != 0 {
		t.Errorf("Expected 0 catalogs, got %d", len(catalogs))
	}
}

func TestDiscoverCatalogs_Error(t *testing.T) {
	mockEngine := &mockQueryEngine{
		shouldError: true,
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	_, err := discovery.DiscoverCatalogs()

	if err == nil {
		t.Fatal("Expected error from DiscoverCatalogs, got nil")
	}
}

func TestDiscoverCatalogs_InvalidDataType(t *testing.T) {
	mockEngine := &mockQueryEngine{
		catalogsResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Catalog": "postgresql"},
				{"Catalog": 12345}, // Invalid: should be string
				{"Catalog": "mysql"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	catalogs, err := discovery.DiscoverCatalogs()

	if err != nil {
		t.Fatalf("DiscoverCatalogs failed: %v", err)
	}

	// Should skip the invalid entry
	if len(catalogs) != 2 {
		t.Errorf("Expected 2 catalogs (skipping invalid entry), got %d", len(catalogs))
	}
}

func TestDiscoverSchemas_Success(t *testing.T) {
	mockEngine := &mockQueryEngine{
		schemasResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Schema": "public"},
				{"Schema": "information_schema"},
				{"Schema": "pg_catalog"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	schemas, err := discovery.DiscoverSchemas("postgresql")

	if err != nil {
		t.Fatalf("DiscoverSchemas failed: %v", err)
	}

	if len(schemas) != 3 {
		t.Errorf("Expected 3 schemas, got %d", len(schemas))
	}

	expectedNames := []string{"public", "information_schema", "pg_catalog"}
	for i, schema := range schemas {
		if schema.Name != expectedNames[i] {
			t.Errorf("Expected schema name '%s', got '%s'", expectedNames[i], schema.Name)
		}
		if schema.CatalogName != "postgresql" {
			t.Errorf("Expected catalog name 'postgresql', got '%s'", schema.CatalogName)
		}
		if schema.Metadata == nil {
			t.Error("Expected Metadata map to be initialized")
		}
	}
}

func TestDiscoverSchemas_Empty(t *testing.T) {
	mockEngine := &mockQueryEngine{
		schemasResult: datasync.QueryResult{
			Rows: []map[string]interface{}{},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	schemas, err := discovery.DiscoverSchemas("postgresql")

	if err != nil {
		t.Fatalf("DiscoverSchemas failed: %v", err)
	}

	if len(schemas) != 0 {
		t.Errorf("Expected 0 schemas, got %d", len(schemas))
	}
}

func TestDiscoverSchemas_Error(t *testing.T) {
	mockEngine := &mockQueryEngine{
		shouldError: true,
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	_, err := discovery.DiscoverSchemas("postgresql")

	if err == nil {
		t.Fatal("Expected error from DiscoverSchemas, got nil")
	}
}

func TestDiscoverSchemas_InvalidDataType(t *testing.T) {
	mockEngine := &mockQueryEngine{
		schemasResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Schema": "public"},
				{"Schema": true}, // Invalid: should be string
				{"Schema": "information_schema"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	schemas, err := discovery.DiscoverSchemas("postgresql")

	if err != nil {
		t.Fatalf("DiscoverSchemas failed: %v", err)
	}

	// Should skip the invalid entry
	if len(schemas) != 2 {
		t.Errorf("Expected 2 schemas (skipping invalid entry), got %d", len(schemas))
	}
}

func TestDiscoverTables_Success(t *testing.T) {
	mockEngine := &mockQueryEngine{
		tablesResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Table": "users"},
				{"Table": "products"},
				{"Table": "orders"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	tables, err := discovery.DiscoverTables("postgresql", "public")

	if err != nil {
		t.Fatalf("DiscoverTables failed: %v", err)
	}

	if len(tables) != 3 {
		t.Errorf("Expected 3 tables, got %d", len(tables))
	}

	expectedNames := []string{"users", "products", "orders"}
	for i, table := range tables {
		if table.Name != expectedNames[i] {
			t.Errorf("Expected table name '%s', got '%s'", expectedNames[i], table.Name)
		}
		if table.SchemaName != "public" {
			t.Errorf("Expected schema name 'public', got '%s'", table.SchemaName)
		}
		if table.CatalogName != "postgresql" {
			t.Errorf("Expected catalog name 'postgresql', got '%s'", table.CatalogName)
		}
		if table.Metadata == nil {
			t.Error("Expected Metadata map to be initialized")
		}
	}
}

func TestDiscoverTables_Empty(t *testing.T) {
	mockEngine := &mockQueryEngine{
		tablesResult: datasync.QueryResult{
			Rows: []map[string]interface{}{},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	tables, err := discovery.DiscoverTables("postgresql", "public")

	if err != nil {
		t.Fatalf("DiscoverTables failed: %v", err)
	}

	if len(tables) != 0 {
		t.Errorf("Expected 0 tables, got %d", len(tables))
	}
}

func TestDiscoverTables_Error(t *testing.T) {
	mockEngine := &mockQueryEngine{
		shouldError: true,
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	_, err := discovery.DiscoverTables("postgresql", "public")

	if err == nil {
		t.Fatal("Expected error from DiscoverTables, got nil")
	}
}

func TestDiscoverTables_InvalidDataType(t *testing.T) {
	mockEngine := &mockQueryEngine{
		tablesResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Table": "users"},
				{"Table": 12345}, // Invalid: should be string
				{"Table": "products"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	tables, err := discovery.DiscoverTables("postgresql", "public")

	if err != nil {
		t.Fatalf("DiscoverTables failed: %v", err)
	}

	// Should skip the invalid entry
	if len(tables) != 2 {
		t.Errorf("Expected 2 tables (skipping invalid entry), got %d", len(tables))
	}
}

func TestDiscoverColumns_Success(t *testing.T) {
	mockEngine := &mockQueryEngine{
		columnsResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Column": "id", "Type": "bigint"},
				{"Column": "name", "Type": "varchar(255)"},
				{"Column": "created_at", "Type": "timestamp"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	columns, err := discovery.DiscoverColumns("postgresql", "public", "users")

	if err != nil {
		t.Fatalf("DiscoverColumns failed: %v", err)
	}

	if len(columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(columns))
	}

	expectedNames := []string{"id", "name", "created_at"}
	expectedTypes := []string{"bigint", "varchar(255)", "timestamp"}
	for i, column := range columns {
		if column.Name != expectedNames[i] {
			t.Errorf("Expected column name '%s', got '%s'", expectedNames[i], column.Name)
		}
		if column.DataType != expectedTypes[i] {
			t.Errorf("Expected data type '%s', got '%s'", expectedTypes[i], column.DataType)
		}
		if column.TableName != "users" {
			t.Errorf("Expected table name 'users', got '%s'", column.TableName)
		}
		if column.SchemaName != "public" {
			t.Errorf("Expected schema name 'public', got '%s'", column.SchemaName)
		}
		if column.CatalogName != "postgresql" {
			t.Errorf("Expected catalog name 'postgresql', got '%s'", column.CatalogName)
		}
		if column.Metadata == nil {
			t.Error("Expected Metadata map to be initialized")
		}
	}
}

func TestDiscoverColumns_Empty(t *testing.T) {
	mockEngine := &mockQueryEngine{
		columnsResult: datasync.QueryResult{
			Rows: []map[string]interface{}{},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	columns, err := discovery.DiscoverColumns("postgresql", "public", "users")

	if err != nil {
		t.Fatalf("DiscoverColumns failed: %v", err)
	}

	if len(columns) != 0 {
		t.Errorf("Expected 0 columns, got %d", len(columns))
	}
}

func TestDiscoverColumns_Error(t *testing.T) {
	mockEngine := &mockQueryEngine{
		shouldError: true,
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	_, err := discovery.DiscoverColumns("postgresql", "public", "users")

	if err == nil {
		t.Fatal("Expected error from DiscoverColumns, got nil")
	}
}

func TestDiscoverColumns_InvalidDataType(t *testing.T) {
	mockEngine := &mockQueryEngine{
		columnsResult: datasync.QueryResult{
			Rows: []map[string]interface{}{
				{"Column": "id", "Type": "bigint"},
				{"Column": 12345, "Type": "varchar"}, // Invalid: Column should be string
				{"Column": "name", "Type": 999},      // Invalid: Type should be string
				{"Column": "email", "Type": "varchar(255)"},
			},
		},
	}

	discovery := NewTrinoMetadataDiscovery(mockEngine)
	columns, err := discovery.DiscoverColumns("postgresql", "public", "users")

	if err != nil {
		t.Fatalf("DiscoverColumns failed: %v", err)
	}

	// Should skip the invalid entries
	if len(columns) != 2 {
		t.Errorf("Expected 2 columns (skipping invalid entries), got %d", len(columns))
	}
}
