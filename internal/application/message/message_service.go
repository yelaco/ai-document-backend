package message

import (
	"context"
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type MessageService struct{}

func NewMessageService() interfaces.MessageService {
	return &MessageService{}
}

// CreateMessage implements interfaces.MessageService.
func (m *MessageService) CreateMessage(ctx context.Context, chatID types.ChatID, role string, content string, timestamp time.Time) error {
	panic("unimplemented")
}

// DeleteMessagesByChatID implements interfaces.MessageService.
func (m *MessageService) DeleteMessagesByChatID(ctx context.Context, chatID types.ChatID) error {
	panic("unimplemented")
}

// DeleteMessagesByIDs implements interfaces.MessageService.
func (m *MessageService) DeleteMessagesByIDs(ctx context.Context, messageIDs []types.MessageID) error {
	panic("unimplemented")
}

// GetMessageByID implements interfaces.MessageService.
func (m *MessageService) GetMessageByID(ctx context.Context, messageID types.MessageID) (entity.Message, error) {
	panic("unimplemented")
}

// GetPaginatedMessages implements interfaces.MessageService.
func (m *MessageService) GetPaginatedMessages(ctx context.Context, chatID types.ChatID, page int64, pageSize int64) ([]entity.Message, int64, error) {
	panic("unimplemented")
}

// ProcessMessage implements interfaces.MessageService.
func (m *MessageService) ProcessMessage(ctx context.Context, input string) (string, error) {
	panic("unimplemented")
}
