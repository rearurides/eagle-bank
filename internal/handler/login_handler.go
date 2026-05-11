package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rearurides/eagle-bank/internal/service"
	"github.com/rearurides/eagle-bank/pkg/token"
)

type loginHandler struct {
	svc      userService
	tm       token.Manager
	validate *validator.Validate
}

func NewLoginHandler(
	svc *service.UserService,
	v *validator.Validate,
	tm token.Manager,
) *loginHandler {
	return &loginHandler{svc: svc, validate: v, tm: tm}
}

type loginResponse struct {
	Token  *string `json:"token"`
	UserID *string `json:"userId"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password"`
}

func (h *loginHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, errorResponse{
			Message: "invalid request body",
		})
		return
	}

	if errs := h.validate.Struct(body); errs != nil {
		valErrs, ok := errors.AsType[validator.ValidationErrors](errs)
		if !ok {
			panic(fmt.Sprintf("unexpected validation error type: %T", errs))
		}

		writeError(w, http.StatusBadRequest, errorResponse{
			Message: "invalid auth details",
			Details: formatValidationErrs(valErrs),
		})
		return
	}

	user, err := h.svc.Login(service.LoginInput{Email: body.Email})
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(w, http.StatusBadRequest, errorResponse{
				Message: ErrInvalidCredentials.Error(),
			})
			return
		}

		writeError(w, http.StatusInternalServerError, errorResponse{
			Message: "internal server error",
		})
		return
	}

	tokenString, err := h.tm.Generate(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errorResponse{
			Message: "internal server error",
		})
		return
	}

	response := loginResponse{
		Token:  &tokenString,
		UserID: &user.ID,
	}

	writeJSON(w, http.StatusOK, response)
}
