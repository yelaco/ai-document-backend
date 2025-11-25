// Package token provides functionality for managing and verifying tokens, including JWT and PASETO tokens.
package token

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(sub, email, role string) (string, error)

	// VeriyToken checks if the token is valid or not
	VerifyToken(token string, skipExp bool) (*Claims, error)

	GetPublicKey() string
}
