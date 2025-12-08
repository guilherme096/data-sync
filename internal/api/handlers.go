package api

import (
	"encoding/json"
	_ "github.com/guilherme096/data-sync/pkg/data-sync"
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

type QueryRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.engine.ExecuteQuery(req.Query, req.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleListCatalogs(w http.ResponseWriter, r *http.Request) {
	catalogs, err := s.storage.ListCatalogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalogs)
}

func (s *Server) handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}

	catalog, err := s.storage.GetCatalog(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalog)
}

func (s *Server) handleListSchemas(w http.ResponseWriter, r *http.Request) {
	catalogName := r.PathValue("name")
	if catalogName == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}

	schemas, err := s.storage.ListSchemas(catalogName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemas)
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	if err := s.sync.SyncAll(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Metadata sync completed successfully",
	})
}

func (s *Server) handleDiscoverTables(w http.ResponseWriter, r *http.Request) {
	catalogName := r.PathValue("catalog")
	schemaName := r.PathValue("schema")

	if catalogName == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}
	if schemaName == "" {
		http.Error(w, "schema name is required", http.StatusBadRequest)
		return
	}

	tables, err := s.discovery.DiscoverTables(catalogName, schemaName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

func (s *Server) handleDiscoverColumns(w http.ResponseWriter, r *http.Request) {
	catalogName := r.PathValue("catalog")
	schemaName := r.PathValue("schema")
	tableName := r.PathValue("table")

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

	columns, err := s.discovery.DiscoverColumns(catalogName, schemaName, tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}
