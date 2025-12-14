package chatbot

type ChatMessage struct {
	Role    string
	Content string
}

// AgentResponse contains the agent's text response and any tool execution results
type AgentResponse struct {
	Message     string
	ToolResults []ToolResult
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolName string
	Data     interface{}
}

// QueryGenerationResponse contains a generated SQL query and explanation
type QueryGenerationResponse struct {
	Message      string
	GeneratedSQL string
}

type AgentActions interface {
	SendMessage(message string) (string, error)
	SendMessageWithHistory(message string, history []ChatMessage) (string, error)
	SendMessageWithTools(message string, history []ChatMessage, tools ToolExecutor) (*AgentResponse, error)
	SendMessageForQueryGeneration(message string, history []ChatMessage, tools ToolExecutor) (*QueryGenerationResponse, error)
}
