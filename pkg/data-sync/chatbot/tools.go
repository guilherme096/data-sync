package chatbot

import (
	"fmt"

	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/query"
	"google.golang.org/genai"
)

// ToolExecutor executes tools called by the Gemini agent
type ToolExecutor interface {
	ExecuteTool(toolName string, arguments map[string]interface{}) (interface{}, error)
}

// DefaultToolExecutor implements ToolExecutor with access to query translator and metadata discovery
type DefaultToolExecutor struct {
	translator query.QueryTranslator
	discovery  discovery.MetadataDiscovery
	storage    interface {
		ListGlobalTables() ([]*models.GlobalTable, error)
	}
}

// NewToolExecutor creates a new tool executor with required dependencies
func NewToolExecutor(translator query.QueryTranslator, discovery discovery.MetadataDiscovery, storage interface {
	ListGlobalTables() ([]*models.GlobalTable, error)
}) ToolExecutor {
	return &DefaultToolExecutor{
		translator: translator,
		discovery:  discovery,
		storage:    storage,
	}
}

// ExecuteTool routes tool calls to the appropriate handler
func (te *DefaultToolExecutor) ExecuteTool(toolName string, arguments map[string]interface{}) (interface{}, error) {
	switch toolName {
	case "executeGlobalQuery":
		return te.executeGlobalQuery(arguments)
	case "discoverMetadata":
		return te.discoverMetadata(arguments)
	case "listGlobalTables":
		return te.listGlobalTables(arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

// executeGlobalQuery executes a SQL query on global tables
func (te *DefaultToolExecutor) executeGlobalQuery(args map[string]interface{}) (interface{}, error) {
	query, ok := args["query"].(string)
	if !ok {
		return map[string]interface{}{
			"error":      "Invalid query parameter",
			"suggestion": "Please provide a valid SQL query string",
		}, nil
	}

	result, err := te.translator.TranslateAndExecute(query)
	if err != nil {
		// Return structured error that Gemini can explain to user
		return map[string]interface{}{
			"error":      err.Error(),
			"suggestion": "Check that the table names and column names are correct, and the SQL syntax is valid",
		}, nil
	}

	return result, nil
}

// listGlobalTables lists all available global tables in the system
func (te *DefaultToolExecutor) listGlobalTables(args map[string]interface{}) (interface{}, error) {
	tables, err := te.storage.ListGlobalTables()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, nil
	}

	// Format the response
	tableList := make([]map[string]string, len(tables))
	for i, table := range tables {
		tableList[i] = map[string]string{
			"name":        table.Name,
			"description": table.Description,
		}
	}

	return map[string]interface{}{
		"tables": tableList,
		"count":  len(tables),
	}, nil
}

// discoverMetadata discovers metadata about data sources
func (te *DefaultToolExecutor) discoverMetadata(args map[string]interface{}) (interface{}, error) {
	level, ok := args["level"].(string)
	if !ok {
		return map[string]interface{}{
			"error":      "Invalid level parameter",
			"suggestion": "Please specify one of: catalogs, schemas, tables, columns",
		}, nil
	}

	switch level {
	case "catalogs":
		catalogs, err := te.discovery.DiscoverCatalogs()
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}, nil
		}
		return map[string]interface{}{
			"level":    "catalogs",
			"catalogs": catalogs,
		}, nil

	case "schemas":
		catalog, ok := args["catalog"].(string)
		if !ok {
			return map[string]interface{}{
				"error":      "Missing catalog parameter",
				"suggestion": "Please specify a catalog name",
			}, nil
		}
		schemas, err := te.discovery.DiscoverSchemas(catalog)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}, nil
		}
		return map[string]interface{}{
			"level":   "schemas",
			"catalog": catalog,
			"schemas": schemas,
		}, nil

	case "tables":
		catalog, catalogOk := args["catalog"].(string)
		schema, schemaOk := args["schema"].(string)
		if !catalogOk || !schemaOk {
			return map[string]interface{}{
				"error":      "Missing catalog or schema parameter",
				"suggestion": "Please specify both catalog and schema names",
			}, nil
		}
		tables, err := te.discovery.DiscoverTables(catalog, schema)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}, nil
		}
		return map[string]interface{}{
			"level":   "tables",
			"catalog": catalog,
			"schema":  schema,
			"tables":  tables,
		}, nil

	case "columns":
		catalog, catalogOk := args["catalog"].(string)
		schema, schemaOk := args["schema"].(string)
		table, tableOk := args["table"].(string)
		if !catalogOk || !schemaOk || !tableOk {
			return map[string]interface{}{
				"error":      "Missing catalog, schema, or table parameter",
				"suggestion": "Please specify catalog, schema, and table names",
			}, nil
		}
		columns, err := te.discovery.DiscoverColumns(catalog, schema, table)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}, nil
		}
		return map[string]interface{}{
			"level":   "columns",
			"catalog": catalog,
			"schema":  schema,
			"table":   table,
			"columns": columns,
		}, nil

	default:
		return map[string]interface{}{
			"error":      fmt.Sprintf("Invalid level: %s", level),
			"suggestion": "Please use one of: catalogs, schemas, tables, columns",
		}, nil
	}
}

