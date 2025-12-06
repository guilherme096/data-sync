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
	return m.schemasResult, nil
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
