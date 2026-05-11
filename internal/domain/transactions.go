package domain

import (
	"time"
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
) *Transaction {
	now := time.Now().UTC()

	return &Transaction{
		Amount:          amount,
		Currency:        Currency(currency),
		TransactionType: TransactionType(transactionType),
		Reference:       reference,
		MinorUnit:       getMinorUnit(Currency(currency)),
		CreatedAt:       now,
	}
}
