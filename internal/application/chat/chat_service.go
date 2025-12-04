package chat

import (
	"context"
	"fmt"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/aigateway/prompts"
	reqContext "github.com/yelaco/ai-document-backend/internal/infrastructure/context"
)

type ChatService struct {
	chatRepo       interfaces.ChatRepository
	ragRetriever   interfaces.RagRetriever
	messageService interfaces.MessageService
	aiGateway      interfaces.AIGateway
}

func NewChatService(chatRepo interfaces.ChatRepository, ragRetriever interfaces.RagRetriever, messageService interfaces.MessageService) interfaces.ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		ragRetriever:   ragRetriever,
		messageService: messageService,
	}
}

// AnswerQuestion implements interfaces.ChatService.
func (c *ChatService) AnswerQuestion(ctx context.Context, chatID types.ChatID, question string) (<-chan dtos.AnswerQuestionResult, error) {
	userID := reqContext.UserIDMustFromContext(ctx)
	chat, err := c.chatRepo.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("ChatService.AnswerQuestion: failed to get chat: %w", err)
	}

	relevantContent, err := c.ragRetriever.RetrieveDocumentContent(ctx, chat.DocumentID, question, 10)
	if err != nil {
		return nil, fmt.Errorf("ChatService.AnswerQuestion: failed to retrieve document content: %w", err)
	}

	// TODO: retrieve chat message for furthur context
	prompt := prompts.BuildAskDocumentPrompt(relevantContent, question, nil, nil)
	result, err := c.aiGateway.GenerateTextStream(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("ChatService.AnswerQuestion: failed to generate text: %w", err)
	}

	resultCh := make(chan dtos.AnswerQuestionResult)
	resultCh <- dtos.AnswerQuestionResult{
		Status: "completed",
		Text:   result,
	}
	defer close(resultCh)

	return resultCh, nil
}

// CreateChat implements interfaces.ChatService.
func (c *ChatService) CreateChat(ctx context.Context, params dtos.CreateChatParams) (string, error) {
	panic("unimplemented")
}

// DeleteChat implements interfaces.ChatService.
func (c *ChatService) DeleteChat(ctx context.Context, chatID types.ChatID) error {
	panic("unimplemented")
}

// DeleteChatsByDocumentID implements interfaces.ChatService.
func (c *ChatService) DeleteChatsByDocumentID(ctx context.Context, documentID types.DocumentID) error {
	panic("unimplemented")
}

// DeleteChatsByUserID implements interfaces.ChatService.
func (c *ChatService) DeleteChatsByUserID(ctx context.Context, userID types.UserID) error {
	panic("unimplemented")
}

// GetChatByID implements interfaces.ChatService.
func (c *ChatService) GetChatByID(ctx context.Context, chatID types.ChatID) (entity.Chat, error) {
	panic("unimplemented")
}

// GetPaginatedChats implements interfaces.ChatService.
func (c *ChatService) GetPaginatedChats(ctx context.Context, page int64, pageSize int64) ([]entity.Chat, int64, error) {
	panic("unimplemented")
}

// UpdateChat implements interfaces.ChatService.
func (c *ChatService) UpdateChat(ctx context.Context, chatID types.ChatID, params dtos.UpdateChatParams) error {
	panic("unimplemented")
}
