package auth

import "github.com/yelaco/ai-document-backend/pkg/util"

type Argon2PasswordHasher struct{}

func NewArgon2PasswordHasher() util.PasswordHasher {
	return &Argon2PasswordHasher{}
}

func (apu *Argon2PasswordHasher) HashPassword(password string) (string, error) {
	// Implement Argon2 hashing here
	return "", nil
}

func (apu *Argon2PasswordHasher) VerifyPassword(hashedPassword, password string) error {
	// Implement Argon2 password verification here
	return nil
}
