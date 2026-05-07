package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/domain/validation"
	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
	"github.com/rearurides/eagle-bank/pkg/token"
)

type transactionsHandler struct {
	service transactionsService
	tm      *token.Manager
}

func newTransactionsHandler(svc *service.TransactionsService, tokenManager *token.Manager) *transactionsHandler {
	return &transactionsHandler{service: svc, tm: tokenManager}
}

func (h *transactionsHandler) handleCreateTransactions(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.PathValue("accountNumber")
	if accountNumber == "" {
		writeError(w, http.StatusBadRequest, BadRequestMessage{
			Message: "account number is required",
		})
		return
	}

	var body CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, BadRequestMessage{
			Message: "invalid request body",
		})
		return
	}

	userId, ok := middleware.GetUserID(r)
	if !ok || userId == "" {
		writeError(w, http.StatusUnauthorized, BadRequestMessage{
			Message: "unauthorized",
		})
		return
	}

	input := service.CreateTransactionInput{
		Amount:          body.Amount,
		AccountNumber:   accountNumber,
		Currency:        body.Currency,
		TransactionType: body.TransactionType,
		Reference:       body.Reference,
		UserID:          userId,
	}

	transaction, err := h.service.CreateTransaction(&input)
	if err != nil {
		if valErr, ok := errors.AsType[*validation.ValidationError](err); ok {
			writeError(w, http.StatusBadRequest, BadRequestMessage{
				Message: valErr.Message,
				Details: valErr.Items,
			})
			return
		}

		if errors.Is(err, domain.ErrInsufficientFunds) {
			writeError(w, http.StatusBadRequest, BadRequestMessage{
				Message: err.Error(),
			})
			return
		}

		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "failed to process transaction",
		})
		return
	}

	response := &TransactionResponse{
		ID:              transaction.TransactionID,
		Amount:          Money(float64(transaction.Amount) / float64(transaction.MinorUnit)),
		Currency:        string(transaction.Currency),
		TransactionType: string(transaction.TransactionType),
		Reference:       transaction.Reference,
		UserID:          userId,
		CreatedAt:       transaction.CreatedAt.Format(TimeLayout),
	}

	writeJSON(w, http.StatusCreated, response)

}
