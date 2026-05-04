package domain

type IUserRepository interface {
	Create(user *User) error
}
