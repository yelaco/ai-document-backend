package auth

import (
	"encoding/base64"

	"github.com/yelaco/ai-document-backend/pkg/util"
)

const RefreshTokenExpirationDays = 7

func createRefreshToken() (string, error) {
	b, err := util.GenerateRandomBytes(64)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
