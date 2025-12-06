package discovery

import (
	"fmt"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
)

// MetadataDiscovery discovers catalogs and schemas from Trino
type MetadataDiscovery interface {
	DiscoverCatalogs() ([]*models.Catalog, error)
	DiscoverSchemas(catalogName string) ([]*models.Schema, error)
}

type trinoMetadataDiscovery struct {
	engine datasync.QueryEngine
}

func NewTrinoMetadataDiscovery(engine datasync.QueryEngine) MetadataDiscovery {
	return &trinoMetadataDiscovery{
		engine: engine,
	}
}

func (d *trinoMetadataDiscovery) DiscoverCatalogs() ([]*models.Catalog, error) {
	result, err := d.engine.ExecuteQuery("SHOW CATALOGS", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover catalogs: %w", err)
	}

	catalogs := make([]*models.Catalog, 0, len(result.Rows))
	for _, row := range result.Rows {
		// Trino returns catalog name in a column (typically "Catalog")
		if catalogName, ok := row["Catalog"].(string); ok {
			catalogs = append(catalogs, &models.Catalog{
				Name:     catalogName,
				Metadata: make(map[string]string),
			})
		}
	}

	return catalogs, nil
}

func (d *trinoMetadataDiscovery) DiscoverSchemas(catalogName string) ([]*models.Schema, error) {
	query := fmt.Sprintf("SHOW SCHEMAS FROM %s", catalogName)
	result, err := d.engine.ExecuteQuery(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover schemas for catalog %s: %w", catalogName, err)
	}

	schemas := make([]*models.Schema, 0, len(result.Rows))
	for _, row := range result.Rows {
		// Trino returns schema name in a column (typically "Schema")
		if schemaName, ok := row["Schema"].(string); ok {
			schemas = append(schemas, &models.Schema{
				Name:        schemaName,
				CatalogName: catalogName,
				Metadata:    make(map[string]string),
			})
		}
	}

	return schemas, nil
}
