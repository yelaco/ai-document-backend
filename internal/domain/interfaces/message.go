package interfaces

import (
	"context"
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type MessageService interface {
	ProcessMessage(ctx context.Context, input string) (string, error)
	CreateMessage(ctx context.Context, chatID types.ChatID, role string, content string, timestamp time.Time) error
	GetPaginatedMessages(ctx context.Context, chatID types.ChatID, page int64, pageSize int64) ([]entity.Message, int64, error)
	GetMessageByID(ctx context.Context, messageID types.MessageID) (entity.Message, error)
	DeleteMessagesByIDs(ctx context.Context, messageIDs []types.MessageID) error
	DeleteMessagesByChatID(ctx context.Context, chatID types.ChatID) error
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.Message) error
	GetPaginatedMessages(ctx context.Context, chatID types.ChatID, skip int64, limit int64) ([]entity.Message, int64, error)
	GetMessageByID(ctx context.Context, messageID types.MessageID) (entity.Message, error)
	DeleteMessagesByIDs(ctx context.Context, messageIDs []types.MessageID) error
	DeleteMessagesByChatID(ctx context.Context, chatID types.ChatID) error
}
