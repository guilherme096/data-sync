package query

import (
	"fmt"
	"strings"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
)

// SQLGenerator generates Trino SQL from resolved components
type SQLGenerator struct{}

// NewSQLGenerator creates a new SQL generator
func NewSQLGenerator() *SQLGenerator {
	return &SQLGenerator{}
}

// GenerateSQL builds a simple SELECT query for a physical table
// For Phase 1: Only supports simple SELECT with optional WHERE and LIMIT clauses
func (g *SQLGenerator) GenerateSQL(
	physicalTable *models.TableMapping,
	columnMap map[string]string,
	globalColumns []string,
	whereClause string,
	limitClause string,
	isSelectAll bool,
) (string, error) {
	// Build fully qualified table name
	fullTableName := fmt.Sprintf("%s.%s.%s",
		physicalTable.CatalogName,
		physicalTable.SchemaName,
		physicalTable.TableName,
	)

	// Build SELECT clause
	var selectClause string
	if isSelectAll {
		selectClause = "*"
	} else {
		// Map global columns to physical columns
		physicalCols := make([]string, len(globalColumns))
		for i, globalCol := range globalColumns {
			physicalCol, exists := columnMap[globalCol]
			if !exists {
				return "", fmt.Errorf("column '%s' not found in column mapping", globalCol)
			}
			physicalCols[i] = physicalCol
		}
		selectClause = strings.Join(physicalCols, ", ")
	}

	// Build the query
	query := fmt.Sprintf("SELECT %s FROM %s", selectClause, fullTableName)

	// Add WHERE clause if present
	if whereClause != "" {
		// For Phase 1, pass through WHERE clause as-is
		// TODO: In later phases, map column names in WHERE clause
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	// Add LIMIT clause if present
	if limitClause != "" {
		query += fmt.Sprintf(" LIMIT %s", limitClause)
	}

	return query, nil
}

// GenerateUnionSQL builds a UNION query from multiple physical tables or a relation
// For Phase 2: Supports UNION of physical tables
func (g *SQLGenerator) GenerateUnionSQL(
	tables []*models.TableMapping,
	columnMaps []map[string]string,
	globalColumns []string,
	whereClause string,
	limitClause string,
	isSelectAll bool,
) (string, error) {
	if len(tables) == 0 {
		return "", fmt.Errorf("no tables provided for UNION")
	}

	if len(tables) != len(columnMaps) {
		return "", fmt.Errorf("mismatch between tables and column maps")
	}

	queries := make([]string, len(tables))

	for i, table := range tables {
		// Generate SELECT for each table (without LIMIT for individual queries)
		query, err := g.GenerateSQL(table, columnMaps[i], globalColumns, whereClause, "", isSelectAll)
		if err != nil {
			return "", fmt.Errorf("failed to generate SQL for table %s.%s.%s: %w",
				table.CatalogName, table.SchemaName, table.TableName, err)
		}
		queries[i] = query
	}

	// Join with UNION
	unionQuery := strings.Join(queries, " UNION ")

	// Add LIMIT clause after UNION
	if limitClause != "" {
		unionQuery += fmt.Sprintf(" LIMIT %s", limitClause)
	}

	return unionQuery, nil
}

// GenerateUnionFromRelation builds a UNION query from a resolved relation
// For Phase 2: Supports simple UNION relations (no nesting)
func (g *SQLGenerator) GenerateUnionFromRelation(
	relation *ResolvedRelation,
	columnMaps []map[string]string,
	globalColumns []string,
	whereClause string,
	limitClause string,
	isSelectAll bool,
) (string, error) {
	if relation.RelationType != "UNION" {
		return "", fmt.Errorf("relation type '%s' not supported for UNION generation", relation.RelationType)
	}

	// Get physical tables from relation
	tables := []*RelationNode{relation.LeftNode, relation.RightNode}

	if len(columnMaps) != len(tables) {
		return "", fmt.Errorf("mismatch between relation tables and column maps")
	}

	queries := make([]string, len(tables))

	for i, node := range tables {
		if node.Type != NodeTypePhysical {
			return "", fmt.Errorf("nested relations not supported in Phase 2")
		}

		// Build table mapping from node
		tableMapping := &models.TableMapping{
			CatalogName: node.Catalog,
			SchemaName:  node.Schema,
			TableName:   node.Table,
		}

		// Generate SELECT for this table (without LIMIT for individual queries)
		query, err := g.GenerateSQL(tableMapping, columnMaps[i], globalColumns, whereClause, "", isSelectAll)
		if err != nil {
			return "", fmt.Errorf("failed to generate SQL for table %s.%s.%s: %w",
				node.Catalog, node.Schema, node.Table, err)
		}
		queries[i] = query
	}

	// Join with UNION
	unionQuery := strings.Join(queries, " UNION ")

	// Add LIMIT clause after UNION
	if limitClause != "" {
		unionQuery += fmt.Sprintf(" LIMIT %s", limitClause)
	}

	return unionQuery, nil
}

// GenerateJoinFromRelation builds a JOIN query from a resolved relation
// For Phase 3: Supports simple JOIN relations (no nesting)
func (g *SQLGenerator) GenerateJoinFromRelation(
	relation *ResolvedRelation,
	columnMaps []map[string]string,
	globalColumns []string,
	whereClause string,
	limitClause string,
	isSelectAll bool,
) (string, error) {
	if relation.RelationType != "JOIN" {
		return "", fmt.Errorf("relation type '%s' not supported for JOIN generation", relation.RelationType)
	}

	if relation.JoinColumn == nil {
		return "", fmt.Errorf("JOIN relation requires join columns")
	}

	// Verify both nodes are physical
	if relation.LeftNode.Type != NodeTypePhysical || relation.RightNode.Type != NodeTypePhysical {
		return "", fmt.Errorf("nested relations not supported in Phase 3")
	}

	// Build table names with aliases
	leftTableFull := fmt.Sprintf("%s.%s.%s",
		relation.LeftNode.Catalog,
		relation.LeftNode.Schema,
		relation.LeftNode.Table)
	rightTableFull := fmt.Sprintf("%s.%s.%s",
		relation.RightNode.Catalog,
		relation.RightNode.Schema,
		relation.RightNode.Table)

	leftAlias := "t1"
	rightAlias := "t2"

	// Build SELECT clause with qualified column names
	var selectClause string
	if isSelectAll {
		// For SELECT *, we need to qualify columns to avoid ambiguity
		selectClause = fmt.Sprintf("%s.*, %s.*", leftAlias, rightAlias)
	} else {
		// Map each global column to its physical column with table alias
		if len(columnMaps) != 2 {
			return "", fmt.Errorf("expected 2 column maps for JOIN, got %d", len(columnMaps))
		}

		selectParts := make([]string, 0)
		for _, globalCol := range globalColumns {
			// Check left table first
			if physicalCol, exists := columnMaps[0][globalCol]; exists {
				selectParts = append(selectParts, fmt.Sprintf("%s.%s", leftAlias, physicalCol))
			} else if physicalCol, exists := columnMaps[1][globalCol]; exists {
				// Check right table
				selectParts = append(selectParts, fmt.Sprintf("%s.%s", rightAlias, physicalCol))
			} else {
				return "", fmt.Errorf("column '%s' not found in either table", globalCol)
			}
		}
		selectClause = strings.Join(selectParts, ", ")
	}

	// Build the JOIN query
	query := fmt.Sprintf("SELECT %s FROM %s %s JOIN %s %s ON %s.%s = %s.%s",
		selectClause,
		leftTableFull,
		leftAlias,
		rightTableFull,
		rightAlias,
		leftAlias,
		relation.JoinColumn.Left,
		rightAlias,
		relation.JoinColumn.Right,
	)

	// Add WHERE clause if present
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	// Add LIMIT clause if present
	if limitClause != "" {
		query += fmt.Sprintf(" LIMIT %s", limitClause)
	}

	return query, nil
}
