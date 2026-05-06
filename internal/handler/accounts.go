package handler

import (
	"fmt"

	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/service"
)

type CreateAccountRequest struct {
	Name        string `json:"name" validate:"required"`
	AccountType string `json:"accountType"`
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

type Money float32

func (m Money) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 8)
	buf = fmt.Appendf(buf, "%.2f", m)

	return buf, nil
}
