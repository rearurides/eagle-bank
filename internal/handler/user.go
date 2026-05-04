package handler

import (
	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/service"
)

type Addr struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	Line3    string `json:"line3"`
	Town     string `json:"town"`
	County   string `json:"county"`
	PostCode string `json:"postcode"`
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     Addr   `json:"address"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     Addr   `json:"address"`
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
