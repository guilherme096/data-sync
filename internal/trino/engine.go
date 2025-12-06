package trino

import (
	"database/sql"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	_ "github.com/trinodb/trino-go-client/trino"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(connString string) (*Engine, error) {
	db, err := sql.Open("trino", connString)
	if err != nil {
		return nil, err
	}
	return &Engine{db: db}, nil
}

func (e *Engine) Close() error {
	return e.db.Close()
}

func (e *Engine) ExecuteQuery(query string, params map[string]interface{}) (datasync.QueryResult, error) {
	stmt, err := e.db.Prepare(query)
	if err != nil {
		return datasync.QueryResult{}, err
	}
	defer stmt.Close()

	args := make([]interface{}, 0, len(params))
	for _, v := range params {
		args = append(args, v)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return datasync.QueryResult{}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return datasync.QueryResult{}, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return datasync.QueryResult{}, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return datasync.QueryResult{Rows: results}, nil
}
