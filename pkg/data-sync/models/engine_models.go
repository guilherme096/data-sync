package models

type Catalog struct {
	Name     string
	Metadata map[string]string
}

type Schema struct {
	Name        string
	CatalogName string
	Metadata    map[string]string
}

type Table struct {
	Name        string
	SchemaName  string
	CatalogName string
	Metadata    map[string]string
}

type Column struct {
	Name        string
	TableName   string
	SchemaName  string
	CatalogName string
	DataType    string
	Metadata    map[string]string
}
