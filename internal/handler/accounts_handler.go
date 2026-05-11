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

type accountsHandler struct {
	svc      accountsService
	validate *validator.Validate
}

func NewAccountsHandler(
	svc *service.AccountsService,
	v *validator.Validate,
) *accountsHandler {
	return &accountsHandler{svc: svc, validate: v}
}

func (h *accountsHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var body CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, errorResponse{
			Message: ErrInvalidRequestBody.Error(),
		})
		return
	}

	tokenUserID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errorResponse{
			Message: ErrUnauthorized.Error(),
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

	input := service.CreateAccountInput{
		UserID:      tokenUserID,
		Name:        body.Name,
		AccountType: body.AccountType,
	}

	account, err := h.svc.CreateAccount(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errorResponse{
			Message: ErrInternalSever.Error(),
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

func (h *accountsHandler) HandleGetAccountByNumber(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.PathValue("accountNumber")

	tokenUserID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errorResponse{
			Message: ErrUnauthorized.Error(),
		})
		return
	}

	account, err := h.svc.GetAccountByNumber(tokenUserID, accountNumber)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			writeError(w, http.StatusNotFound, errorResponse{
				Message: "account not found",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, errorResponse{
			Message: err.Error(),
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
