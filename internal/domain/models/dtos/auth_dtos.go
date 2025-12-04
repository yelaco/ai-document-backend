package dtos

type RegisterUserParams struct {
	Email    string
	FullName string
	Password string
}

type LoginUserParams struct {
	Email    string
	Password string
}

type RefreshFlowParams struct {
	OldRefreshToken string
}
