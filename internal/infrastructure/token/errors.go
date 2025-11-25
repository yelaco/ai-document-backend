package token

import "fmt"

var (
	ErrInvalidToken            = fmt.Errorf("invalid token")
	ErrUnexpectedSigningMethod = fmt.Errorf("unexpected signing method")
	ErrExpiredToken            = fmt.Errorf("token is expired")
)
