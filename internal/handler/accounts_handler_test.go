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

	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/domain/validation"
	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
)

type mockAccountsService struct {
	createAccount      func(input service.CreateAccountInput) (*domain.Account, error)
	getAccountByNumber func(userId, accountNumber string) (*domain.Account, error)
}

func (m *mockAccountsService) CreateAccount(input service.CreateAccountInput) (*domain.Account, error) {
	return m.createAccount(input)
}

func (m *mockAccountsService) GetAccountByNumber(userId, accountNumber string) (*domain.Account, error) {
	return m.getAccountByNumber(userId, accountNumber)
}

func newTestAccountsHandler(svc accountsService) *accountsHandler {
	return &accountsHandler{service: svc}
}

func TestAccountsHandler_handleCreateAccount(t *testing.T) {
	testCases := []struct {
		userID     string
		name       string
		input      string
		createFunc func(input service.CreateAccountInput) (*domain.Account, error)
		wantStatus int
		wantBody   *AccountResponse
		wantErr    *BadRequestMessage
	}{
		{
			userID: "usr-abc123",
			name:   "successful account creation",
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
			userID: "usr-abc123",
			name:   "validation error",
			input: `{
				"name": "",
				"accountType": "savings"
			}`,
			createFunc: func(input service.CreateAccountInput) (*domain.Account, error) {
				return nil, &validation.ValidationError{
					Message: "invalid account",
					Items: []validation.ValidationItem{
						{Field: "name", Message: "name is required"},
					},
				}
			},
			wantStatus: http.StatusBadRequest,
			wantErr: &BadRequestMessage{
				Message: "invalid account",
				Details: []validation.ValidationItem{
					{Field: "name", Message: "name is required"},
				},
			},
		},
		{
			userID: "usr-abc123",
			name:   "invalid JSON",
			input:  "{invalid json}",
			createFunc: func(input service.CreateAccountInput) (*domain.Account, error) {
				return nil, &validation.ValidationError{
					Message: "invalid account",
					Items: []validation.ValidationItem{
						{Field: "name", Message: "name is required"},
					},
				}
			},
			wantStatus: http.StatusBadRequest,
			wantErr: &BadRequestMessage{
				Message: "invalid request body",
				Details: nil,
			},
		},
		{
			userID: "usr-abc123",
			name:   "service error",
			input: `{
				"name": "My Savings Account",
				"accountType": "savings"
			}`,
			createFunc: func(input service.CreateAccountInput) (*domain.Account, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &BadRequestMessage{
				Message: "failed to create account",
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
			mux.HandleFunc("GET /v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				// Add user ID to context (mimicking middleware)
				ctx := context.WithValue(r.Context(), middleware.UserIDKey, tc.userID)
				r = r.WithContext(ctx)
				handler.handleCreateAccount(w, r)
			})

			req := httptest.NewRequest(http.MethodGet, "/v1/accounts", strings.NewReader(tc.input))
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

			if tc.wantErr != nil {
				var responseBody BadRequestMessage
				if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
					t.Fatalf("failed to unmarshal error body: %v", err)
				}
				if !reflect.DeepEqual(&responseBody, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, &responseBody)
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
		wantErr       *BadRequestMessage
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
			wantErr: &BadRequestMessage{
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
			wantErr: &BadRequestMessage{
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
			wantErr: &BadRequestMessage{
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
				handler.handleGetAccountByNumber(w, r)
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
