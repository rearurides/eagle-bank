package service

import (
	"errors"
	"testing"

	"github.com/rearurides/eagle-bank/internal/domain"
)

type mockUserRepository struct {
	createFunc     func(user *domain.User) error
	getByEmailFunc func(email string) (*domain.User, error)
	getByIDFunc    func(id string) (*domain.User, error)
}

func (m *mockUserRepository) Create(user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(user)
	}
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(email)
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByID(id string) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, domain.ErrUserNotFound
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{}
}

func TestUserService_CreateUser(t *testing.T) {
	testCases := []struct {
		name       string
		createUser func(input CreateUserInput) (*domain.User, error)
		wantErr    bool
	}{
		{
			name: "successful user creation",
			createUser: func(input CreateUserInput) (*domain.User, error) {
				return &domain.User{
					ID:          "usr-abc123",
					Name:        input.Name,
					Email:       input.Email,
					PhoneNumber: input.PhoneNumber,
					Addr:        input.Address,
				}, nil
			},
			wantErr: false,
		},
		{
			name: "repository error",
			createUser: func(input CreateUserInput) (*domain.User, error) {
				return nil, errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			mockRepo.createFunc = func(user *domain.User) error {
				if tc.createUser != nil {
					_, err := tc.createUser(CreateUserInput{
						Name:        user.Name,
						Email:       user.Email,
						PhoneNumber: user.PhoneNumber,
						Address:     user.Addr,
					})
					return err
				}
				return nil
			}
			svc := NewUserService(mockRepo)

			input := CreateUserInput{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				PhoneNumber: "+441234567890",
				Address: domain.Addr{
					Line1:    "123 Main St",
					Town:     "Anytown",
					County:   "Anycounty",
					PostCode: "12345",
				},
			}
			_, err := svc.CreateUser(input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
