package domain

import (
	"time"
)

type AccountType string

const (
	AccountTypePersonal AccountType = "personal"
	AccountTypeSavings  AccountType = "savings"
)

type Currency string

const (
	GBP Currency = "GBP"
	USD Currency = "USD"
	EUR Currency = "EUR"
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

// getMinorUnit returns the minor unit multiplier for a given currency.
// For example, for GBP, USD, and EUR, the minor unit is 100 (i.e., 1 pound/dollar/euro = 100 pence/cents).
func getMinorUnit(currency Currency) int64 {
	switch currency {
	case GBP, USD, EUR:
		return 100
	default:
		return 100
	}
}
