package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/domain/validation"
	"github.com/rearurides/eagle-bank/internal/service"
)

// mock implementation of the UserService interface
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

func newTestUserHandler(svc userService) *userHandler {
	return &userHandler{service: svc}
}

var stubUser = &domain.User{
	ID:          "usr-abc123",
	Name:        "John Doe",
	Email:       "john@example.com",
	PhoneNumber: "+447911123456",
	Addr: domain.Addr{
		Line1:    "123 Main St",
		Town:     "London",
		County:   "Greater London",
		PostCode: "SW1A 1AA",
	},
	CreatedAt: time.Now().UTC(),
	UpdatedAt: time.Now().UTC(),
}

func TestHandleCreateUser_Success(t *testing.T) {
	svc := &mockUserService{
		createUser: func(input service.CreateUserInput) (*domain.User, error) {
			return stubUser, nil
		},
	}

	body := CreateUserRequest{
		Name:        "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "+447911123456",
		Address: Addr{
			Line1:    "123 Main St",
			Town:     "London",
			County:   "Greater London",
			PostCode: "SW1A 1AA",
		},
	}

	rr := makeRequest(t, svc, body)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	var response UserResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != stubUser.ID {
		t.Errorf("expected ID %s, got %s", stubUser.ID, response.ID)
	}
	if response.Email != stubUser.Email {
		t.Errorf("expected email %s, got %s", stubUser.Email, response.Email)
	}
}

func TestHandleCreateUser_ValidationError(t *testing.T) {
	svc := &mockUserService{
		createUser: func(input service.CreateUserInput) (*domain.User, error) {
			return nil, &validation.ValidationError{
				Message: "invalid user",
				Items:   []validation.ValidationItem{{Field: "email", Message: "email is required"}},
			}
		},
	}

	rr := makeRequest(t, svc, CreateUserRequest{})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleCreateUser_InternalServerError(t *testing.T) {
	svc := &mockUserService{
		createUser: func(input service.CreateUserInput) (*domain.User, error) {
			return nil, errors.New("db connection lost")
		},
	}

	rr := makeRequest(t, svc, CreateUserRequest{
		Name:        "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "+447911123456",
	})

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestHandleCreateUser_InvalidJSON(t *testing.T) {
	svc := &mockUserService{}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	newTestUserHandler(svc).handleCreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func makeRequest(t *testing.T, svc userService, body CreateUserRequest) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	newTestUserHandler(svc).handleCreateUser(rr, req)
	return rr
}
