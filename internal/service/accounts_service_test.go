package service

import (
	"testing"

	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/domain/validation"
)

type mockAccountRepository struct {
	createFunc             func(account *domain.Account) error
	getByAccountNumberFunc func(userId, accountNumber string) (*domain.Account, error)
}

func (m *mockAccountRepository) Create(account *domain.Account) error {
	if m.createFunc != nil {
		return m.createFunc(account)
	}
	return nil
}

func (m *mockAccountRepository) GetByAccountNumber(userId, accountNumber string) (*domain.Account, error) {
	if m.getByAccountNumberFunc != nil {
		return m.getByAccountNumberFunc(userId, accountNumber)
	}
	return nil, domain.ErrAccountNotFound
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{}
}

func TestAccountsService_CreateAccount(t *testing.T) {
	testCases := []struct {
		name       string
		input      CreateAccountInput
		createFunc func(account *domain.Account) error
		wantErr    bool
	}{
		{
			name: "successful account creation",
			input: CreateAccountInput{
				UserID:      "usr-abc123",
				Name:        "My Savings Account",
				AccountType: "savings",
			},
			createFunc: func(account *domain.Account) error {
				account.AccountNumber = "12345678"
				return nil
			},
			wantErr: false,
		},
		{
			name: "validation error",
			input: CreateAccountInput{
				UserID:      "usr-abc123",
				Name:        "",
				AccountType: "savings",
			},
			createFunc: func(account *domain.Account) error {
				return &validation.ValidationError{
					Message: "invalid account",
					Items: []validation.ValidationItem{
						{Field: "name", Message: "name is required"},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "repository error",
			input: CreateAccountInput{
				UserID:      "usr-abc123",
				Name:        "My Savings Account",
				AccountType: "savings",
			},
			createFunc: func(account *domain.Account) error {
				return domain.ErrAccountNotFound
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := newMockAccountRepository()
			mockRepo.createFunc = tc.createFunc
			service := NewAccountsService(mockRepo)

			account, err := service.CreateAccount(tc.input)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if account != nil && account.AccountNumber == "" {
				t.Errorf("expected account number to be set, got empty string")
			}
		})
	}
}

func TestAccountsService_GetAccountByNumber(t *testing.T) {
	testCases := []struct {
		name                   string
		userId                 string
		accountNumber          string
		getByAccountNumberFunc func(userId, accountNumber string) (*domain.Account, error)
		wantErr                bool
	}{
		{
			name:          "account found",
			userId:        "usr-abc123",
			accountNumber: "12345678",
			getByAccountNumberFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return &domain.Account{
					AccountNumber: accountNumber,
					UserID:        userId,
					Name:          "My Savings Account",
					AccountType:   "savings",
				}, nil
			},
			wantErr: false,
		},
		{
			name:          "account not found",
			userId:        "usr-abc123",
			accountNumber: "99999999",
			getByAccountNumberFunc: func(userId, accountNumber string) (*domain.Account, error) {
				return nil, domain.ErrAccountNotFound
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := newMockAccountRepository()
			mockRepo.getByAccountNumberFunc = tc.getByAccountNumberFunc
			service := NewAccountsService(mockRepo)

			account, err := service.GetAccountByNumber(tc.userId, tc.accountNumber)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if account != nil && account.AccountNumber != tc.accountNumber {
				t.Errorf("expected account number %s, got %s", tc.accountNumber, account.AccountNumber)
			}
		})
	}
}
