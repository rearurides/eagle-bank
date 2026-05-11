package domain

import (
	"time"
)

type AccountType string

const (
	AccountTypePersonal AccountType = "personal"
	AccountTypeSavings  AccountType = "savings"
)

type Account struct {
	ID            int64
	AccountNumber string
	SortCode      string
	UserID        string
	Name          string
	AccountType   AccountType
	Balance       int64
	Currency      Currency
	MinorUnit     int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewAccount(
	userId,
	name string,
	accountType string,
	currency Currency,
) *Account {
	now := time.Now().UTC()

	return &Account{
		UserID:      userId,
		Name:        name,
		AccountType: AccountType(accountType),
		Currency:    currency,
		Balance:     0,
		MinorUnit:   getMinorUnit(currency),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}