// BuildToolDeclarations returns the tool declarations for Gemini function calling
func BuildToolDeclarations() []*genai.Tool {
	return []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				buildListGlobalTablesTool(),
				buildExecuteGlobalQueryTool(),
				buildDiscoverMetadataTool(),
			},
		},
	}
}

// buildListGlobalTablesTool creates the tool declaration for listing global tables
func buildListGlobalTablesTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "listGlobalTables",
		Description: "Lists all available global tables in the system. Use this tool when the user asks 'what tables do I have?', 'show me available tables', or when you need to know what tables exist before querying. This is the FIRST tool you should use when unsure what data is available.",
		Parameters: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: map[string]*genai.Schema{},
		},
	}
}

// buildExecuteGlobalQueryTool creates the tool declaration for executing queries
func buildExecuteGlobalQueryTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "executeGlobalQuery",
		Description: "Executes a SQL query on global tables. Use this when the user wants to retrieve, count, filter, or analyze data from global tables. The query should use global table names (e.g., 'SELECT * FROM customers'). This tool will automatically translate the global query to the appropriate physical tables and execute it.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"query": {
					Type:        genai.TypeString,
					Description: "SQL query using global table names. Example: 'SELECT * FROM customers WHERE country = 'USA' LIMIT 10'. Supports SELECT, WHERE, JOIN, GROUP BY, ORDER BY, and LIMIT clauses.",
				},
			},
			Required: []string{"query"},
		},
	}
}

// buildDiscoverMetadataTool creates the tool declaration for discovering metadata
func buildDiscoverMetadataTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "discoverMetadata",
		Description: "Discovers metadata about data sources including catalogs, schemas, tables, and columns. Use this when the user asks about available data sources, table structures, or column information. This helps users understand what data they can query.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"level": {
					Type:        genai.TypeString,
					Description: "Level of metadata to discover. Options: 'catalogs' (lists all data source catalogs), 'schemas' (lists schemas in a catalog), 'tables' (lists tables in a schema), 'columns' (lists columns in a table)",
					Enum:        []string{"catalogs", "schemas", "tables", "columns"},
				},
				"catalog": {
					Type:        genai.TypeString,
					Description: "Catalog name (required for schemas, tables, columns levels). Example: 'postgresql', 'mysql'",
				},
				"schema": {
					Type:        genai.TypeString,
					Description: "Schema name (required for tables, columns levels). Example: 'public', 'testdb'",
				},
				"table": {
					Type:        genai.TypeString,
					Description: "Table name (required for columns level). Example: 'customers', 'orders'",
				},
			},
			Required: []string{"level"},
		},
	}
}
