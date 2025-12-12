package chatbot

type AgentActions interface {
	SendMessage(message string) (string, error)
}
