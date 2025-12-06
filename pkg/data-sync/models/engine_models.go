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
