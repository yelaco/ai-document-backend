package services

import "github.com/yelaco/ai-document-backend/internal/domain/interfaces"

type AuthService struct {
	userRepo interfaces.UserRepository
}
