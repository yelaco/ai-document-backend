package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type ChatService interface {
	AnswerQuestion(ctx context.Context, chatID uuid.UUID, question string) (<-chan dtos.AnswerQuestionResult, error)
	CreateChat(ctx context.Context, params dtos.CreateChatParams) (string, error)
	GetPaginatedChats(ctx context.Context, page int64, pageSize int64) ([]entity.Chat, int64, error)
	UpdateChat(ctx context.Context, chatID uuid.UUID, params dtos.UpdateChatParams) error
	GetChatByID(ctx context.Context, chatID uuid.UUID) (entity.Chat, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	DeleteChatsByDocumentID(ctx context.Context, documentID uuid.UUID) error
	DeleteChatsByUserID(ctx context.Context, userID uuid.UUID) error
}

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *entity.Chat) error
	GetChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (entity.Chat, error)
	GetPaginatedChats(ctx context.Context, userID uuid.UUID, skip int64, limit int64) ([]entity.Chat, int64, error)
	UpdateChat(ctx context.Context, chatID uuid.UUID, params dtos.UpdateChatParams) error
	DeleteChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	DeleteChatsByDocumentID(ctx context.Context, documentID uuid.UUID) error
	DeleteChatsByUserID(ctx context.Context, userID uuid.UUID) error
}
