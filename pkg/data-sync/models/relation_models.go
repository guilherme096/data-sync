package models

type TableSource struct {
	Type       string `json:"type"` // "physical" or "relation"
	Catalog    string `json:"catalog,omitempty"`
	Schema     string `json:"schema,omitempty"`
	Table      string `json:"table,omitempty"`
	RelationID string `json:"relationId,omitempty"`
}

type JoinColumn struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

type TableRelation struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	LeftTable    TableSource  `json:"leftTable"`
	RightTable   TableSource  `json:"rightTable"`
	RelationType string       `json:"relationType"` // "JOIN" or "UNION"
	JoinColumn   *JoinColumn  `json:"joinColumn,omitempty"`
	Description  string       `json:"description,omitempty"`
}
