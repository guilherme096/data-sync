package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/matching"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type RelationRouter struct {
	storage   storage.MetadataStorage
	discovery discovery.MetadataDiscovery
	matcher   *matching.Matcher
}

func NewRelationRouter(storage storage.MetadataStorage, discovery discovery.MetadataDiscovery, matcher *matching.Matcher) *RelationRouter {
	return &RelationRouter{
		storage:   storage,
		discovery: discovery,
		matcher:   matcher,
	}
}

func (r *RelationRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /relations", r.handleCreateRelation)
	mux.HandleFunc("GET /relations", r.handleListRelations)
	mux.HandleFunc("GET /relations/{id}", r.handleGetRelation)
	mux.HandleFunc("DELETE /relations/{id}", r.handleDeleteRelation)
	mux.HandleFunc("POST /relations/auto-match", r.handleAutoMatch)
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

// AutoMatchRequest specifies parameters for auto-matching
type AutoMatchRequest struct {
	MaxSuggestions int  `json:"maxSuggestions"` // Optional, defaults to 5
	AutoCreate     bool `json:"autoCreate"`     // If true, create relations immediately
}

// AutoMatchResponse returns suggestions and optionally created relations
type AutoMatchResponse struct {
	Suggestions      []matching.RelationSuggestion `json:"suggestions"`
	CreatedRelations []*models.TableRelation       `json:"createdRelations,omitempty"`
	Errors           []string                      `json:"errors,omitempty"`
}

func (r *RelationRouter) handleAutoMatch(w http.ResponseWriter, req *http.Request) {
	var matchReq AutoMatchRequest
	if err := json.NewDecoder(req.Body).Decode(&matchReq); err != nil {
		// Use defaults if no body provided or invalid JSON
		matchReq.MaxSuggestions = 5
		matchReq.AutoCreate = true
	}

	if matchReq.MaxSuggestions <= 0 {
		matchReq.MaxSuggestions = 5
	}

	// Gather metadata for matching context
	ctx, err := r.buildMatchingContext(matchReq.MaxSuggestions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to build matching context: %v", err), http.StatusInternalServerError)
		return
	}

	// Get suggestions from matching service
	suggestions, err := r.matcher.SuggestRelations(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get suggestions: %v", err), http.StatusInternalServerError)
		return
	}

	response := AutoMatchResponse{
		Suggestions: suggestions,
	}

	// Auto-create relations if requested
	if matchReq.AutoCreate {
		createdRelations, errors := r.createSuggestedRelations(suggestions)
		response.CreatedRelations = createdRelations
		response.Errors = errors
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (r *RelationRouter) buildMatchingContext(maxSuggestions int) (matching.MatchingContext, error) {
	// Discover all physical tables
	catalogs, err := r.discovery.DiscoverCatalogs()
	if err != nil {
		return matching.MatchingContext{}, err
	}

	var physicalTables []matching.PhysicalTableInfo

	for _, catalog := range catalogs {
		schemas, err := r.discovery.DiscoverSchemas(catalog.Name)
		if err != nil {
			continue // Skip catalogs with errors
		}

		for _, schema := range schemas {
			tables, err := r.discovery.DiscoverTables(catalog.Name, schema.Name)
			if err != nil {
				continue
			}

			for _, table := range tables {
				columns, err := r.discovery.DiscoverColumns(catalog.Name, schema.Name, table.Name)
				if err != nil {
					continue
				}

				columnInfo := make([]matching.ColumnInfo, len(columns))
				for i, col := range columns {
					columnInfo[i] = matching.ColumnInfo{
						Name:     col.Name,
						DataType: col.DataType,
					}
				}

				physicalTables = append(physicalTables, matching.PhysicalTableInfo{
					Catalog: catalog.Name,
					Schema:  schema.Name,
					Table:   table.Name,
					Columns: columnInfo,
				})
			}
		}
	}

	// Get existing relations
	existingRelations, err := r.storage.ListTableRelations()
	if err != nil {
		return matching.MatchingContext{}, err
	}

	return matching.MatchingContext{
		PhysicalTables:    physicalTables,
		ExistingRelations: existingRelations,
		MaxSuggestions:    maxSuggestions,
	}, nil
}

func (r *RelationRouter) createSuggestedRelations(suggestions []matching.RelationSuggestion) ([]*models.TableRelation, []string) {
	var createdRelations []*models.TableRelation
	var errors []string

	existingRelations, _ := r.storage.ListTableRelations()
	existingNames := make(map[string]bool)
	for _, rel := range existingRelations {
		existingNames[rel.Name] = true
	}

	for _, suggestion := range suggestions {
		relation := suggestion.ToTableRelation()

		// Handle duplicate names
		originalName := relation.Name
		suffix := 1
		for existingNames[relation.Name] {
			relation.Name = fmt.Sprintf("%s_v%d", originalName, suffix)
			suffix++
		}

		relation.ID = fmt.Sprintf("auto_%d", time.Now().UnixNano())

		// Validate relation before creating
		if err := r.validateRelation(relation); err != nil {
			errors = append(errors, fmt.Sprintf("Invalid relation '%s': %v", relation.Name, err))
			continue
		}

		if err := r.storage.CreateTableRelation(relation); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to create '%s': %v", relation.Name, err))
			continue
		}

		// Auto-create global table (existing logic)
		if err := r.autoCreateGlobalTableFromRelation(relation); err != nil {
			fmt.Printf("Warning: failed to auto-create global table for '%s': %v\n", relation.Name, err)
		}

		existingNames[relation.Name] = true
		createdRelations = append(createdRelations, relation)
	}

	return createdRelations, errors
}

func (r *RelationRouter) validateRelation(relation *models.TableRelation) error {
	// Validate left table exists
	if relation.LeftTable.Type == "physical" {
		if _, err := r.discovery.DiscoverColumns(
			relation.LeftTable.Catalog,
			relation.LeftTable.Schema,
			relation.LeftTable.Table,
		); err != nil {
			return fmt.Errorf("left table not found: %w", err)
		}
	}

	// Validate right table exists
	if relation.RightTable.Type == "physical" {
		if _, err := r.discovery.DiscoverColumns(
			relation.RightTable.Catalog,
			relation.RightTable.Schema,
			relation.RightTable.Table,
		); err != nil {
			return fmt.Errorf("right table not found: %w", err)
		}
	}

	return nil
}
