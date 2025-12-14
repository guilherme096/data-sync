package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/chatbot"
	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/query"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history,omitempty"`
}

type ChatResponse struct {
	Message     string       `json:"message"`
	ToolResults []ToolResult `json:"toolResults,omitempty"`
}

type ToolResult struct {
	ToolName string      `json:"toolName"`
	Data     interface{} `json:"data"`
}

type QueryGenerationResponse struct {
	Message      string `json:"message"`
	GeneratedSQL string `json:"generatedSQL"`
}

type ChatbotRouter struct {
	agent      chatbot.AgentActions
	translator query.QueryTranslator
	discovery  discovery.MetadataDiscovery
	storage    storage.MetadataStorage
}

func NewChatbotRouter(agent chatbot.AgentActions, translator query.QueryTranslator, discovery discovery.MetadataDiscovery, storage storage.MetadataStorage) *ChatbotRouter {
	return &ChatbotRouter{
		agent:      agent,
		translator: translator,
		discovery:  discovery,
		storage:    storage,
	}
}

func (r *ChatbotRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/chatbot/message", r.handleSendMessage)
	mux.HandleFunc("/chatbot/generate-query", r.handleGenerateQuery)
}

func (r *ChatbotRouter) handleSendMessage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var chatReq ChatRequest
	if err := json.NewDecoder(req.Body).Decode(&chatReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if chatReq.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Convert history to chatbot format
	var history []chatbot.ChatMessage
	for _, msg := range chatReq.History {
		history = append(history, chatbot.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Create tool executor with translator, discovery, and storage
	toolExecutor := chatbot.NewToolExecutor(r.translator, r.discovery, r.storage)

	// Get response from chatbot with tools
	agentResponse, err := r.agent.SendMessageWithTools(chatReq.Message, history, toolExecutor)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert tool results to API format
	var toolResults []ToolResult
	for _, tr := range agentResponse.ToolResults {
		toolResults = append(toolResults, ToolResult{
			ToolName: tr.ToolName,
			Data:     tr.Data,
		})
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Message:     agentResponse.Message,
		ToolResults: toolResults,
	})
}

func (r *ChatbotRouter) handleGenerateQuery(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var chatReq ChatRequest
	if err := json.NewDecoder(req.Body).Decode(&chatReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if chatReq.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Convert history to chatbot format
	var history []chatbot.ChatMessage
	for _, msg := range chatReq.History {
		history = append(history, chatbot.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Create query generator tool executor (no query execution)
	toolExecutor := chatbot.NewQueryGeneratorToolExecutor(r.discovery, r.storage)

	// Get query generation response from chatbot
	queryGenResponse, err := r.agent.SendMessageForQueryGeneration(chatReq.Message, history, toolExecutor)
	if err != nil {
		http.Error(w, "Failed to generate query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(QueryGenerationResponse{
		Message:      queryGenResponse.Message,
		GeneratedSQL: queryGenResponse.GeneratedSQL,
	})
}
