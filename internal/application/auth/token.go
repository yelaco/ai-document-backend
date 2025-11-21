package auth

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	AccessTokenExpirationMinutes = 60
	RefreshTokenExpirationDays   = 30
	TokenIssuer                  = "ai-document-backend"
)

type Claims struct {
	Sub   string
	Exp   int64
	Iss   string
	Role  string
	Email string
}

func createAccessToken(userID string, email string, role string) (string, error) {
	// exp := time.Now().Add(AccessTokenExpirationMinutes * time.Minute).Unix()
	// claims := Claims{
	// 	Sub:   userID,
	// 	Exp:   exp,
	// 	Iss:   TokenIssuer,
	// 	Email: email,
	// 	Role:  role,
	// }
	return "access_token_placeholder", nil // Replace with actual token generation logic
}

func createRefreshToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
