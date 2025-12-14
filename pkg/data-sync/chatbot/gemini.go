package chatbot

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	ctx    context.Context
}

func NewGeminiClient() (*GeminiClient, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &GeminiClient{client: client, ctx: ctx}, nil
}

func (g *GeminiClient) SendMessage(message string) (string, error) {
	res, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", genai.Text(message), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return res.Text(), nil
}

func (g *GeminiClient) SendMessageWithHistory(message string, history []ChatMessage) (string, error) {
	// Convert history to Gemini format
	var contents []*genai.Content

	for _, msg := range history {
		// Map frontend roles to Gemini roles
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		// Create content for this message with the correct role
		msgContents := genai.Text(msg.Content)
		if len(msgContents) > 0 {
			msgContents[0].Role = role
			contents = append(contents, msgContents...)
		}
	}

	// Add current user message
	userContents := genai.Text(message)
	if len(userContents) > 0 {
		userContents[0].Role = "user"
		contents = append(contents, userContents...)
	}

	// Generate response with conversation context
	res, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with history: %w", err)
	}

	return res.Text(), nil
}

func (g *GeminiClient) SendMessageWithTools(message string, history []ChatMessage, toolExecutor ToolExecutor) (*AgentResponse, error) {
	// Convert history to Gemini format
	var contents []*genai.Content

	for _, msg := range history {
		// Map frontend roles to Gemini roles
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		// Create content for this message with the correct role
		msgContents := genai.Text(msg.Content)
		if len(msgContents) > 0 {
			msgContents[0].Role = role
			contents = append(contents, msgContents...)
		}
	}

	// Add current user message
	userContents := genai.Text(message)
	if len(userContents) > 0 {
		userContents[0].Role = "user"
		contents = append(contents, userContents...)
	}

	// Build tool declarations
	tools := BuildToolDeclarations()

	// Track tool results for response
	var toolResults []ToolResult

	// System instruction to guide Gemini on how to use the tools
	systemInstruction := `You are a helpful data assistant with access to a federated data system.

IMPORTANT GUIDELINES:
1. When users ask what tables are available (e.g., "what tables do I have?", "show me available tables"), ALWAYS use listGlobalTables first.
2. When users ask to see, query, or analyze data (e.g., "show me clients", "how many orders"), use executeGlobalQuery with SQL queries on GLOBAL TABLES.
3. If you encounter an error saying a table doesn't exist, use listGlobalTables to see what tables are actually available.
4. Global tables abstract physical data sources - users query logical table names like "clients" or "orders", not physical catalog.schema.table names.
5. Always use executeGlobalQuery for data retrieval queries - construct proper SQL SELECT statements.
6. The discoverMetadata tool is for exploring physical catalogs/schemas/tables/columns - use it when users ask about the underlying data sources.
7. Provide friendly, conversational responses that explain the data you found.

Example interactions:
- "Show me all clients" → listGlobalTables (to verify "clients" exists), then executeGlobalQuery with "SELECT * FROM clients"
- "How many orders are there?" → executeGlobalQuery with "SELECT COUNT(*) FROM orders"
- "What tables do I have?" → listGlobalTables
- "Show me clients from USA" → executeGlobalQuery with "SELECT * FROM clients WHERE country = 'USA'"
- "What catalogs exist?" → discoverMetadata with level="catalogs"`

	// Generate content with tools - may require multiple rounds
	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		// Generate response with tools
		res, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{
			Tools:             tools,
			SystemInstruction: genai.NewContentFromText(systemInstruction, "system"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate content with tools: %w", err)
		}

		// Check if Gemini wants to call a function
		if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
			return &AgentResponse{
				Message:     "",
				ToolResults: toolResults,
			}, nil
		}

		// Look for function calls in the response
		var hasFunctionCalls bool
		for _, part := range res.Candidates[0].Content.Parts {
			if part.FunctionCall != nil {
				hasFunctionCalls = true
				break
			}
		}

		// If no function calls, we have the final response
		if !hasFunctionCalls {
			return &AgentResponse{
				Message:     res.Text(),
				ToolResults: toolResults,
			}, nil
		}

		// Execute each function call and collect responses
		for _, part := range res.Candidates[0].Content.Parts {
			if fc := part.FunctionCall; fc != nil {
				// Execute the tool
				result, err := toolExecutor.ExecuteTool(fc.Name, fc.Args)
				if err != nil {
					// Create error response
					result = map[string]interface{}{
						"error": err.Error(),
					}
				}

				// Track tool result
				toolResults = append(toolResults, ToolResult{
					ToolName: fc.Name,
					Data:     result,
				})

				// Convert result to map[string]any for FunctionResponse
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					// If result is not a map, wrap it
					resultMap = map[string]interface{}{
						"result": result,
					}
				}

				// Add function response to conversation
				functionContent := genai.NewContentFromFunctionResponse(fc.Name, resultMap, "function")
				contents = append(contents, functionContent)
			}
		}

		// Continue to next iteration to get Gemini's response based on function results
	}

	// Max iterations reached
	return &AgentResponse{
		Message:     "I apologize, but I needed to use too many tools to answer your question. Please try rephrasing or breaking down your request.",
		ToolResults: toolResults,
	}, nil
}

