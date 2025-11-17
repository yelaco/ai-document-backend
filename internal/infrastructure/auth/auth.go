package auth

type Role int

const (
	RoleUser Role = iota
	RoleAdmin
)
