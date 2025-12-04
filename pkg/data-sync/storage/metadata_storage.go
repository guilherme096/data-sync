package storage

import "github.com/guilherme096/data-sync/pkg/data-sync/models"

type MetadataStorage interface {
	CreateCatalog(catalog *models.Catalog) error
	GetCatalog(id int) (*models.Catalog, error)
	ListCatalogs() ([]*models.Catalog, error)
	CreateSchema(schema *models.Schema) error
	GetSchema(id int) (*models.Schema, error)
	ListSchemas(catalogID int) ([]*models.Schema, error)
}
