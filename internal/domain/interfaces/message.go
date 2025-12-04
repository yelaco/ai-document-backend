package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type MessageService interface {
	ProcessMessage(ctx context.Context, input string) (string, error)
	CreateMessage(ctx context.Context, chatID uuid.UUID, role string, content string, timestamp time.Time) error
	GetPaginatedMessages(ctx context.Context, chatID uuid.UUID, page int64, pageSize int64) ([]entity.Message, int64, error)
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (entity.Message, error)
	DeleteMessagesByIDs(ctx context.Context, messageIDs []uuid.UUID) error
	DeleteMessagesByChatID(ctx context.Context, chatID uuid.UUID) error
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.Message) error
	GetPaginatedMessages(ctx context.Context, chatID uuid.UUID, skip int64, limit int64) ([]entity.Message, int64, error)
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (entity.Message, error)
	DeleteMessagesByIDs(ctx context.Context, messageIDs []uuid.UUID) error
	DeleteMessagesByChatID(ctx context.Context, chatID uuid.UUID) error
}
