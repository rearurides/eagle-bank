package handler

import (
	"context"
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
	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
)

func newTestAccountsHandler(svc accountsService) *accountsHandler {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &accountsHandler{svc: svc, validate: validate}
}

func TestAccountsHandler_handleCreateAccount(t *testing.T) {
	testCases := []struct {
		userID     string
		name       string
		input      string
		createFunc func(input service.CreateAccountInput) (*domain.Account, error)
		wantStatus int
		wantBody   *AccountResponse
		wantErr    *badRequestErrorResponse
	}{
		{
			name:   "successful account creation",
			userID: "usr-abc123",
			input: `{
				"name": "My Savings Account",
				"accountType": "savings"
			}`,
			createFunc: func(input service.CreateAccountInput) (*domain.Account, error) {
				return &domain.Account{
					UserID:        input.UserID,
					Name:          input.Name,
					AccountType:   domain.AccountType(input.AccountType),
					Currency:      domain.GBP,
					AccountNumber: "12345678",
					SortCode:      "10-10-10",
				}, nil
			},
			wantStatus: http.StatusCreated,
			wantBody: &AccountResponse{
				AccountNumber: "12345678",
				SortCode:      "10-10-10",
				Name:          "My Savings Account",
				AccountType:   "savings",
				Balance:       0.00,
				Currency:      "GBP",
				CreatedAt:     "0001-01-01T00:00:00.000Z",
				UpdatedAt:     "0001-01-01T00:00:00.000Z",
			},
		},
		{
			name:       "validation error",
			userID:     "usr-abc123",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
			wantErr: &badRequestErrorResponse{
				Message: "invalid account",
				Details: []validationItem{
					{Field: "name", Message: "must not be blank", Type: "required"},
					{Field: "accountType", Message: "must not be blank", Type: "required"},
				},
			},
		},
		{
			name:   "validation error",
			userID: "usr-abc123",
			input: `{
				"name": "test",
				"accountType": "host"
			}`,
			wantStatus: http.StatusBadRequest,
			wantErr: &badRequestErrorResponse{
				Message: "invalid account",
				Details: []validationItem{
					{Field: "accountType", Message: "must be one of: personal, savings", Type: "invalid_enum"},
				},
			},
		},
		{
			name:   "unathorized error",
			userID: "",
			input: `{
				"name": "test",
				"accountType": "personal"
			}`,
			wantStatus: http.StatusBadRequest,
			wantErr: &badRequestErrorResponse{
				Message: ErrUnauthorized.Error(),
			},
		},
		{
			userID:     "usr-abc123",
			name:       "invalid JSON",
			input:      "{invalid json}",
			wantStatus: http.StatusBadRequest,
			wantErr: &badRequestErrorResponse{
				Message: ErrInvalidRequestBody.Error(),
				Details: nil,
			},
		},
		{
			name:   "service error",
			userID: "usr-abc123",
			input: `{
				"name": "My Savings Account",
				"accountType": "savings"
			}`,
			createFunc: func(input service.CreateAccountInput) (*domain.Account, error) {
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
			mockSvc := &mockAccountsService{
				createAccount: tc.createFunc,
			}
			handler := newTestAccountsHandler(mockSvc)

			// Create a mux to handle path parameters
			mux := http.NewServeMux()
			mux.HandleFunc("POST /v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				// Add user ID to context (mimicking middleware)
				ctx := context.WithValue(r.Context(), middleware.UserIDKey, tc.userID)
				r = r.WithContext(ctx)
				handler.HandleCreateAccount(w, r)
			})

			req := httptest.NewRequest(http.MethodPost, "/v1/accounts", strings.NewReader(tc.input))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}

			if tc.wantBody != nil {
				var responseBody AccountResponse
				if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				if diff := cmp.Diff(*tc.wantBody, responseBody); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}

			if tc.wantErr != nil {
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

func TestAccountsHandler_handleGetAccountByNumber(t *testing.T) {
	testCases := []struct {
		name          string
		userID        string
		accountNumber string
		getFunc       func(userId, accountNumber string) (*domain.Account, error)
		wantStatus    int
		wantBody      *AccountResponse
		wantErr       *badRequestErrorResponse
	}{
		{
			name:          "account found",
			userID:        "usr-abc123",
			accountNumber: "12345678",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return &domain.Account{
					UserID:        userId,
					AccountNumber: accountNumber,
					SortCode:      "10-10-10",
					Name:          "My Savings Account",
					AccountType:   "savings",
					Currency:      domain.GBP,
				}, nil
			},
			wantStatus: http.StatusOK,
			wantBody: &AccountResponse{
				AccountNumber: "12345678",
				SortCode:      "10-10-10",
				Name:          "My Savings Account",
				AccountType:   "savings",
				Balance:       0.00,
				Currency:      "GBP",
				CreatedAt:     "0001-01-01T00:00:00.000Z",
				UpdatedAt:     "0001-01-01T00:00:00.000Z",
			},
		},
		{
			name:          "account not found",
			userID:        "usr-abc123",
			accountNumber: "99999999",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, domain.ErrAccountNotFound
			},
			wantStatus: http.StatusNotFound,
			wantErr: &badRequestErrorResponse{
				Message: "account not found",
				Details: nil,
			},
		},
		{
			name:          "invalid user id",
			userID:        "",
			accountNumber: "12345678",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, nil
			},
			wantStatus: http.StatusUnauthorized,
			wantErr: &badRequestErrorResponse{
				Message: "invalid user id",
				Details: nil,
			},
		},
		{
			name:          "service error",
			userID:        "usr-abc123",
			accountNumber: "12345678",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &badRequestErrorResponse{
				Message: "unexpected error",
				Details: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAccountsService{
				getAccountByNumber: tc.getFunc,
			}
			handler := newTestAccountsHandler(mockSvc)

			// Create a mux to handle path parameters
			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/accounts/{accountNumber}", func(w http.ResponseWriter, r *http.Request) {
				// Add user ID to context (mimicking middleware)
				ctx := context.WithValue(r.Context(), middleware.UserIDKey, tc.userID)
				r = r.WithContext(ctx)
				handler.HandleGetAccountByNumber(w, r)
			})

			req := httptest.NewRequest(http.MethodGet, "/v1/accounts/"+tc.accountNumber, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantBody != nil {
				var responseBody AccountResponse
				if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				if !reflect.DeepEqual(&responseBody, tc.wantBody) {
					t.Errorf("expected body %v, got %v", tc.wantBody, &responseBody)
				}
			}
		})
	}
}
