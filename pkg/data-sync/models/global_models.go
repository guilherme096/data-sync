package models

// GlobalTable represents a logical table that abstracts multiple local tables
type GlobalTable struct {
	Name        string
	Description string
}

// GlobalColumn represents a column in a global table
type GlobalColumn struct {
	GlobalTableName string
	Name            string
	DataType        string
	Description     string
}

// ColumnRelationship represents a foreign key relationship between global table columns
type ColumnRelationship struct {
	SourceGlobalTableName  string
	SourceGlobalColumnName string
	TargetGlobalTableName  string
	TargetGlobalColumnName string
	RelationshipName       string // Optional: user-defined name
	Description            string // Optional: description of the relationship
}
