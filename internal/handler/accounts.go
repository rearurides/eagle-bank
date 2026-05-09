package handler

import (
	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/service"
)

type CreateAccountRequest struct {
	Name        string `json:"name" validate:"required"`
	AccountType string `json:"accountType" validate:"required,oneof=personal savings"`
}

type AccountResponse struct {
	AccountNumber string `json:"accountNumber"`
	SortCode      string `json:"sortCode"`
	Name          string `json:"name"`
	AccountType   string `json:"accountType"`
	Balance       Money  `json:"balance"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"createdTimestamp"`
	UpdatedAt     string `json:"updatedTimestamp"`
}

type accountsService interface {
	CreateAccount(input service.CreateAccountInput) (*domain.Account, error)
	GetAccountByNumber(userId, accountNumber string) (*domain.Account, error)
}

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
