package matching

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/guilherme096/data-sync/pkg/data-sync/chatbot"
	"github.com/guilherme096/data-sync/pkg/data-sync/models"
	"google.golang.org/genai"
)

type GeminiMatchingStrategy struct {
	client *chatbot.GeminiClient
}

func NewGeminiMatchingStrategy(client *chatbot.GeminiClient) *GeminiMatchingStrategy {
	return &GeminiMatchingStrategy{client: client}
}

func (s *GeminiMatchingStrategy) SuggestRelations(ctx MatchingContext) ([]RelationSuggestion, error) {
	// Build prompt with metadata
	prompt := s.buildPrompt(ctx)

	// Call Gemini with structured output request
	response, err := s.callGeminiForSuggestions(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions from Gemini: %w", err)
	}

	return response, nil
}

func (s *GeminiMatchingStrategy) buildPrompt(ctx MatchingContext) string {
	// Serialize metadata to JSON for clarity
	tablesJSON, _ := json.MarshalIndent(ctx.PhysicalTables, "", "  ")
	relationsJSON, _ := json.MarshalIndent(ctx.ExistingRelations, "", "  ")

	return fmt.Sprintf(`You are an expert data architect analyzing a federated data system.

TASK: Suggest table relations (JOINs and UNIONs) to create unified global tables representing single entities.

IMPORTANT: Both JOINs and UNIONs are used to create a single global table representing ONE entity, NOT to relate different entities together.

AVAILABLE PHYSICAL TABLES:
%s

EXISTING RELATIONS:
%s

GUIDELINES:

1. **UNION Relations** (Vertical Stacking - Same Schema):
   - Use UNIONs to combine tables representing the SAME entity with the SAME schema
   - Example: "customers_2023" UNION "customers_2024" → creates global table "customers"
   - Look for tables with:
     - Similar/identical column names and types
     - Same semantic meaning (temporal partitions, regional splits, etc.)
   - Relation name: The entity name (e.g., "customers", "orders", "products")

2. **JOIN Relations** (Horizontal Merge - Same Entity, Different Sources):
   - Use JOINs to combine DIFFERENT representations of the SAME entity from different catalogs/schemas
   - Example: "postgres.public.customers" JOIN "mysql.main.customers" → creates global table "customers"
   - This merges additional columns/attributes about the same entity from different data sources
   - Join on the most compatible attribute:
     - Primary key columns (typically "id")
     - Columns with matching names and types (e.g., both have "customer_id")
     - If table names are similar but columns differ, join on common identifier
   - **DO NOT** suggest JOINs between different entities (e.g., customers JOIN orders)
   - Relation name: The entity name (e.g., "customers", "users", "products")

3. **Entity Identification**:
   - Identify tables representing the same entity by:
     - Similar table names (e.g., "customers", "customer", "clients")
     - Similar column structures (even if from different schemas)
     - Semantic similarity in naming
   - Each relation should represent a single business entity

4. **Nested Relations**: You can suggest relations that use existing relations as sources, not just physical tables.

5. **Relation Names**:
   - Name should be the entity name, in singular form if possible
   - Examples: "customers" (not "postgres_mysql_customers_join"), "orders", "users", "products"
   - Use lowercase, simple names

6. **Confidence**: Assign confidence scores:
   - 0.9-1.0: Very high confidence (exact name/type matches, clear same entity)
   - 0.7-0.9: High confidence (semantic matches, similar structures)
   - 0.5-0.7: Medium confidence (heuristic matches)
   - Below 0.5: Don't suggest

7. **Limit**: Suggest up to %d relations, prioritized by confidence.

REQUIRED OUTPUT FORMAT (JSON):
{
  "suggestions": [
    {
      "name": "customers",
      "leftTable": {
        "type": "physical",
        "catalog": "postgres",
        "schema": "public",
        "table": "customers"
      },
      "rightTable": {
        "type": "physical",
        "catalog": "mysql",
        "schema": "main",
        "table": "customers"
      },
      "relationType": "JOIN",
      "joinColumn": {
        "left": "id",
        "right": "id"
      },
      "description": "Merges customer data from PostgreSQL and MySQL sources on primary key",
      "confidence": 0.95
    }
  ]
}

Generate ONLY the JSON output, no additional text.`,
		string(tablesJSON),
		string(relationsJSON),
		ctx.MaxSuggestions)
}

// GeminiSuggestionsResponse matches the JSON structure we expect
type GeminiSuggestionsResponse struct {
	Suggestions []struct {
		Name       string `json:"name"`
		LeftTable  struct {
			Type       string `json:"type"`
			Catalog    string `json:"catalog,omitempty"`
			Schema     string `json:"schema,omitempty"`
			Table      string `json:"table,omitempty"`
			RelationID string `json:"relationId,omitempty"`
		} `json:"leftTable"`
		RightTable struct {
			Type       string `json:"type"`
			Catalog    string `json:"catalog,omitempty"`
			Schema     string `json:"schema,omitempty"`
			Table      string `json:"table,omitempty"`
			RelationID string `json:"relationId,omitempty"`
		} `json:"rightTable"`
		RelationType string `json:"relationType"`
		JoinColumn   *struct {
			Left  string `json:"left"`
			Right string `json:"right"`
		} `json:"joinColumn,omitempty"`
		Description string  `json:"description"`
		Confidence  float64 `json:"confidence"`
	} `json:"suggestions"`
}

func (s *GeminiMatchingStrategy) callGeminiForSuggestions(prompt string) ([]RelationSuggestion, error) {
	// Use Gemini's JSON mode for structured output
	ctx := context.Background()

	systemInstruction := "You are a data architecture expert. Always respond with valid JSON matching the requested schema."

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		SystemInstruction: genai.NewContentFromText(
			systemInstruction,
			"system",
		),
	}

	res, err := s.client.SendMessageWithConfig(ctx, prompt, config)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var geminiResp GeminiSuggestionsResponse
	if err := json.Unmarshal([]byte(res), &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	// Convert to our domain model
	suggestions := make([]RelationSuggestion, 0, len(geminiResp.Suggestions))
	for _, gemSuggestion := range geminiResp.Suggestions {
		suggestion := RelationSuggestion{
			Name:         gemSuggestion.Name,
			RelationType: gemSuggestion.RelationType,
			Description:  gemSuggestion.Description,
			Confidence:   gemSuggestion.Confidence,
		}

		// Convert left table
		suggestion.LeftTable = models.TableSource{
			Type:       gemSuggestion.LeftTable.Type,
			Catalog:    gemSuggestion.LeftTable.Catalog,
			Schema:     gemSuggestion.LeftTable.Schema,
			Table:      gemSuggestion.LeftTable.Table,
			RelationID: gemSuggestion.LeftTable.RelationID,
		}

		// Convert right table
		suggestion.RightTable = models.TableSource{
			Type:       gemSuggestion.RightTable.Type,
			Catalog:    gemSuggestion.RightTable.Catalog,
			Schema:     gemSuggestion.RightTable.Schema,
			Table:      gemSuggestion.RightTable.Table,
			RelationID: gemSuggestion.RightTable.RelationID,
		}

		// Convert join column if present
		if gemSuggestion.JoinColumn != nil {
			suggestion.JoinColumn = &models.JoinColumn{
				Left:  gemSuggestion.JoinColumn.Left,
				Right: gemSuggestion.JoinColumn.Right,
			}
		}

		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}
