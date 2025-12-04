package models

type Catalog struct {
	ID       int
	Name     string
	metadata string
}

type Schema struct {
	ID        int
	Name      string
	CatalogID int
	metadata  string
}
