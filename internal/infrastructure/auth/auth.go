package auth

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func IsValidRole(role Role) bool {
	switch role {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

type AuthClaims struct {
	Email string
	Role  Role
}
