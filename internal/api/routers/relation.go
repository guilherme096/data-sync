package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type RelationRouter struct {
	storage   storage.MetadataStorage
	discovery discovery.MetadataDiscovery
}

func NewRelationRouter(storage storage.MetadataStorage, discovery discovery.MetadataDiscovery) *RelationRouter {
	return &RelationRouter{
		storage:   storage,
		discovery: discovery,
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

	// Create the table relation
	if err := r.storage.CreateTableRelation(&relation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Automatically create global table and column mappings
	if err := r.autoCreateGlobalTableFromRelation(&relation); err != nil {
		// Log error but don't fail the relation creation
		fmt.Printf("Warning: failed to auto-create global table for relation '%s': %v\n", relation.Name, err)
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

// autoCreateGlobalTableFromRelation automatically creates a global table and column mappings
// when a relation is created
func (r *RelationRouter) autoCreateGlobalTableFromRelation(relation *models.TableRelation) error {
	// Check if global table already exists
	existingTable, _ := r.storage.GetGlobalTable(relation.Name)
	if existingTable != nil {
		// Table already exists, skip creation
		return nil
	}

	// Create global table
	globalTable := &models.GlobalTable{
		Name:        relation.Name,
		Description: fmt.Sprintf("Auto-generated from %s relation", relation.RelationType),
	}
	if err := r.storage.CreateGlobalTable(globalTable); err != nil {
		return fmt.Errorf("failed to create global table: %w", err)
	}

	// Discover and create columns from physical tables
	if err := r.discoverAndCreateColumnsFromRelation(relation); err != nil {
		return fmt.Errorf("failed to discover columns: %w", err)
	}

	return nil
}

// discoverAndCreateColumnsFromRelation discovers columns from physical tables in the relation
func (r *RelationRouter) discoverAndCreateColumnsFromRelation(relation *models.TableRelation) error {
	// Collect all physical tables from the relation
	var physicalTables []struct {
		catalog string
		schema  string
		table   string
	}

	// Add left table if physical
	if relation.LeftTable.Type == "physical" {
		physicalTables = append(physicalTables, struct {
			catalog string
			schema  string
			table   string
		}{
			catalog: relation.LeftTable.Catalog,
			schema:  relation.LeftTable.Schema,
			table:   relation.LeftTable.Table,
		})
	}

	// Add right table if physical
	if relation.RightTable.Type == "physical" {
		physicalTables = append(physicalTables, struct {
			catalog string
			schema  string
			table   string
		}{
			catalog: relation.RightTable.Catalog,
			schema:  relation.RightTable.Schema,
			table:   relation.RightTable.Table,
		})
	}

	if len(physicalTables) == 0 {
		return fmt.Errorf("no physical tables found in relation")
	}

	// Discover columns from the first table and create global columns
	firstTable := physicalTables[0]
	columns, err := r.discovery.DiscoverColumns(firstTable.catalog, firstTable.schema, firstTable.table)
	if err != nil {
		return fmt.Errorf("failed to discover columns from %s.%s.%s: %w",
			firstTable.catalog, firstTable.schema, firstTable.table, err)
	}

	// Create global columns
	for _, col := range columns {
		globalColumn := &models.GlobalColumn{
			GlobalTableName: relation.Name,
			Name:            col.Name,
			DataType:        col.DataType,
			Description:     fmt.Sprintf("Auto-discovered from %s.%s.%s", firstTable.catalog, firstTable.schema, firstTable.table),
		}

		if err := r.storage.CreateGlobalColumn(globalColumn); err != nil {
			// Continue even if column creation fails (might already exist)
			fmt.Printf("Warning: failed to create global column '%s': %v\n", col.Name, err)
			continue
		}

		// Create column mappings for all physical tables
		for _, physTable := range physicalTables {
			mapping := &models.ColumnMapping{
				GlobalTableName:  relation.Name,
				GlobalColumnName: col.Name,
				CatalogName:      physTable.catalog,
				SchemaName:       physTable.schema,
				TableName:        physTable.table,
				ColumnName:       col.Name,
			}

			if err := r.storage.CreateColumnMapping(mapping); err != nil {
				fmt.Printf("Warning: failed to create column mapping for '%s' in %s.%s.%s: %v\n",
					col.Name, physTable.catalog, physTable.schema, physTable.table, err)
			}
		}
	}

	return nil
}
