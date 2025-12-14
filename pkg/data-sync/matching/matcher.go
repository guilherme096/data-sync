package matching

import "github.com/guilherme096/data-sync/pkg/data-sync/models"

// RelationSuggestion represents a single suggested table relation
type RelationSuggestion struct {
	Name         string
	LeftTable    models.TableSource
	RightTable   models.TableSource
	RelationType string // "JOIN" or "UNION"
	JoinColumn   *models.JoinColumn
	Description  string
	Confidence   float64 // 0.0 to 1.0
}

// MatchingStrategy defines the interface for relation matching strategies
type MatchingStrategy interface {
	// SuggestRelations analyzes available tables and relations to suggest new relations
	SuggestRelations(ctx MatchingContext) ([]RelationSuggestion, error)
}

// MatchingContext provides context for the matching operation
type MatchingContext struct {
	// Available physical tables (from metadata discovery)
	PhysicalTables []PhysicalTableInfo

	// Existing relations that can be used as sources
	ExistingRelations []*models.TableRelation

	// Maximum number of suggestions to return
	MaxSuggestions int
}

// PhysicalTableInfo contains detailed information about a physical table
type PhysicalTableInfo struct {
	Catalog string
	Schema  string
	Table   string
	Columns []ColumnInfo
}

// ColumnInfo contains column metadata
type ColumnInfo struct {
	Name     string
	DataType string
}

// Matcher orchestrates the matching process using a strategy
type Matcher struct {
	strategy MatchingStrategy
}

// NewMatcher creates a new matcher with the given strategy
func NewMatcher(strategy MatchingStrategy) *Matcher {
	return &Matcher{strategy: strategy}
}

// SuggestRelations delegates to the strategy
func (m *Matcher) SuggestRelations(ctx MatchingContext) ([]RelationSuggestion, error) {
	return m.strategy.SuggestRelations(ctx)
}
