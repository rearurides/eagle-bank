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
		tokenUserID string
		name        string
		input       string
		createFunc  func(input service.CreateAccountInput) (*domain.Account, error)
		wantStatus  int
		wantBody    *AccountResponse
		wantErr     *errorResponse
	}{
		{
			name:        "successful account creation",
			tokenUserID: "usr-abc123",
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
			name:        "validation error",
			tokenUserID: "usr-abc123",
			input:       `{}`,
			wantStatus:  http.StatusBadRequest,
			wantErr: &errorResponse{
				Message: "invalid account",
				Details: []validationItem{
					{Field: "name", Message: "must not be blank", Type: "required"},
					{Field: "accountType", Message: "must not be blank", Type: "required"},
				},
			},
		},
		{
			name:        "validation error",
			tokenUserID: "usr-abc123",
			input: `{
				"name": "test",
				"accountType": "host"
			}`,
			wantStatus: http.StatusBadRequest,
			wantErr: &errorResponse{
				Message: "invalid account",
				Details: []validationItem{
					{Field: "accountType", Message: "must be one of: personal, savings", Type: "invalid_enum"},
				},
			},
		},
		{
			name:        "unathorized error",
			tokenUserID: "",
			input: `{
				"name": "test",
				"accountType": "personal"
			}`,
			wantStatus: http.StatusUnauthorized,
			wantErr: &errorResponse{
				Message: ErrUnauthorized.Error(),
			},
		},
		{
			name:        "invalid JSON",
			tokenUserID: "usr-abc123",
			input:       "{invalid json}",
			wantStatus:  http.StatusBadRequest,
			wantErr: &errorResponse{
				Message: ErrInvalidRequestBody.Error(),
				Details: nil,
			},
		},
		{
			name:        "service error",
			tokenUserID: "usr-abc123",
			input: `{
				"name": "My Savings Account",
				"accountType": "savings"
			}`,
			createFunc: func(input service.CreateAccountInput) (*domain.Account, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &errorResponse{
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
				handler.HandleCreateAccount(w, r)
			})

			req := newAuthRequest(t, http.MethodPost, "/v1/accounts", tc.input, tc.tokenUserID)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assertResponse(t, w, tc.wantStatus, tc.wantBody, tc.wantErr)

		})
	}
}

func TestAccountsHandler_handleGetAccountByNumber(t *testing.T) {
	testCases := []struct {
		name          string
		userID        string
		tokenUserID   string
		accountNumber string
		getFunc       func(userId, accountNumber string) (*domain.Account, error)
		wantStatus    int
		wantBody      *AccountResponse
		wantErr       *errorResponse
	}{
		{
			name:          "account found",
			userID:        "usr-abc123",
			tokenUserID:   "usr-abc123",
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
			tokenUserID:   "usr-abc123",
			accountNumber: "99999999",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, domain.ErrAccountNotFound
			},
			wantStatus: http.StatusNotFound,
			wantErr: &errorResponse{
				Message: "account not found",
				Details: nil,
			},
		},
		{
			name:          "unauthorized - no token",
			userID:        "",
			accountNumber: "12345678",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, nil
			},
			wantStatus: http.StatusUnauthorized,
			wantErr: &errorResponse{
				Message: "unauthorized",
				Details: nil,
			},
		},
		{
			name:          "service error",
			userID:        "usr-abc123",
			tokenUserID:   "usr-abc123",
			accountNumber: "12345678",
			getFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &errorResponse{
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

			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/accounts/{accountNumber}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAccountByNumber(w, r)
			})

			req := newAuthRequest(t, http.MethodGet, "/v1/accounts/"+tc.accountNumber, "", tc.tokenUserID)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assertResponse(t, w, tc.wantStatus, tc.wantBody, tc.wantErr)
		})
	}
}
