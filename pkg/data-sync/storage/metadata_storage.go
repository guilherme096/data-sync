package storage

import "github.com/guilherme096/data-sync/pkg/data-sync/models"

type MetadataStorage interface {
	CreateCatalog(catalog *models.Catalog) error
	GetCatalog(name string) (*models.Catalog, error)
	ListCatalogs() ([]*models.Catalog, error)
	CreateSchema(schema *models.Schema) error
	GetSchema(catalogName, schemaName string) (*models.Schema, error)
	ListSchemas(catalogName string) ([]*models.Schema, error)
}
