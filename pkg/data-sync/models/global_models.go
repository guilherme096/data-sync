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
