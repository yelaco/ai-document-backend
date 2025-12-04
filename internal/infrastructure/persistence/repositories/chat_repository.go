package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/database/sqlc"
)

type ChatRepository struct {
	connPool *pgxpool.Pool
	queries  *sqlc.Queries
}

func NewChatRepository(connPool *pgxpool.Pool) interfaces.ChatRepository {
	return &ChatRepository{
		connPool: connPool,
		queries:  sqlc.New(connPool),
	}
}

// CreateChat implements interfaces.ChatRepository.
func (c *ChatRepository) CreateChat(ctx context.Context, chat *entity.Chat) error {
	panic("unimplemented")
}

// DeleteChat implements interfaces.ChatRepository.
func (c *ChatRepository) DeleteChat(ctx context.Context, chatID types.ChatID, userID types.UserID) error {
	panic("unimplemented")
}

// DeleteChatsByDocumentID implements interfaces.ChatRepository.
func (c *ChatRepository) DeleteChatsByDocumentID(ctx context.Context, documentID types.DocumentID) error {
	panic("unimplemented")
}

// DeleteChatsByUserID implements interfaces.ChatRepository.
func (c *ChatRepository) DeleteChatsByUserID(ctx context.Context, userID types.UserID) error {
	panic("unimplemented")
}

// GetChat implements interfaces.ChatRepository.
func (c *ChatRepository) GetChat(ctx context.Context, chatID types.ChatID, userID types.UserID) (entity.Chat, error) {
	panic("unimplemented")
}

// GetPaginatedChats implements interfaces.ChatRepository.
func (c *ChatRepository) GetPaginatedChats(ctx context.Context, userID types.UserID, skip int64, limit int64) ([]entity.Chat, int64, error) {
	panic("unimplemented")
}

// UpdateChat implements interfaces.ChatRepository.
func (c *ChatRepository) UpdateChat(ctx context.Context, chatID types.ChatID, params dtos.UpdateChatParams) error {
	panic("unimplemented")
}
