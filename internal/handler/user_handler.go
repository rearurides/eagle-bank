package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rearurides/eagle-bank/internal/domain"
	"github.com/rearurides/eagle-bank/internal/domain/validation"
	"github.com/rearurides/eagle-bank/internal/service"
)

const TimeLayout = "2006-01-02T15:04:05.000Z"

type userHandler struct {
	service userService
}

func newUserHandler(svc *service.UserService) *userHandler {
	return &userHandler{service: svc}
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
