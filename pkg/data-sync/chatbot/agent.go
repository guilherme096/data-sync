package chatbot

type ChatMessage struct {
	Role    string
	Content string
}

type AgentActions interface {
	SendMessage(message string) (string, error)
	SendMessageWithHistory(message string, history []ChatMessage) (string, error)
}
