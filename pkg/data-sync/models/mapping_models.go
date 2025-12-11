package models

// TableMapping links a local table to a global table
type TableMapping struct {
	GlobalTableName string
	CatalogName     string
	SchemaName      string
	TableName       string
}

// ColumnMapping links a local column to a global column
type ColumnMapping struct {
	GlobalTableName  string
	GlobalColumnName string
	CatalogName      string
	SchemaName       string
	TableName        string
	ColumnName       string
}
