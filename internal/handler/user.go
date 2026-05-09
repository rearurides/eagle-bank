package handler

import (
	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/service"
)

type Addr struct {
	Line1    string `json:"line1" validate:"required"`
	Line2    string `json:"line2"`
	Line3    string `json:"line3"`
	Town     string `json:"town" validate:"required"`
	County   string `json:"county" validate:"required"`
	PostCode string `json:"postcode" validate:"required"`
}

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Address     Addr   `json:"address"`
	PhoneNumber string `json:"phoneNumber" validate:"required,e164"`
	Email       string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Addr        Addr   `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	CreatedAt   string `json:"createdTimestamp"`
	UpdatedAt   string `json:"updatedTimestamp"`
}

type userService interface {
	CreateUser(input service.CreateUserInput) (*domain.User, error)
	Login(input service.LoginInput) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
}

// mocks
type mockUserService struct {
	createUser  func(input service.CreateUserInput) (*domain.User, error)
	login       func(input service.LoginInput) (*domain.User, error)
	getUserByID func(id string) (*domain.User, error)
}

func (m *mockUserService) CreateUser(input service.CreateUserInput) (*domain.User, error) {
	return m.createUser(input)
}

func (m *mockUserService) Login(input service.LoginInput) (*domain.User, error) {
	return m.login(input)
}

func (m *mockUserService) GetUserByID(id string) (*domain.User, error) {
	return m.getUserByID(id)
}
