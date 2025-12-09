package storage

import "github.com/guilherme096/data-sync/pkg/data-sync/models"

type MetadataStorage interface {
	// Catalog operations
	CreateCatalog(catalog *models.Catalog) error
	UpdateCatalog(catalog *models.Catalog) error
	UpsertCatalog(catalog *models.Catalog) error
	GetCatalog(name string) (*models.Catalog, error)
	ListCatalogs() ([]*models.Catalog, error)

	// Schema operations
	CreateSchema(schema *models.Schema) error
	UpdateSchema(schema *models.Schema) error
	UpsertSchema(schema *models.Schema) error
	GetSchema(catalogName, schemaName string) (*models.Schema, error)
	ListSchemas(catalogName string) ([]*models.Schema, error)

	// Table operations
	CreateTable(table *models.Table) error
	UpdateTable(table *models.Table) error
	UpsertTable(table *models.Table) error
	GetTable(catalogName, schemaName, tableName string) (*models.Table, error)
	ListTables(catalogName, schemaName string) ([]*models.Table, error)

	// Column operations
	CreateColumn(column *models.Column) error
	UpdateColumn(column *models.Column) error
	UpsertColumn(column *models.Column) error
	GetColumn(catalogName, schemaName, tableName, columnName string) (*models.Column, error)
	ListColumns(catalogName, schemaName, tableName string) ([]*models.Column, error)
}
