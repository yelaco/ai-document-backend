package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/database/sqlc"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/mappers"
)

type PostgresUserRepository struct {
	// Add necessary fields, e.g., database connection
	querier sqlc.Querier
}

func NewUserRepository() interfaces.UserRepository {
	return &PostgresUserRepository{}
}

// FindByID implements interfaces.UserRepository.
func (p *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	user, err := p.querier.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	return mappers.UserEntityFromDBModel(user), nil
}
