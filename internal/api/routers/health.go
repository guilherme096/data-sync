package routers

import (
	"net/http"
)

type HealthRouter struct{}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", r.handleHealth)
}

func (r *HealthRouter) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
