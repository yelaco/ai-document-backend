package interfaces

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type ChatService interface {
	AnswerQuestion(ctx context.Context, chatID types.ChatID, question string) (<-chan dtos.AnswerQuestionResult, error)
	CreateChat(ctx context.Context, params dtos.CreateChatParams) (string, error)
	GetPaginatedChats(ctx context.Context, page int64, pageSize int64) ([]entity.Chat, int64, error)
	UpdateChat(ctx context.Context, chatID types.ChatID, params dtos.UpdateChatParams) error
	GetChatByID(ctx context.Context, chatID types.ChatID) (entity.Chat, error)
	DeleteChat(ctx context.Context, chatID types.ChatID) error
	DeleteChatsByDocumentID(ctx context.Context, documentID types.DocumentID) error
	DeleteChatsByUserID(ctx context.Context, userID types.UserID) error
}

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *entity.Chat) error
	GetChat(ctx context.Context, chatID types.ChatID, userID types.UserID) (entity.Chat, error)
	GetPaginatedChats(ctx context.Context, userID types.UserID, skip int64, limit int64) ([]entity.Chat, int64, error)
	UpdateChat(ctx context.Context, chatID types.ChatID, params dtos.UpdateChatParams) error
	DeleteChat(ctx context.Context, chatID types.ChatID, userID types.UserID) error
	DeleteChatsByDocumentID(ctx context.Context, documentID types.DocumentID) error
	DeleteChatsByUserID(ctx context.Context, userID types.UserID) error
}
