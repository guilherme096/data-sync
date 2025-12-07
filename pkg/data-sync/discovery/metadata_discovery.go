package discovery

import (
	"fmt"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
)

// MetadataDiscovery discovers catalogs, schemas, tables, and columns from Trino
type MetadataDiscovery interface {
	DiscoverCatalogs() ([]*models.Catalog, error)
	DiscoverSchemas(catalogName string) ([]*models.Schema, error)
	DiscoverTables(catalogName, schemaName string) ([]*models.Table, error)
	DiscoverColumns(catalogName, schemaName, tableName string) ([]*models.Column, error)
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

func (d *trinoMetadataDiscovery) DiscoverTables(catalogName, schemaName string) ([]*models.Table, error) {
	query := fmt.Sprintf("SHOW TABLES FROM %s.%s", catalogName, schemaName)
	result, err := d.engine.ExecuteQuery(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover tables for catalog %s and schema %s: %w", catalogName, schemaName, err)
	}

	tables := make([]*models.Table, 0, len(result.Rows))
	for _, row := range result.Rows {
		if tableName, ok := row["Table"].(string); ok {
			tables = append(tables, &models.Table{
				Name:        tableName,
				SchemaName:  schemaName,
				CatalogName: catalogName,
				Metadata:    make(map[string]string),
			})
		}
	}

	return tables, nil
}

func (d *trinoMetadataDiscovery) DiscoverColumns(catalogName, schemaName, tableName string) ([]*models.Column, error) {
	query := fmt.Sprintf("DESCRIBE %s.%s.%s", catalogName, schemaName, tableName)
	result, err := d.engine.ExecuteQuery(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover columns for table %s.%s.%s: %w", catalogName, schemaName, tableName, err)
	}

	columns := make([]*models.Column, 0, len(result.Rows))
	for _, row := range result.Rows {
		columnName, nameOk := row["Column"].(string)
		dataType, typeOk := row["Type"].(string)

		if nameOk && typeOk {
			columns = append(columns, &models.Column{
				Name:        columnName,
				TableName:   tableName,
				SchemaName:  schemaName,
				CatalogName: catalogName,
				DataType:    dataType,
				Metadata:    make(map[string]string),
			})
		}
	}

	return columns, nil
}
