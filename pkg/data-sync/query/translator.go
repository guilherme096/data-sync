package query

import (
	"fmt"
	"time"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

// QueryTranslator translates queries on global tables to Trino SQL
type QueryTranslator interface {
	Translate(globalQuery string) (trinoQuery string, error error)
	TranslateAndExecute(globalQuery string) (*QueryResult, error)
}

// QueryResult contains the results of a translated and executed query
type QueryResult struct {
	GeneratedSQL  string                   `json:"generatedSQL"`
	Rows          []map[string]interface{} `json:"rows"`
	RowCount      int                      `json:"rowCount"`
	ExecutionTime string                   `json:"executionTime"`
}

// Translator implements QueryTranslator
type Translator struct {
	parser       *QueryParser
	resolver     *TableResolver
	columnMapper *ColumnMapper
	generator    *SQLGenerator
	engine       datasync.QueryEngine
	storage      storage.MetadataStorage
}

// NewTranslator creates a new query translator
func NewTranslator(storage storage.MetadataStorage, engine datasync.QueryEngine) *Translator {
	return &Translator{
		parser:       NewQueryParser(),
		resolver:     NewTableResolver(storage),
		columnMapper: NewColumnMapper(storage),
		generator:    NewSQLGenerator(),
		engine:       engine,
		storage:      storage,
	}
}

// Translate converts a query on global tables to executable Trino SQL
func (t *Translator) Translate(globalQuery string) (string, error) {
	// 1. Parse the query
	parsed, err := t.parser.Parse(globalQuery)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	// 2. Resolve global table to physical table
	physicalTable, err := t.resolver.ResolveGlobalTable(parsed.TableName)
	if err != nil {
		return "", fmt.Errorf("resolution error: %w", err)
	}

	// 3. Get columns to map
	columnsToMap := parsed.Columns
	if parsed.IsSelectAll {
		// Get all global columns for SELECT *
		allCols, err := t.columnMapper.GetAllColumns(parsed.TableName)
		if err != nil {
			return "", fmt.Errorf("column resolution error: %w", err)
		}
		columnsToMap = allCols
	}

	// 4. Map columns
	columnMap, err := t.columnMapper.MapColumns(parsed.TableName, columnsToMap, physicalTable)
	if err != nil {
		return "", fmt.Errorf("column mapping error: %w", err)
	}

	// 5. Generate SQL
	trinoSQL, err := t.generator.GenerateSQL(
		physicalTable,
		columnMap,
		columnsToMap,
		parsed.WhereClause,
		parsed.LimitClause,
		parsed.IsSelectAll,
	)
	if err != nil {
		return "", fmt.Errorf("SQL generation error: %w", err)
	}

	return trinoSQL, nil
}

// TranslateAndExecute translates the query and executes it against Trino
func (t *Translator) TranslateAndExecute(globalQuery string) (*QueryResult, error) {
	startTime := time.Now()

	// Try Phase 2 translation first (supports UNION, etc.)
	trinoSQL, err := t.TranslateAdvanced(globalQuery)
	if err != nil {
		// Fall back to Phase 1 translation
		trinoSQL, err = t.Translate(globalQuery)
		if err != nil {
			return nil, err
		}
	}

	// Execute the query
	result, err := t.engine.ExecuteQuery(trinoSQL, nil)
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}

	executionTime := time.Since(startTime)

	return &QueryResult{
		GeneratedSQL:  trinoSQL,
		Rows:          result.Rows,
		RowCount:      len(result.Rows),
		ExecutionTime: fmt.Sprintf("%dms", executionTime.Milliseconds()),
	}, nil
}

// TranslateAdvanced converts a query on global tables to executable Trino SQL
// Phase 2: Supports UNION relations and multiple table mappings
func (t *Translator) TranslateAdvanced(globalQuery string) (string, error) {
	// 1. Parse the query
	parsed, err := t.parser.Parse(globalQuery)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	// 2. Resolve global table using advanced resolver
	resolved, err := t.resolver.ResolveGlobalTableAdvanced(parsed.TableName)
	if err != nil {
		return "", fmt.Errorf("resolution error: %w", err)
	}

	// 3. Get columns to map
	columnsToMap := parsed.Columns
	if parsed.IsSelectAll {
		// Get all global columns for SELECT *
		allCols, err := t.columnMapper.GetAllColumns(parsed.TableName)
		if err != nil {
			return "", fmt.Errorf("column resolution error: %w", err)
		}
		columnsToMap = allCols
	}

	// 4. Handle different resolution types
	if resolved.IsRelation {
		// Explicit relation (UNION or JOIN)
		return t.translateRelation(parsed, resolved.Relation, columnsToMap)
	} else if resolved.MultipleMappings != nil {
		// Multiple mappings - auto-generate UNION
		return t.translateMultipleMappings(parsed, resolved.MultipleMappings, columnsToMap)
	} else {
		// Single mapping - use simple translation
		columnMap, err := t.columnMapper.MapColumns(parsed.TableName, columnsToMap, resolved.SingleMapping)
		if err != nil {
			return "", fmt.Errorf("column mapping error: %w", err)
		}

		return t.generator.GenerateSQL(
			resolved.SingleMapping,
			columnMap,
			columnsToMap,
			parsed.WhereClause,
			parsed.LimitClause,
			parsed.IsSelectAll,
		)
	}
}

// translateRelation handles translation for explicit relations
func (t *Translator) translateRelation(
	parsed *ParsedQuery,
	relation *ResolvedRelation,
	columnsToMap []string,
) (string, error) {
	// Map columns for the relation
	columnMaps, err := t.columnMapper.MapColumnsForRelation(parsed.TableName, columnsToMap, relation)
	if err != nil {
		return "", fmt.Errorf("column mapping error: %w", err)
	}

	// Generate SQL based on relation type
	switch relation.RelationType {
	case "UNION":
		return t.generator.GenerateUnionFromRelation(
			relation,
			columnMaps,
			columnsToMap,
			parsed.WhereClause,
			parsed.LimitClause,
			parsed.IsSelectAll,
		)
	case "JOIN":
		return t.generator.GenerateJoinFromRelation(
			relation,
			columnMaps,
			columnsToMap,
			parsed.WhereClause,
			parsed.LimitClause,
			parsed.IsSelectAll,
		)
	default:
		return "", fmt.Errorf("unsupported relation type: %s", relation.RelationType)
	}
}

// translateMultipleMappings handles translation for multiple table mappings (auto-UNION)
func (t *Translator) translateMultipleMappings(
	parsed *ParsedQuery,
	mappings []*models.TableMapping,
	columnsToMap []string,
) (string, error) {
	// Map columns for each table
	columnMaps, err := t.columnMapper.MapColumnsForMultipleTables(parsed.TableName, columnsToMap, mappings)
	if err != nil {
		return "", fmt.Errorf("column mapping error: %w", err)
	}

	// Generate UNION SQL
	return t.generator.GenerateUnionSQL(
		mappings,
		columnMaps,
		columnsToMap,
		parsed.WhereClause,
		parsed.LimitClause,
		parsed.IsSelectAll,
	)
}
