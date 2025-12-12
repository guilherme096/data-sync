package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/chatbot"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

type ChatbotRouter struct {
	agent chatbot.AgentActions
}

func NewChatbotRouter(agent chatbot.AgentActions) *ChatbotRouter {
	return &ChatbotRouter{agent: agent}
}

func (r *ChatbotRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/chatbot/message", r.handleSendMessage)
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

	// Get response from chatbot
	response, err := r.agent.SendMessage(chatReq.Message)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{Message: response})
}
