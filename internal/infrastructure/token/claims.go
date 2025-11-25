package token

import (
	"time"
)

const (
	AccessTokenExpirationMinutes = 60
	RefreshTokenExpirationDays   = 30
	TokenIssuer                  = "ai-document-backend"
	TokenAudience                = "ai-document-backend-users"
)

type Claims struct {
	Sub   string
	Role  string
	Email string
	Exp   time.Time
	Iss   string
}
