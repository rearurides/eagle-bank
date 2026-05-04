package handler

import (
	"net/http"

	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
	"github.com/rearurides/eagle-bank/pkg/token"
)

func NewRouter(
	userService *service.UserService,
	tm *token.Manager,
) http.Handler {
	mux := http.NewServeMux()

	userHandler := newUserHandler(userService, tm)

	mux.HandleFunc("POST /v1/users", userHandler.handleCreateUser)
	mux.HandleFunc("POST /v1/auth/login", userHandler.handleLogin)

	protected := http.NewServeMux()
	// add protected routes
	protected.HandleFunc("GET /v1/users/{userId}", userHandler.handleGetUserByID)

	mux.Handle("/v1/users/{userId}", middleware.Auth(tm)(protected))

	return middleware.Chain(mux,
		middleware.Logging,
		middleware.RecoverPanic,
	)
}
