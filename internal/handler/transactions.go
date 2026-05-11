package handler

import (
	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/service"
)

type CreateTransactionRequest struct {
	Amount          float64 `json:"amount" validate:"required,gt=0,decimal2"`
	Currency        string  `json:"currency" validate:"required,oneof=GBP USD"`
	TransactionType string  `json:"type" validate:"required,oneof=deposit withdrawal"`
	Reference       string  `json:"reference"`
}

type TransactionResponse struct {
	ID              string `json:"id"`
	Amount          Money  `json:"amount"`
	Currency        string `json:"currency"`
	TransactionType string `json:"type"`
	Reference       string `json:"reference"`
	UserID          string `json:"userId"`
	CreatedAt       string `json:"createdTimestamp"`
}

type transactionsService interface {
	CreateTransaction(tan *service.CreateTransactionInput) (*domain.Transaction, error)
}
