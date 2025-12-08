package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type CatalogsRouter struct {
	storage storage.MetadataStorage
}

func NewCatalogsRouter(storage storage.MetadataStorage) *CatalogsRouter {
	return &CatalogsRouter{
		storage: storage,
	}
}

func (r *CatalogsRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /catalogs", r.handleListCatalogs)
	mux.HandleFunc("GET /catalogs/{name}", r.handleGetCatalog)
	mux.HandleFunc("GET /catalogs/{name}/schemas", r.handleListSchemas)
}

func (r *CatalogsRouter) handleListCatalogs(w http.ResponseWriter, req *http.Request) {
	catalogs, err := r.storage.ListCatalogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalogs)
}

func (r *CatalogsRouter) handleGetCatalog(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	if name == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}

	catalog, err := r.storage.GetCatalog(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalog)
}

func (r *CatalogsRouter) handleListSchemas(w http.ResponseWriter, req *http.Request) {
	catalogName := req.PathValue("name")
	if catalogName == "" {
		http.Error(w, "catalog name is required", http.StatusBadRequest)
		return
	}

	schemas, err := r.storage.ListSchemas(catalogName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemas)
}
