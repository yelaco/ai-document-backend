package handlers

import (
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
)

type ChatHandler struct {
	chatService interfaces.ChatService
}

func NewChatHandler(chatService interfaces.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}
