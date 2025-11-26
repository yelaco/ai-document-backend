package auth

import (
	"errors"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var ErrInvalidRole = errors.New("invalid role")

func ParseRole(role string) (Role, error) {
	switch role {
	case "user":
		return RoleUser, nil
	case "admin":
		return RoleAdmin, nil
	default:
		return "", ErrInvalidRole
	}
}

type AuthClaims struct {
	Email string
	Role  Role
}
