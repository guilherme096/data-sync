package storage

import (
	"testing"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
)

func TestCreateAndGetCatalog(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{
		Name:     "postgresql",
		Metadata: map[string]string{"type": "relational"},
	}

	err := storage.CreateCatalog(catalog)
	if err != nil {
		t.Fatalf("CreateCatalog failed: %v", err)
	}

	retrieved, err := storage.GetCatalog("postgresql")
	if err != nil {
		t.Fatalf("GetCatalog failed: %v", err)
	}

	if retrieved.Name != "postgresql" {
		t.Errorf("Expected catalog name 'postgresql', got '%s'", retrieved.Name)
	}
}

func TestCreateCatalog_EmptyName(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{
		Name:     "",
		Metadata: map[string]string{},
	}

	err := storage.CreateCatalog(catalog)
	if err == nil {
		t.Fatal("Expected error for empty catalog name, got nil")
	}
}

func TestCreateCatalog_Duplicate(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{
		Name:     "mysql",
		Metadata: map[string]string{},
	}

	err := storage.CreateCatalog(catalog)
	if err != nil {
		t.Fatalf("First CreateCatalog failed: %v", err)
	}

	err = storage.CreateCatalog(catalog)
	if err == nil {
		t.Fatal("Expected error for duplicate catalog, got nil")
	}
}

func TestGetCatalog_NotFound(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	_, err := storage.GetCatalog("nonexistent")
	if err == nil {
		t.Fatal("Expected error for nonexistent catalog, got nil")
	}
}

func TestListCatalogs(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalogs := []*models.Catalog{
		{Name: "postgresql", Metadata: map[string]string{}},
		{Name: "mysql", Metadata: map[string]string{}},
		{Name: "mongodb", Metadata: map[string]string{}},
	}

	for _, catalog := range catalogs {
		err := storage.CreateCatalog(catalog)
		if err != nil {
			t.Fatalf("CreateCatalog failed: %v", err)
		}
	}

	retrieved, err := storage.ListCatalogs()
	if err != nil {
		t.Fatalf("ListCatalogs failed: %v", err)
	}

	if len(retrieved) != 3 {
		t.Errorf("Expected 3 catalogs, got %d", len(retrieved))
	}
}

func TestCreateAndGetSchema(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	// Create catalog first
	catalog := &models.Catalog{
		Name:     "postgresql",
		Metadata: map[string]string{},
	}
	err := storage.CreateCatalog(catalog)
	if err != nil {
		t.Fatalf("CreateCatalog failed: %v", err)
	}

	// Create schema
	schema := &models.Schema{
		Name:        "public",
		CatalogName: "postgresql",
		Metadata:    map[string]string{"owner": "postgres"},
	}

	err = storage.CreateSchema(schema)
	if err != nil {
		t.Fatalf("CreateSchema failed: %v", err)
	}

	retrieved, err := storage.GetSchema("postgresql", "public")
	if err != nil {
		t.Fatalf("GetSchema failed: %v", err)
	}

	if retrieved.Name != "public" || retrieved.CatalogName != "postgresql" {
		t.Errorf("Expected schema 'public' in catalog 'postgresql', got '%s' in '%s'",
			retrieved.Name, retrieved.CatalogName)
	}
}

func TestCreateSchema_EmptyName(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{Name: "postgresql", Metadata: map[string]string{}}
	storage.CreateCatalog(catalog)

	schema := &models.Schema{
		Name:        "",
		CatalogName: "postgresql",
		Metadata:    map[string]string{},
	}

	err := storage.CreateSchema(schema)
	if err == nil {
		t.Fatal("Expected error for empty schema name, got nil")
	}
}

func TestCreateSchema_CatalogNotFound(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	schema := &models.Schema{
		Name:        "public",
		CatalogName: "nonexistent",
		Metadata:    map[string]string{},
	}

	err := storage.CreateSchema(schema)
	if err == nil {
		t.Fatal("Expected error for nonexistent catalog, got nil")
	}
}

func TestCreateSchema_Duplicate(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{Name: "postgresql", Metadata: map[string]string{}}
	storage.CreateCatalog(catalog)

	schema := &models.Schema{
		Name:        "public",
		CatalogName: "postgresql",
		Metadata:    map[string]string{},
	}

	err := storage.CreateSchema(schema)
	if err != nil {
		t.Fatalf("First CreateSchema failed: %v", err)
	}

	err = storage.CreateSchema(schema)
	if err == nil {
		t.Fatal("Expected error for duplicate schema, got nil")
	}
}

func TestListSchemas(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{Name: "postgresql", Metadata: map[string]string{}}
	storage.CreateCatalog(catalog)

	schemas := []*models.Schema{
		{Name: "public", CatalogName: "postgresql", Metadata: map[string]string{}},
		{Name: "information_schema", CatalogName: "postgresql", Metadata: map[string]string{}},
		{Name: "pg_catalog", CatalogName: "postgresql", Metadata: map[string]string{}},
	}

	for _, schema := range schemas {
		err := storage.CreateSchema(schema)
		if err != nil {
			t.Fatalf("CreateSchema failed: %v", err)
		}
	}

	retrieved, err := storage.ListSchemas("postgresql")
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}

	if len(retrieved) != 3 {
		t.Errorf("Expected 3 schemas, got %d", len(retrieved))
	}
}

func TestListSchemas_EmptyCatalog(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	catalog := &models.Catalog{Name: "postgresql", Metadata: map[string]string{}}
	storage.CreateCatalog(catalog)

	schemas, err := storage.ListSchemas("postgresql")
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}

	if len(schemas) != 0 {
		t.Errorf("Expected 0 schemas for empty catalog, got %d", len(schemas))
	}
}

func TestListSchemas_NonexistentCatalog(t *testing.T) {
	storage := NewMemoryMetadataStorage()

	schemas, err := storage.ListSchemas("nonexistent")
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}

	if len(schemas) != 0 {
		t.Errorf("Expected 0 schemas for nonexistent catalog, got %d", len(schemas))
	}
}
