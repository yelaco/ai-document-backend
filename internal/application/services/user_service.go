package services

import "github.com/yelaco/ai-document-backend/internal/domain/interfaces"

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUser(id int) string {
	return "User"
}
