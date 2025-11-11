package interfaces

type UserRepository interface {
	GetUserByID(id int64) (User, error)
}

type UserService interface{}
