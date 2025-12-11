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

	// Local table operations (for discovered tables from data sources)
	CreateTable(table *models.Table) error
	UpdateTable(table *models.Table) error
	UpsertTable(table *models.Table) error
	GetTable(catalogName, schemaName, tableName string) (*models.Table, error)
	ListTables(catalogName, schemaName string) ([]*models.Table, error)

	// Local column operations (for discovered columns from data sources)
	CreateColumn(column *models.Column) error
	UpdateColumn(column *models.Column) error
	UpsertColumn(column *models.Column) error
	GetColumn(catalogName, schemaName, tableName, columnName string) (*models.Column, error)
	ListColumns(catalogName, schemaName, tableName string) ([]*models.Column, error)

	// Global table operations
	CreateGlobalTable(table *models.GlobalTable) error
	GetGlobalTable(name string) (*models.GlobalTable, error)
	ListGlobalTables() ([]*models.GlobalTable, error)
	DeleteGlobalTable(name string) error

	// Global column operations
	CreateGlobalColumn(column *models.GlobalColumn) error
	ListGlobalColumns(globalTableName string) ([]*models.GlobalColumn, error)
	DeleteGlobalColumn(globalTableName, columnName string) error

	// Table mapping operations
	CreateTableMapping(mapping *models.TableMapping) error
	ListTableMappings(globalTableName string) ([]*models.TableMapping, error)
	DeleteTableMapping(globalTableName, catalog, schema, table string) error

	// Column mapping operations
	CreateColumnMapping(mapping *models.ColumnMapping) error
	ListColumnMappings(globalTableName, globalColumnName string) ([]*models.ColumnMapping, error)
	DeleteColumnMapping(globalTableName, globalColumnName, catalog, schema, table, column string) error
}
