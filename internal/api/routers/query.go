package routers

import (
	"encoding/json"
	"net/http"
	"strings"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
)

type QueryRouter struct {
	engine datasync.QueryEngine
}

func NewQueryRouter(engine datasync.QueryEngine) *QueryRouter {
	return &QueryRouter{
		engine: engine,
	}
}

func (r *QueryRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /query", r.handleQuery)
}

type QueryRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
}

type QueryResponse struct {
	Rows     []map[string]interface{} `json:"rows"`
	RowCount int                      `json:"rowCount"`
}

func (r *QueryRouter) handleQuery(w http.ResponseWriter, req *http.Request) {
	var queryReq QueryRequest
	if err := json.NewDecoder(req.Body).Decode(&queryReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Strip trailing semicolon (Trino doesn't accept them)
	query := strings.TrimSpace(queryReq.Query)
	query = strings.TrimSuffix(query, ";")

	result, err := r.engine.ExecuteQuery(query, queryReq.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format response to match global query endpoint format
	response := QueryResponse{
		Rows:     result.Rows,
		RowCount: len(result.Rows),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
