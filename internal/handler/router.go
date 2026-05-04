package handler

import (
	"net/http"

	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
)

func NewRouter(userService *service.UserService) http.Handler {
	mux := http.NewServeMux()

	userHandler := newUserHandler(userService)

	mux.HandleFunc("POST /v1/users", userHandler.handleCreateUser)

	return middleware.Chain(mux,
		middleware.Logging,
		middleware.RecoverPanic,
	)
}
