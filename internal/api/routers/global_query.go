package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/query"
)

type GlobalQueryRouter struct {
	translator query.QueryTranslator
}

func NewGlobalQueryRouter(translator query.QueryTranslator) *GlobalQueryRouter {
	return &GlobalQueryRouter{
		translator: translator,
	}
}

func (r *GlobalQueryRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /query/global", r.handleGlobalQuery)
}

type GlobalQueryRequest struct {
	Query string `json:"query"`
}

func (r *GlobalQueryRouter) handleGlobalQuery(w http.ResponseWriter, req *http.Request) {
	var queryReq GlobalQueryRequest
	if err := json.NewDecoder(req.Body).Decode(&queryReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if queryReq.Query == "" {
		http.Error(w, "query field is required", http.StatusBadRequest)
		return
	}

	result, err := r.translator.TranslateAndExecute(queryReq.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
