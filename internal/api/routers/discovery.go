package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
)

type DiscoveryRouter struct {
	discovery discovery.MetadataDiscovery
}

func NewDiscoveryRouter(discovery discovery.MetadataDiscovery) *DiscoveryRouter {
	return &DiscoveryRouter{
		discovery: discovery,
	}
}

func (r *DiscoveryRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /discover/catalogs/{catalog}/schemas/{schema}/tables", r.handleDiscoverTables)
	mux.HandleFunc("GET /discover/catalogs/{catalog}/schemas/{schema}/tables/{table}/columns", r.handleDiscoverColumns)
}

func (r *DiscoveryRouter) handleDiscoverTables(w http.ResponseWriter, req *http.Request) {
	catalogName := req.PathValue("catalog")
	schemaName := req.PathValue("schema")

	if catalogName == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}
	if schemaName == "" {
		http.Error(w, "schema name is required", http.StatusBadRequest)
		return
	}

	tables, err := r.discovery.DiscoverTables(catalogName, schemaName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

func (r *DiscoveryRouter) handleDiscoverColumns(w http.ResponseWriter, req *http.Request) {
	catalogName := req.PathValue("catalog")
	schemaName := req.PathValue("schema")
	tableName := req.PathValue("table")

	if catalogName == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}
	if schemaName == "" {
		http.Error(w, "schema name is required", http.StatusBadRequest)
		return
	}
	if tableName == "" {
		http.Error(w, "table name is required", http.StatusBadRequest)
		return
	}

	columns, err := r.discovery.DiscoverColumns(catalogName, schemaName, tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}
