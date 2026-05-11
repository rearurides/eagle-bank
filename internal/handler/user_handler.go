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

type userHandler struct {
	svc      userService
	validate *validator.Validate
}

func NewUserHandler(
	svc *service.UserService,
	v *validator.Validate,
) *userHandler {
	return &userHandler{svc: svc, validate: v}
}

func (h *userHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, errorResponse{
			Message: ErrInvalidRequestBody.Error(),
		})
		return
	}

	if errs := h.validate.Struct(body); errs != nil {
		valErrs, ok := errors.AsType[validator.ValidationErrors](errs)
		if !ok {
			panic(fmt.Sprintf("unexpected validation error type: %T", errs))
		}

		writeError(w, http.StatusBadRequest, errorResponse{
			Message: "invalid user",
			Details: formatValidationErrs(valErrs),
		})
		return
	}

	input := service.CreateUserInput{
		Name:        body.Name,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		Address: domain.Addr{
			Line1:    body.Address.Line1,
			Line2:    body.Address.Line2,
			Line3:    body.Address.Line3,
			Town:     body.Address.Town,
			County:   body.Address.County,
			PostCode: body.Address.PostCode,
		},
	}
	user, err := h.svc.CreateUser(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errorResponse{
			Message: ErrInternalSever.Error(),
		})

		return
	}

	response := UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Addr: Addr{ //
			Line1:    user.Addr.Line1,
			Line2:    user.Addr.Line2,
			Line3:    user.Addr.Line3,
			Town:     user.Addr.Town,
			County:   user.Addr.County,
			PostCode: user.Addr.PostCode,
		},
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.Format(TimeLayout),
		UpdatedAt:   user.UpdatedAt.Format(TimeLayout),
	}

	writeJSON(w, http.StatusCreated, response)

}

func (h *userHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	tokenUserID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errorResponse{
			Message: ErrUnauthorized.Error(),
		})
		return
	}

	if tokenUserID != userId {
		writeError(w, http.StatusForbidden, errorResponse{
			Message: ErrForbidden.Error(),
		})
		return
	}

	user, err := h.svc.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, errorResponse{
				Message: "user not found",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, errorResponse{
			Message: ErrInternalSever.Error(),
		})
		return
	}

	response := UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Addr: Addr{
			Line1:    user.Addr.Line1,
			Line2:    user.Addr.Line2,
			Line3:    user.Addr.Line3,
			Town:     user.Addr.Town,
			County:   user.Addr.County,
			PostCode: user.Addr.PostCode,
		},
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.Format(TimeLayout),
		UpdatedAt:   user.UpdatedAt.Format(TimeLayout),
	}

	writeJSON(w, http.StatusOK, response)
}
