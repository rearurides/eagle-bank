package service

import (
	"github.com/rearurides/eagle-bank/internal/domain"
)

type UserService struct {
	repo domain.IUserRepository
}

// NewUserService creates a new instance of UserService with the given user repository.
func NewUserService(repo domain.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

type CreateUserInput struct {
	Name        string
	Email       string
	PhoneNumber string
	Address     domain.Addr
}

// CreateUser creates a new user based on the provided input. It validates the input and generates a new user ID.
func (s *UserService) CreateUser(req CreateUserInput) (*domain.User, error) {
	user, valErr := domain.NewUser(
		domain.GenerateID("usr"),
		req.Name,
		req.Email,
		req.PhoneNumber,
		req.Address,
	)
	if valErr != nil {
		return nil, valErr
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
