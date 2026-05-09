package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
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
		wantErr    *badRequestErrorResponse
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
			wantErr: &badRequestErrorResponse{
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
			wantErr: &badRequestErrorResponse{
				Message: ErrInvalidRequestBody.Error(),
				Details: nil,
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
			wantErr: &badRequestErrorResponse{
				Message: ErrInternalSever.Error(),
				Details: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockUserService{
				createUser: tc.createFunc,
			}
			handler := newTestUserHandler(mockSvc)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleCreateUser(w, r)
			})

			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(tc.input))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}

			if tc.wantBody != nil {
				var responseBody UserResponse
				if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				if diff := cmp.Diff(*tc.wantBody, responseBody); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			} else if tc.wantErr != nil {
				var responseBody badRequestErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
					t.Fatalf("failed to unmarshal error body: %v", err)
				}
				if diff := cmp.Diff(*tc.wantErr, responseBody); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
