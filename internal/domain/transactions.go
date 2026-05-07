package domain

import (
	"time"

	"github.com/rearurides/eagle-bank/internal/domain/validation"
)

type TransactionType string

const (
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdrawal = "withdrawal"
)

type Transaction struct {
	TransactionID   string
	AccountID       int64
	Amount          int64
	TransactionType TransactionType
	Reference       string
	Currency        Currency
	MinorUnit       int64
	CreatedAt       time.Time
}

func NewTransaction(
	amount int64,
	currency string,
	transactionType string,
	reference string,
) (*Transaction, *validation.ValidationError) {
	v := validation.NewValidator().
		Required("amount", amount, "amount is required").
		Required("currency", currency, "currency is required").
		Required("type", transactionType, "type is required")

	if v.HasErrors() {
		return nil, v.ToError("invalid transaction")
	}

	v.ValidEnum("type", transactionType, []string{string(TransactionTypeDeposit), string(TransactionTypeWithdrawal)}, "invalid transaction type")
	// TODO: check for correct currency

	if v.HasErrors() {
		return nil, v.ToError("invalid transaction")
	}

	now := time.Now().UTC()

	return &Transaction{
		Amount:          amount,
		Currency:        Currency(currency),
		TransactionType: TransactionType(transactionType),
		Reference:       reference,
		MinorUnit:       getMinorUnit(Currency(currency)),
		CreatedAt:       now,
	}, nil
}
