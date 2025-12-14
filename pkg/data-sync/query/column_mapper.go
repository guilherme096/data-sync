package query

import (
	"fmt"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

// ColumnMapper maps global columns to physical columns
type ColumnMapper struct {
	storage storage.MetadataStorage
}

// NewColumnMapper creates a new column mapper
func NewColumnMapper(storage storage.MetadataStorage) *ColumnMapper {
	return &ColumnMapper{
		storage: storage,
	}
}

// MapColumns maps a list of global column names to their physical counterparts
// Returns a map of global column name -> physical column name
func (m *ColumnMapper) MapColumns(globalTableName string, globalColumns []string, physicalTable *models.TableMapping) (map[string]string, error) {
	columnMap := make(map[string]string)

	for _, globalCol := range globalColumns {
		physicalCol, err := m.mapSingleColumn(globalTableName, globalCol, physicalTable)
		if err != nil {
			return nil, err
		}
		columnMap[globalCol] = physicalCol
	}

	return columnMap, nil
}

// mapSingleColumn maps a single global column to its physical column
func (m *ColumnMapper) mapSingleColumn(globalTableName, globalColumnName string, physicalTable *models.TableMapping) (string, error) {
	// Get column mappings for this global column
	mappings, err := m.storage.ListColumnMappings(globalTableName, globalColumnName)
	if err != nil {
		return "", fmt.Errorf("failed to get mappings for column '%s' in global table '%s': %w", globalColumnName, globalTableName, err)
	}

	if len(mappings) == 0 {
		return "", fmt.Errorf("no column mapping found for '%s.%s'", globalTableName, globalColumnName)
	}

	// Find the mapping that matches our physical table
	for _, mapping := range mappings {
		if mapping.CatalogName == physicalTable.CatalogName &&
			mapping.SchemaName == physicalTable.SchemaName &&
			mapping.TableName == physicalTable.TableName {
			return mapping.ColumnName, nil
		}
	}

	return "", fmt.Errorf("no column mapping found for '%s.%s' in physical table '%s.%s.%s'",
		globalTableName, globalColumnName,
		physicalTable.CatalogName, physicalTable.SchemaName, physicalTable.TableName)
}

// GetAllColumns retrieves all global columns for a global table
func (m *ColumnMapper) GetAllColumns(globalTableName string) ([]string, error) {
	globalColumns, err := m.storage.ListGlobalColumns(globalTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for global table '%s': %w", globalTableName, err)
	}

	columnNames := make([]string, len(globalColumns))
	for i, col := range globalColumns {
		columnNames[i] = col.Name
	}

	return columnNames, nil
}

// MapColumnsForMultipleTables maps global columns to physical columns for multiple tables
// Used for UNION relations - returns a slice of column maps, one per table
func (m *ColumnMapper) MapColumnsForMultipleTables(
	globalTableName string,
	globalColumns []string,
	physicalTables []*models.TableMapping,
) ([]map[string]string, error) {
	columnMaps := make([]map[string]string, len(physicalTables))

	for i, table := range physicalTables {
		columnMap, err := m.MapColumns(globalTableName, globalColumns, table)
		if err != nil {
			return nil, fmt.Errorf("failed to map columns for table %s.%s.%s: %w",
				table.CatalogName, table.SchemaName, table.TableName, err)
		}
		columnMaps[i] = columnMap
	}

	return columnMaps, nil
}

// MapColumnsForRelation maps global columns for a resolved relation
// Returns column maps for both left and right tables
func (m *ColumnMapper) MapColumnsForRelation(
	globalTableName string,
	globalColumns []string,
	relation *ResolvedRelation,
) ([]map[string]string, error) {
	// For Phase 2, only support physical tables in relations
	if relation.LeftNode.Type != NodeTypePhysical || relation.RightNode.Type != NodeTypePhysical {
		return nil, fmt.Errorf("nested relations not supported in Phase 2")
	}

	// Create table mappings from relation nodes
	leftTable := &models.TableMapping{
		GlobalTableName: globalTableName,
		CatalogName:     relation.LeftNode.Catalog,
		SchemaName:      relation.LeftNode.Schema,
		TableName:       relation.LeftNode.Table,
	}

	rightTable := &models.TableMapping{
		GlobalTableName: globalTableName,
		CatalogName:     relation.RightNode.Catalog,
		SchemaName:      relation.RightNode.Schema,
		TableName:       relation.RightNode.Table,
	}

	tables := []*models.TableMapping{leftTable, rightTable}
	return m.MapColumnsForMultipleTables(globalTableName, globalColumns, tables)
}
