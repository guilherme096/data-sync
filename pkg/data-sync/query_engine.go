package datasync

type QueryEngine interface {
	ExecuteQuery(query string, params map[string]interface{}) (QueryResult, error)
}

type QueryResult struct {
	Rows []map[string]interface{}
}

type QueryRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
}
