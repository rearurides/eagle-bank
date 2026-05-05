package domain

import (
	"time"

	"github.com/rearurides/eagle-bank/internal/domain/validation"
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
	userID,
	name string,
	accountType string,
	currency Currency,
) (*Account, *validation.ValidationError) {
	v := validation.NewValidator().
		Required("name", name, "name is required").
		ValidEnum("accountType", accountType, []string{string(AccountTypePersonal), string(AccountTypeSavings)}, "invalid account type")

	if v.HasErrors() {
		return nil, v.ToError("invalid account")
	}

	return &Account{
		UserID:      userID,
		Name:        name,
		AccountType: AccountType(accountType),
		Currency:    currency,
		Balance:     0,
		MinorUnit:   getMinorUnit(currency),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
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
