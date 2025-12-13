package query

import (
	"fmt"
	"regexp"
	"strings"
)

// ParsedQuery represents a simple parsed SQL query
type ParsedQuery struct {
	TableName    string
	Columns      []string
	WhereClause  string
	LimitClause  string
	IsSelectAll  bool
}

// QueryParser handles parsing of SQL queries
type QueryParser struct{}

// NewQueryParser creates a new query parser
func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

// Parse extracts table name, columns, and WHERE clause from a simple SELECT query
// For Phase 1, supports only: SELECT col1, col2 FROM table WHERE condition
func (p *QueryParser) Parse(query string) (*ParsedQuery, error) {
	query = strings.TrimSpace(query)

	// Basic validation
	if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		return nil, fmt.Errorf("only SELECT queries are supported")
	}

	// Remove trailing semicolon if present
	query = strings.TrimSuffix(query, ";")

	// Extract FROM clause and table name
	fromRegex := regexp.MustCompile(`(?i)\s+FROM\s+([a-zA-Z0-9_]+)`)
	fromMatch := fromRegex.FindStringSubmatch(query)
	if len(fromMatch) < 2 {
		return nil, fmt.Errorf("invalid query: FROM clause not found")
	}
	tableName := fromMatch[1]

	// Extract SELECT columns
	selectRegex := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM`)
	selectMatch := selectRegex.FindStringSubmatch(query)
	if len(selectMatch) < 2 {
		return nil, fmt.Errorf("invalid query: SELECT clause not found")
	}

	columnsStr := strings.TrimSpace(selectMatch[1])
	isSelectAll := columnsStr == "*"

	var columns []string
	if !isSelectAll {
		// Split columns by comma
		for _, col := range strings.Split(columnsStr, ",") {
			columns = append(columns, strings.TrimSpace(col))
		}
	}

	// Extract LIMIT clause (optional)
	limitClause := ""
	limitRegex := regexp.MustCompile(`(?i)\s+LIMIT\s+(\d+)`)
	limitMatch := limitRegex.FindStringSubmatch(query)
	if len(limitMatch) >= 2 {
		limitClause = limitMatch[1]
	}

	// Extract WHERE clause (optional) - must come before LIMIT
	whereClause := ""
	whereRegex := regexp.MustCompile(`(?i)\s+WHERE\s+(.+?)(?:\s+LIMIT\s+\d+)?$`)
	whereMatch := whereRegex.FindStringSubmatch(query)
	if len(whereMatch) >= 2 {
		whereClause = strings.TrimSpace(whereMatch[1])
	}

	return &ParsedQuery{
		TableName:   tableName,
		Columns:     columns,
		WhereClause: whereClause,
		LimitClause: limitClause,
		IsSelectAll: isSelectAll,
	}, nil
}
