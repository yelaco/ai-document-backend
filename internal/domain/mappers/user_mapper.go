package mappers

import (
	entity "github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	database "github.com/yelaco/ai-document-backend/internal/infrastructure/database/sqlc"
)

func UserEntityFromDBModel(user database.User) entity.User {
	return entity.User{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		PasswordHash: user.PasswordHash,
		Role:         auth.RoleUser,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
