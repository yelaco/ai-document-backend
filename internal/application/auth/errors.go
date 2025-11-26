package auth

import "errors"

var (
	AuthErrorUserNotFound       = errors.New("user not found")
	AuthErrorInvalidCredentials = errors.New("invalid credentials")
	AuthErrorUserIDMissing      = errors.New("missing user ID in context")
)
