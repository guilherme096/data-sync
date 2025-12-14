package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type RelationRouter struct {
	storage storage.MetadataStorage
}

func NewRelationRouter(storage storage.MetadataStorage) *RelationRouter {
	return &RelationRouter{
		storage: storage,
	}
}

func (r *RelationRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /relations", r.handleCreateRelation)
	mux.HandleFunc("GET /relations", r.handleListRelations)
	mux.HandleFunc("GET /relations/{id}", r.handleGetRelation)
	mux.HandleFunc("DELETE /relations/{id}", r.handleDeleteRelation)
}

func (r *RelationRouter) handleCreateRelation(w http.ResponseWriter, req *http.Request) {
	var relation models.TableRelation
	if err := json.NewDecoder(req.Body).Decode(&relation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.storage.CreateTableRelation(&relation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(relation)
}

func (r *RelationRouter) handleListRelations(w http.ResponseWriter, req *http.Request) {
	relations, err := r.storage.ListTableRelations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relations)
}

func (r *RelationRouter) handleGetRelation(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		http.Error(w, "relation ID is required", http.StatusBadRequest)
		return
	}

	relation, err := r.storage.GetTableRelation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relation)
}

func (r *RelationRouter) handleDeleteRelation(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		http.Error(w, "relation ID is required", http.StatusBadRequest)
		return
	}

	if err := r.storage.DeleteTableRelation(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
