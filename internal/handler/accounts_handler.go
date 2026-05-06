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

type accountsHandler struct {
	service accountsService
	tm      *token.Manager
}

func newAccountsHandler(svc *service.AccountsService, tokenManager *token.Manager) *accountsHandler {
	return &accountsHandler{service: svc, tm: tokenManager}
}

func (h *accountsHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var body CreateAccountRequest
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

	input := service.CreateAccountInput{
		UserID:      userId,
		Name:        body.Name,
		AccountType: body.AccountType,
	}

	account, err := h.service.CreateAccount(input)
	if err != nil {
		if valErr, ok := errors.AsType[*validation.ValidationError](err); ok {
			writeError(w, http.StatusBadRequest, BadRequestMessage{
				Message: valErr.Message,
				Details: valErr.Items,
			})
			return
		}

		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "failed to create account",
		})
		return
	}

	response := &AccountResponse{
		AccountNumber: account.AccountNumber,
		SortCode:      account.SortCode,
		Name:          account.Name,
		AccountType:   string(account.AccountType),
		Balance:       Money(account.Balance),
		Currency:      string(account.Currency),
		CreatedAt:     account.CreatedAt.Format(TimeLayout),
		UpdatedAt:     account.UpdatedAt.Format(TimeLayout),
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *accountsHandler) handleGetAccountByNumber(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.PathValue("accountNumber")
	if accountNumber == "" {
		writeError(w, http.StatusBadRequest, BadRequestMessage{
			Message: "account number is required",
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

	account, err := h.service.GetAccountByNumber(userId, accountNumber)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			writeError(w, http.StatusNotFound, BadRequestMessage{
				Message: "account not found",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "failed to retrieve account",
		})
		return
	}

	response := &AccountResponse{
		AccountNumber: account.AccountNumber,
		SortCode:      account.SortCode,
		Name:          account.Name,
		AccountType:   string(account.AccountType),
		Balance:       Money(account.Balance),
		Currency:      string(account.Currency),
		CreatedAt:     account.CreatedAt.Format(TimeLayout),
		UpdatedAt:     account.UpdatedAt.Format(TimeLayout),
	}

	writeJSON(w, http.StatusOK, response)
}
