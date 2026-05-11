package service

import (
	"fmt"

	"github.com/rearurides/eagle-bank/internal/domain"
)

const (
	MinBalance int64 = 0
	MaxBalance int64 = 1000000
)

type TransactionsService struct {
	tanRepo domain.ITransactionsRepository
	accRepo domain.IAccountRepository
}

// NewTransactionService creates a new instance of TransactionService with the given transaction repository.
func NewTransactionsService(
	tanRepo domain.ITransactionsRepository,
	accountRepo domain.IAccountRepository,
) *TransactionsService {
	return &TransactionsService{tanRepo: tanRepo, accRepo: accountRepo}
}

type CreateTransactionInput struct {
	AccountNumber   string
	Amount          float64
	Currency        string
	TransactionType string
	Reference       string
	UserID          string
}

func (s *TransactionsService) CreateTransaction(tan *CreateTransactionInput) (*domain.Transaction, error) {
	amount := int64(tan.Amount * 100) // TODO: get Minor units

	transaction := domain.NewTransaction(
		amount,
		tan.Currency,
		tan.TransactionType,
		tan.Reference,
	)

	acc, err := s.accRepo.GetByAccountNumber(tan.UserID, tan.AccountNumber)
	if err != nil {
		return nil, err
	}

	switch transaction.TransactionType {
	case domain.TransactionTypeDeposit:
		if acc.Balance+transaction.Amount > MaxBalance {
			return nil, fmt.Errorf("deposit would exceed maximum balance of  %.2f", float64(MaxBalance/transaction.MinorUnit))
		}
	case domain.TransactionTypeWithdrawal:
		if acc.Balance-transaction.Amount < MinBalance {
			return nil, domain.ErrInsufficientFunds
		}
	}
	transaction.AccountID = acc.ID
	transaction.TransactionID = domain.GenerateID("tan")

	ops := map[domain.TransactionType]func(*domain.Transaction) error{
		domain.TransactionTypeDeposit:    s.tanRepo.Deposit,
		domain.TransactionTypeWithdrawal: s.tanRepo.Withdraw,
	}

	op, ok := ops[transaction.TransactionType]
	if !ok {
		return nil, fmt.Errorf("unsupported transaction type: %s", transaction.TransactionType)
	}

	if err := op(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}
