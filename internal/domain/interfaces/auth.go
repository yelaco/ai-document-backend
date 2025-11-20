package interfaces

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, email string, fullName string, password string) error
	LoginUser(ctx context.Context, email string, password string) (string, string, error)
}
