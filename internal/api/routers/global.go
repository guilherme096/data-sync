package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type GlobalRouter struct {
	storage storage.MetadataStorage
}

func NewGlobalRouter(storage storage.MetadataStorage) *GlobalRouter {
	return &GlobalRouter{
		storage: storage,
	}
}

func (r *GlobalRouter) RegisterRoutes(mux *http.ServeMux) {
	// Global table routes
	mux.HandleFunc("POST /global/tables", r.handleCreateGlobalTable)
	mux.HandleFunc("GET /global/tables", r.handleListGlobalTables)
	mux.HandleFunc("GET /global/tables/{name}", r.handleGetGlobalTable)
	mux.HandleFunc("DELETE /global/tables/{name}", r.handleDeleteGlobalTable)

	// Global column routes
	mux.HandleFunc("POST /global/tables/{name}/columns", r.handleCreateGlobalColumn)
	mux.HandleFunc("GET /global/tables/{name}/columns", r.handleListGlobalColumns)
	mux.HandleFunc("DELETE /global/tables/{name}/columns/{column}", r.handleDeleteGlobalColumn)

	// Table mapping routes
	mux.HandleFunc("POST /global/tables/{name}/mappings/tables", r.handleCreateTableMapping)
	mux.HandleFunc("GET /global/tables/{name}/mappings/tables", r.handleListTableMappings)
	mux.HandleFunc("DELETE /global/tables/{name}/mappings/tables", r.handleDeleteTableMapping)

	// Column mapping routes
	mux.HandleFunc("POST /global/tables/{name}/columns/{column}/mappings", r.handleCreateColumnMapping)
	mux.HandleFunc("GET /global/tables/{name}/columns/{column}/mappings", r.handleListColumnMappings)
	mux.HandleFunc("DELETE /global/tables/{name}/columns/{column}/mappings", r.handleDeleteColumnMapping)
}

// ============================================================================
// Global Table Handlers
// ============================================================================

func (r *GlobalRouter) handleCreateGlobalTable(w http.ResponseWriter, req *http.Request) {
	var table models.GlobalTable
	if err := json.NewDecoder(req.Body).Decode(&table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.storage.CreateGlobalTable(&table); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(table)
}

func (r *GlobalRouter) handleListGlobalTables(w http.ResponseWriter, req *http.Request) {
	tables, err := r.storage.ListGlobalTables()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

func (r *GlobalRouter) handleGetGlobalTable(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	if name == "" {
		http.Error(w, "table name is required", http.StatusBadRequest)
		return
	}

	table, err := r.storage.GetGlobalTable(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

func (r *GlobalRouter) handleDeleteGlobalTable(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	if name == "" {
		http.Error(w, "table name is required", http.StatusBadRequest)
		return
	}

	if err := r.storage.DeleteGlobalTable(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// Global Column Handlers
// ============================================================================

func (r *GlobalRouter) handleCreateGlobalColumn(w http.ResponseWriter, req *http.Request) {
	var column models.GlobalColumn
	if err := json.NewDecoder(req.Body).Decode(&column); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Override with path parameter
	column.GlobalTableName = req.PathValue("name")

	if err := r.storage.CreateGlobalColumn(&column); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(column)
}

func (r *GlobalRouter) handleListGlobalColumns(w http.ResponseWriter, req *http.Request) {
	tableName := req.PathValue("name")
	if tableName == "" {
		http.Error(w, "table name is required", http.StatusBadRequest)
		return
	}

	columns, err := r.storage.ListGlobalColumns(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}

func (r *GlobalRouter) handleDeleteGlobalColumn(w http.ResponseWriter, req *http.Request) {
	tableName := req.PathValue("name")
	columnName := req.PathValue("column")

	if tableName == "" || columnName == "" {
		http.Error(w, "table name and column name are required", http.StatusBadRequest)
		return
	}

	if err := r.storage.DeleteGlobalColumn(tableName, columnName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// Table Mapping Handlers
// ============================================================================

func (r *GlobalRouter) handleCreateTableMapping(w http.ResponseWriter, req *http.Request) {
	var mapping models.TableMapping
	if err := json.NewDecoder(req.Body).Decode(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Override with path parameter
	mapping.GlobalTableName = req.PathValue("name")

	if err := r.storage.CreateTableMapping(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapping)
}

func (r *GlobalRouter) handleListTableMappings(w http.ResponseWriter, req *http.Request) {
	tableName := req.PathValue("name")
	if tableName == "" {
		http.Error(w, "table name is required", http.StatusBadRequest)
		return
	}

	mappings, err := r.storage.ListTableMappings(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mappings)
}

func (r *GlobalRouter) handleDeleteTableMapping(w http.ResponseWriter, req *http.Request) {
	tableName := req.PathValue("name")

	var mapping models.TableMapping
	if err := json.NewDecoder(req.Body).Decode(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.storage.DeleteTableMapping(tableName, mapping.CatalogName, mapping.SchemaName, mapping.TableName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// Column Mapping Handlers
// ============================================================================

func (r *GlobalRouter) handleCreateColumnMapping(w http.ResponseWriter, req *http.Request) {
	var mapping models.ColumnMapping
	if err := json.NewDecoder(req.Body).Decode(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Override with path parameters
	mapping.GlobalTableName = req.PathValue("name")
	mapping.GlobalColumnName = req.PathValue("column")

	if err := r.storage.CreateColumnMapping(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapping)
}

func (r *GlobalRouter) handleListColumnMappings(w http.ResponseWriter, req *http.Request) {
	tableName := req.PathValue("name")
	columnName := req.PathValue("column")

	if tableName == "" || columnName == "" {
		http.Error(w, "table name and column name are required", http.StatusBadRequest)
		return
	}

	mappings, err := r.storage.ListColumnMappings(tableName, columnName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mappings)
}

func (r *GlobalRouter) handleDeleteColumnMapping(w http.ResponseWriter, req *http.Request) {
	tableName := req.PathValue("name")
	columnName := req.PathValue("column")

	var mapping models.ColumnMapping
	if err := json.NewDecoder(req.Body).Decode(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.storage.DeleteColumnMapping(
		tableName, columnName,
		mapping.CatalogName, mapping.SchemaName, mapping.TableName, mapping.ColumnName,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
