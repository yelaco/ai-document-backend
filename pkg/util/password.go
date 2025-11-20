package util

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
