package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
)

type transactionsHandler struct {
	service  transactionsService
	validate *validator.Validate
}

func NewTransactionsHandler(
	svc *service.TransactionsService,
	v *validator.Validate,
) *transactionsHandler {
	return &transactionsHandler{service: svc, validate: v}
}

func (h *transactionsHandler) HandleCreateTransactions(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.PathValue("accountNumber")
	if accountNumber == "" {
		writeError(w, http.StatusBadRequest, errorResponse{
			Message: "account number is required",
		})
		return
	}

	var body CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, errorResponse{
			Message: "invalid request body",
		})
		return
	}

	userId, ok := middleware.GetUserID(r)
	if !ok || userId == "" {
		writeError(w, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	if errs := h.validate.Struct(body); errs != nil {
		valErrs, ok := errors.AsType[validator.ValidationErrors](errs)
		if !ok {
			panic(fmt.Sprintf("unexpected validation error type: %T", errs))
		}

		writeError(w, http.StatusBadRequest, errorResponse{
			Message: "invalid account",
			Details: formatValidationErrs(valErrs),
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
		if errors.Is(err, domain.ErrInsufficientFunds) {
			writeError(w, http.StatusBadRequest, errorResponse{
				Message: err.Error(),
			})
			return
		}

		writeError(w, http.StatusInternalServerError, errorResponse{
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