func (g *GeminiClient) SendMessageForQueryGeneration(message string, history []ChatMessage, toolExecutor ToolExecutor) (*QueryGenerationResponse, error) {
	// Convert history to Gemini format
	var contents []*genai.Content

	for _, msg := range history {
		// Map frontend roles to Gemini roles
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		// Create content for this message with the correct role
		msgContents := genai.Text(msg.Content)
		if len(msgContents) > 0 {
			msgContents[0].Role = role
			contents = append(contents, msgContents...)
		}
	}

	// Add current user message
	userContents := genai.Text(message)
	if len(userContents) > 0 {
		userContents[0].Role = "user"
		contents = append(contents, userContents...)
	}

	// Build query generator tool declarations (no execution tools)
	tools := BuildQueryGeneratorToolDeclarations()

	// System instruction focused on query generation
	systemInstruction := `You are an expert SQL query generator for a federated data system with global tables.

IMPORTANT GUIDELINES:
1. Your ONLY job is to generate SQL queries - NEVER execute them.
2. When users ask for queries, ALWAYS use listGlobalTables first to see what tables are available.
3. Use getTableColumns to understand what columns exist in the tables you want to query.
4. Generate queries using GLOBAL TABLE names (e.g., "clients", "orders"), not physical catalog.schema.table names.
5. Your final response MUST include a valid SQL query in a code block (triple backticks with sql).
6. Provide a brief explanation of what the query does.
7. Use discoverMetadata only if the user specifically asks about physical database structure.

QUERY FORMAT:
Always return your SQL query in this format:
` + "```sql\n" + `SELECT column1, column2 FROM table_name WHERE condition
` + "```" + `

Example interactions:
- "Get all clients" → listGlobalTables, getTableColumns("clients"), then respond with:
  "Here's a query to get all clients:
  ` + "```sql\nSELECT * FROM clients\n```" + `"

- "Top 5 customers by revenue" → listGlobalTables, getTableColumns("clients"), then respond with:
  "Here's a query to get the top 5 customers by revenue:
  ` + "```sql\nSELECT name, revenue FROM clients ORDER BY revenue DESC LIMIT 5\n```" + `"

- "Count orders from USA" → listGlobalTables, getTableColumns("orders"), then respond with:
  "Here's a query to count orders from the USA:
  ` + "```sql\nSELECT COUNT(*) as order_count FROM orders WHERE country = 'USA'\n```" + `"`

	// Generate content with tools - may require multiple rounds
	maxIterations := 5
	var lastResponse string
	for i := 0; i < maxIterations; i++ {
		// Generate response with tools
		res, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{
			Tools:             tools,
			SystemInstruction: genai.NewContentFromText(systemInstruction, "system"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate content with tools: %w", err)
		}

		// Check if Gemini wants to call a function
		if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
			return &QueryGenerationResponse{
				Message:      lastResponse,
				GeneratedSQL: "",
			}, nil
		}

		// Look for function calls in the response
		var hasFunctionCalls bool
		for _, part := range res.Candidates[0].Content.Parts {
			if part.FunctionCall != nil {
				hasFunctionCalls = true
				break
			}
		}

		// If no function calls, we have the final response
		if !hasFunctionCalls {
			responseText := res.Text()
			generatedSQL := extractSQLFromResponse(responseText)

			// Debug logging
			fmt.Printf("[Query Generation] Response text: %s\n", responseText)
			fmt.Printf("[Query Generation] Extracted SQL: %s\n", generatedSQL)

			return &QueryGenerationResponse{
				Message:      responseText,
				GeneratedSQL: generatedSQL,
			}, nil
		}

		// Execute each function call and collect responses
		for _, part := range res.Candidates[0].Content.Parts {
			if fc := part.FunctionCall; fc != nil {
				// Execute the tool
				result, err := toolExecutor.ExecuteTool(fc.Name, fc.Args)
				if err != nil {
					// Create error response
					result = map[string]interface{}{
						"error": err.Error(),
					}
				}

				// Convert result to map[string]any for FunctionResponse
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					// If result is not a map, wrap it
					resultMap = map[string]interface{}{
						"result": result,
					}
				}

				// Add function response to conversation
				functionContent := genai.NewContentFromFunctionResponse(fc.Name, resultMap, "function")
				contents = append(contents, functionContent)
			}
		}

		// Store last text response if any
		lastResponse = res.Text()

		// Continue to next iteration to get Gemini's response based on function results
	}

	// Max iterations reached
	return &QueryGenerationResponse{
		Message:      "I apologize, but I needed too many steps to generate your query. Please try a simpler request.",
		GeneratedSQL: "",
	}, nil
}

// extractSQLFromResponse extracts SQL code from markdown code blocks
func extractSQLFromResponse(response string) string {
	// Regular expression to find markdown code blocks.
	// (?i) - case insensitive (for SQL, sql, Sql)
	// (?s) - dot matches newlines
	// ```(?:sql)? - matches ``` optionally followed by sql
	// \s* - matches any leading whitespace/newlines
	// (.*?) - non-greedy match of the content
	// \s*``` - matches any trailing whitespace followed by ```
	re := regexp.MustCompile("(?i)(?s)```(?:sql)?\\s*(.*?)\\s*```")

	matches := re.FindStringSubmatch(response)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}
