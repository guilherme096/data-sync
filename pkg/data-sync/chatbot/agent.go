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

type AgentActions interface {
	SendMessage(message string) (string, error)
	SendMessageWithHistory(message string, history []ChatMessage) (string, error)
	SendMessageWithTools(message string, history []ChatMessage, tools ToolExecutor) (*AgentResponse, error)
}
