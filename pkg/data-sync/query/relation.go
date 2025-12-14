package query

import (
	"fmt"

	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

// RelationNodeType represents the type of a relation node
type RelationNodeType string

const (
	NodeTypePhysical RelationNodeType = "physical"
	NodeTypeRelation RelationNodeType = "relation"
)

// RelationNode represents a node in the relation tree
type RelationNode struct {
	Type RelationNodeType

	// For physical tables
	Catalog string
	Schema  string
	Table   string

	// For relations (nested)
	RelationID   string
	RelationType string // "JOIN" or "UNION"
}

// ResolvedRelation represents a fully resolved relation tree
type ResolvedRelation struct {
	ID           string
	Name         string
	RelationType string // "JOIN" or "UNION"
	LeftNode     *RelationNode
	RightNode    *RelationNode
	JoinColumn   *models.JoinColumn // Only for JOIN relations
}

// RelationResolver resolves table relations to their physical tables
type RelationResolver struct {
	storage storage.MetadataStorage
}

// NewRelationResolver creates a new relation resolver
func NewRelationResolver(storage storage.MetadataStorage) *RelationResolver {
	return &RelationResolver{
		storage: storage,
	}
}

// ResolveRelation resolves a table relation to its components
// Phase 3: Supports both UNION and JOIN with physical tables (no nesting yet)
func (r *RelationResolver) ResolveRelation(relationID string) (*ResolvedRelation, error) {
	// Fetch the relation from storage
	relation, err := r.storage.GetTableRelation(relationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relation '%s': %w", relationID, err)
	}

	if relation == nil {
		return nil, fmt.Errorf("relation '%s' not found", relationID)
	}

	// Support both UNION and JOIN in Phase 3
	if relation.RelationType != "UNION" && relation.RelationType != "JOIN" {
		return nil, fmt.Errorf("relation type '%s' not supported (only UNION and JOIN)", relation.RelationType)
	}

	// For JOIN, validate that join columns are specified
	if relation.RelationType == "JOIN" {
		if relation.JoinColumn == nil {
			return nil, fmt.Errorf("JOIN relation requires join columns to be specified")
		}
		if relation.JoinColumn.Left == "" || relation.JoinColumn.Right == "" {
			return nil, fmt.Errorf("JOIN relation requires both left and right join columns")
		}
	}

	// Resolve left node
	leftNode, err := r.resolveTableSource(&relation.LeftTable)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve left table: %w", err)
	}

	// Resolve right node
	rightNode, err := r.resolveTableSource(&relation.RightTable)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve right table: %w", err)
	}

	return &ResolvedRelation{
		ID:           relation.ID,
		Name:         relation.Name,
		RelationType: relation.RelationType,
		LeftNode:     leftNode,
		RightNode:    rightNode,
		JoinColumn:   relation.JoinColumn,
	}, nil
}

// resolveTableSource resolves a TableSource to a RelationNode
func (r *RelationResolver) resolveTableSource(source *models.TableSource) (*RelationNode, error) {
	return r.resolveTableSourceWithVisited(source, make(map[string]bool))
}

// resolveTableSourceWithVisited resolves a TableSource with circular reference detection
// Phase 4: Supports nested relations recursively
func (r *RelationResolver) resolveTableSourceWithVisited(source *models.TableSource, visited map[string]bool) (*RelationNode, error) {
	switch source.Type {
	case "physical":
		// Physical table - return as-is
		return &RelationNode{
			Type:    NodeTypePhysical,
			Catalog: source.Catalog,
			Schema:  source.Schema,
			Table:   source.Table,
		}, nil

	case "relation":
		// Phase 4: Support nested relations
		if source.RelationID == "" {
			return nil, fmt.Errorf("relation source requires relationId")
		}

		// Check for circular reference
		if visited[source.RelationID] {
			return nil, fmt.Errorf("circular relation detected: relation '%s' references itself", source.RelationID)
		}

		// Mark as visited
		visited[source.RelationID] = true

		// Recursively resolve the nested relation
		nestedRelation, err := r.ResolveRelationWithVisited(source.RelationID, visited)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve nested relation '%s': %w", source.RelationID, err)
		}

		// Return as a relation node
		return &RelationNode{
			Type:         NodeTypeRelation,
			RelationID:   nestedRelation.ID,
			RelationType: nestedRelation.RelationType,
		}, nil

	default:
		return nil, fmt.Errorf("unknown table source type: %s", source.Type)
	}
}

// ResolveRelationWithVisited resolves a relation with circular reference tracking
// Phase 4: Supports nested relations recursively
func (r *RelationResolver) ResolveRelationWithVisited(relationID string, visited map[string]bool) (*ResolvedRelation, error) {
	// Fetch the relation from storage
	relation, err := r.storage.GetTableRelation(relationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relation '%s': %w", relationID, err)
	}

	if relation == nil {
		return nil, fmt.Errorf("relation '%s' not found", relationID)
	}

	// Support both UNION and JOIN
	if relation.RelationType != "UNION" && relation.RelationType != "JOIN" {
		return nil, fmt.Errorf("relation type '%s' not supported (only UNION and JOIN)", relation.RelationType)
	}

	// For JOIN, validate that join columns are specified
	if relation.RelationType == "JOIN" {
		if relation.JoinColumn == nil {
			return nil, fmt.Errorf("JOIN relation requires join columns to be specified")
		}
		if relation.JoinColumn.Left == "" || relation.JoinColumn.Right == "" {
			return nil, fmt.Errorf("JOIN relation requires both left and right join columns")
		}
	}

	// Resolve left node (may be nested)
	leftNode, err := r.resolveTableSourceWithVisited(&relation.LeftTable, visited)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve left table: %w", err)
	}

	// Resolve right node (may be nested)
	rightNode, err := r.resolveTableSourceWithVisited(&relation.RightTable, visited)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve right table: %w", err)
	}

	return &ResolvedRelation{
		ID:           relation.ID,
		Name:         relation.Name,
		RelationType: relation.RelationType,
		LeftNode:     leftNode,
		RightNode:    rightNode,
		JoinColumn:   relation.JoinColumn,
	}, nil
}

// GetPhysicalTables extracts all physical tables from a resolved relation
func (r *RelationResolver) GetPhysicalTables(resolved *ResolvedRelation) ([]*RelationNode, error) {
	tables := []*RelationNode{}

	if resolved.LeftNode.Type == NodeTypePhysical {
		tables = append(tables, resolved.LeftNode)
	}

	if resolved.RightNode.Type == NodeTypePhysical {
		tables = append(tables, resolved.RightNode)
	}

	return tables, nil
}
