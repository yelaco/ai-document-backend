package auth

type Claims struct {
	Sub   string
	Exp   int64
	Iss   string
	Role  string
	Email string
}
