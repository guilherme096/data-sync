package query

import (
	"fmt"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

// ResolvedTableSource represents the result of resolving a global table
type ResolvedTableSource struct {
	IsRelation bool

	// For single mapping (Phase 1)
	SingleMapping *models.TableMapping

	// For multiple mappings (Phase 2 - auto UNION)
	MultipleMappings []*models.TableMapping

	// For explicit relations (Phase 2)
	Relation *ResolvedRelation
}

// TableResolver resolves global tables to physical tables
type TableResolver struct {
	storage          storage.MetadataStorage
	relationResolver *RelationResolver
}

// NewTableResolver creates a new table resolver
func NewTableResolver(storage storage.MetadataStorage) *TableResolver {
	return &TableResolver{
		storage:          storage,
		relationResolver: NewRelationResolver(storage),
	}
}

// ResolveGlobalTable resolves a global table name to its physical table mapping
// For Phase 1, only supports single table mapping (errors if multiple mappings exist)
func (r *TableResolver) ResolveGlobalTable(name string) (*models.TableMapping, error) {
	// Check if global table exists
	globalTable, err := r.storage.GetGlobalTable(name)
	if err != nil {
		return nil, fmt.Errorf("global table '%s' not found: %w", name, err)
	}

	if globalTable == nil {
		return nil, fmt.Errorf("global table '%s' not found", name)
	}

	// Get table mappings
	mappings, err := r.storage.ListTableMappings(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get mappings for global table '%s': %w", name, err)
	}

	if len(mappings) == 0 {
		return nil, fmt.Errorf("no table mappings found for global table '%s'", name)
	}

	// Phase 1: Only support single mapping
	if len(mappings) > 1 {
		return nil, fmt.Errorf("global table '%s' has multiple physical table mappings (not supported in Phase 1)", name)
	}

	return mappings[0], nil
}

// ResolveGlobalTableAdvanced resolves a global table to either mappings or a relation
// Phase 2: Supports multiple mappings (auto-UNION) and explicit relations
func (r *TableResolver) ResolveGlobalTableAdvanced(name string) (*ResolvedTableSource, error) {
	// Check if global table exists
	globalTable, err := r.storage.GetGlobalTable(name)
	if err != nil {
		return nil, fmt.Errorf("global table '%s' not found: %w", name, err)
	}

	if globalTable == nil {
		return nil, fmt.Errorf("global table '%s' not found", name)
	}

	// First, check if there's an explicit relation with this name
	relations, err := r.storage.ListTableRelations()
	if err != nil {
		return nil, fmt.Errorf("failed to list relations: %w", err)
	}

	for _, rel := range relations {
		if rel.Name == name {
			// Found a relation with matching name
			resolved, err := r.relationResolver.ResolveRelation(rel.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve relation: %w", err)
			}
			return &ResolvedTableSource{
				IsRelation: true,
				Relation:   resolved,
			}, nil
		}
	}

	// No relation found, check for table mappings
	mappings, err := r.storage.ListTableMappings(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get mappings for global table '%s': %w", name, err)
	}

	if len(mappings) == 0 {
		return nil, fmt.Errorf("no table mappings or relations found for global table '%s'", name)
	}

	if len(mappings) == 1 {
		// Single mapping
		return &ResolvedTableSource{
			IsRelation:    false,
			SingleMapping: mappings[0],
		}, nil
	}

	// Multiple mappings - will auto-generate UNION
	return &ResolvedTableSource{
		IsRelation:       false,
		MultipleMappings: mappings,
	}, nil
}
