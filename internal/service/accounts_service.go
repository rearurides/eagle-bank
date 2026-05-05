package service

import (
	"github.com/rearurides/eagle-bank/internal/domain"
)

type AccountsService struct {
	repo domain.IAccountRepository
}

// NewAccountsService creates a new instance of AccountsService with the given account repository.
func NewAccountsService(repo domain.IAccountRepository) *AccountsService {
	return &AccountsService{repo: repo}
}

type CreateAccountInput struct {
	UserID      string
	Name        string
	AccountType string
}

// CreateAccount creates a new account based on the provided input. It generates a new account number and sort code.
func (s *AccountsService) CreateAccount(input CreateAccountInput) (*domain.Account, error) {
	// TODO: change to enum
	sorteCode := "10-10-10"

	account, valErr := domain.NewAccount(
		input.UserID, input.Name, input.AccountType, domain.GBP,
	)
	if valErr != nil {
		return nil, valErr
	}

	account.SortCode = sorteCode

	if err := s.repo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}
