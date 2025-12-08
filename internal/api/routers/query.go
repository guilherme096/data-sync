package routers

import (
	"encoding/json"
	"net/http"

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

func (r *QueryRouter) handleQuery(w http.ResponseWriter, req *http.Request) {
	var queryReq QueryRequest
	if err := json.NewDecoder(req.Body).Decode(&queryReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := r.engine.ExecuteQuery(queryReq.Query, queryReq.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
