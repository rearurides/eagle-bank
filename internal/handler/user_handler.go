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

const TimeLayout = "2006-01-02T15:04:05.000Z"

type userHandler struct {
	service      userService
	tokenManager *token.Manager
}

func newUserHandler(svc *service.UserService, tokenManager *token.Manager) *userHandler {
	return &userHandler{service: svc, tokenManager: tokenManager}
}

type loginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *userHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, BadRequestMessage{
			Message: "invalid request body",
		})
		return
	}

	user, err := h.service.Login(service.LoginInput{Email: body.Email})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, BadRequestMessage{
				Message: domain.ErrInvalidCredentials.Error(),
			})
			return
		}

		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "internal server error",
		})
		return
	}

	tokenString, err := h.tokenManager.Generate(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "internal server error",
		})
		return
	}

	response := loginResponse{
		Token: tokenString,
		User: &UserResponse{
			ID:   user.ID,
			Name: user.Name,
			Address: Addr{ //
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
		},
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *userHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, BadRequestMessage{
			Message: "invalid request body",
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
	user, err := h.service.CreateUser(input)
	if err != nil {
		if valErr, ok := errors.AsType[*validation.ValidationError](err); ok {
			writeError(w, http.StatusBadRequest, BadRequestMessage{
				Message: valErr.Message,
				Details: valErr.Items,
			})
			return
		}

		// This could be handled better as this could be a security concern as it leaks information about existing emails in the system.
		// This is just for demonstration purposes. Could be mitigated by request timing and sending email validation.
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			writeError(w, http.StatusBadRequest, BadRequestMessage{
				Message: "invalid request body",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "internal server error",
		})

		return
	}

	response := UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Address: Addr{ //
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

func (h *userHandler) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	if userId == "" {
		writeError(w, http.StatusBadRequest, BadRequestMessage{
			Message: "user ID is required",
		})
		return
	}

	tokenUserID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, BadRequestMessage{
			Message: "unauthorized",
		})
		return
	}

	if tokenUserID != userId {
		writeError(w, http.StatusForbidden, BadRequestMessage{
			Message: "forbidden",
		})
		return
	}

	user, err := h.service.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, BadRequestMessage{
				Message: "user not found",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, BadRequestMessage{
			Message: "internal server error",
		})
		return
	}

	response := UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Address: Addr{ //
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
