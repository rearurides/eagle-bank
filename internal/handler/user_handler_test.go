package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/service"
)

func newTestUserHandler(svc userService) *userHandler {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &userHandler{svc: svc, validate: validate}
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
}

func Test_HandleCreateUser(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		createFunc func(input service.CreateUserInput) (*domain.User, error)
		wantStatus int
		wantBody   *UserResponse
		wantErr    *errorResponse
	}{
		{
			name: "successful user creation",
			input: `{
				"name": "John Doe",
				"address": {
					"line1": "123 Main St",
					"line3": "",
					"town": "London",
					"county": "Greater London",
					"postcode": "SW1A 1AA"
				},
				"phoneNumber": "+447911123456",
				"email": "user@example.com"
			}`,
			createFunc: func(input service.CreateUserInput) (*domain.User, error) {
				return stubUser, nil
			},
			wantStatus: http.StatusCreated,
			wantBody: &UserResponse{
				ID:   "usr-abc123",
				Name: "John Doe",
				Addr: Addr{
					Line1:    "123 Main St",
					Town:     "London",
					County:   "Greater London",
					PostCode: "SW1A 1AA",
				},
				PhoneNumber: "+447911123456",
				Email:       "john@example.com",
				CreatedAt:   "0001-01-01T00:00:00.000Z",
				UpdatedAt:   "0001-01-01T00:00:00.000Z",
			},
		},
		{
			name:       "required fields validation error",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
			wantErr: &errorResponse{
				Message: "invalid user",
				Details: []validationItem{
					{Field: "name", Message: "must not be blank", Type: "required"},
					{Field: "address.line1", Message: "must not be blank", Type: "required"},
					{Field: "address.town", Message: "must not be blank", Type: "required"},
					{Field: "address.county", Message: "must not be blank", Type: "required"},
					{Field: "address.postcode", Message: "must not be blank", Type: "required"},
					{Field: "phoneNumber", Message: "must not be blank", Type: "required"},
					{Field: "email", Message: "must not be blank", Type: "required"},
				},
			},
		},
		{
			name:       "invalid JSON",
			input:      "{invalid json}",
			wantStatus: http.StatusBadRequest,
			wantErr: &errorResponse{
				Message: ErrInvalidRequestBody.Error(),
			},
		},
		{
			name: "service error",
			input: `{
				"name": "Test User",
				"address": {
					"line1": "1 House",
					"line2": "Street st",
					"line3": "",
					"town": "Town",
					"county": "Lancs",
					"postcode": "M1 1AA"
				},
				"phoneNumber": "+441234567890",
				"email": "user@example.com"
			}`,
			createFunc: func(input service.CreateUserInput) (*domain.User, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &errorResponse{
				Message: ErrInternalSever.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockUserService{
				createUser: tc.createFunc,
			}
			handler := newTestUserHandler(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(tc.input))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleCreateUser(w, req)

			assertResponse(t, w, tc.wantStatus, tc.wantBody, tc.wantErr)
		})
	}
}

func Test_HandleGetUserByID(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		tokenUserID string
		getByIdFunc func(input string) (*domain.User, error)
		wantStatus  int
		wantBody    *UserResponse
		wantErr     *errorResponse
	}{
		{
			name:        "successful get user by id",
			userID:      "usr-123abc",
			tokenUserID: "usr-123abc",
			getByIdFunc: func(input string) (*domain.User, error) {
				return stubUser, nil
			},
			wantStatus: http.StatusOK,
			wantBody: &UserResponse{
				ID:   "usr-abc123",
				Name: "John Doe",
				Addr: Addr{
					Line1:    "123 Main St",
					Town:     "London",
					County:   "Greater London",
					PostCode: "SW1A 1AA",
				},
				PhoneNumber: "+447911123456",
				Email:       "john@example.com",
				CreatedAt:   "0001-01-01T00:00:00.000Z",
				UpdatedAt:   "0001-01-01T00:00:00.000Z",
			},
		},
		{
			name:        "forbidden",
			userID:      "usr-123abc",
			tokenUserID: "usr-abc123",
			getByIdFunc: func(input string) (*domain.User, error) {
				return nil, nil
			},
			wantStatus: http.StatusForbidden,
			wantErr: &errorResponse{
				Message: ErrForbidden.Error(),
			},
		},
		{
			name:        "unauthorized - no token",
			userID:      "usr-123abc",
			tokenUserID: "",
			getByIdFunc: func(input string) (*domain.User, error) {
				return nil, nil
			},
			wantStatus: http.StatusUnauthorized,
			wantErr: &errorResponse{
				Message: ErrUnauthorized.Error(),
			},
		},
		{
			name:        "user not found",
			userID:      "usr-abc123",
			tokenUserID: "usr-abc123",
			getByIdFunc: func(input string) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			wantStatus: http.StatusNotFound,
			wantErr: &errorResponse{
				Message: "user not found",
			},
		},
		{
			name:        "service error",
			userID:      "usr-123abc",
			tokenUserID: "usr-123abc",
			getByIdFunc: func(input string) (*domain.User, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &errorResponse{
				Message: ErrInternalSever.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockUserService{
				getUserByID: tc.getByIdFunc,
			}
			handler := newTestUserHandler(mockSvc)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUserByID(w, r)
			})

			req := newAuthRequest(t, http.MethodGet, "/v1/users/"+tc.userID, "", tc.tokenUserID)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assertResponse(t, w, tc.wantStatus, tc.wantBody, tc.wantErr)
		})
	}
}
