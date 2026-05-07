package handler

import (
	"net/http"

	"github.com/rearurides/eagle-bank/internal/handler/middleware"
	"github.com/rearurides/eagle-bank/internal/service"
	"github.com/rearurides/eagle-bank/pkg/token"
)

func NewRouter(
	userService *service.UserService,
	accountService *service.AccountsService,
	transactionsService *service.TransactionsService,
	tm *token.Manager,
) http.Handler {
	mux := http.NewServeMux()

	userHandler := newUserHandler(userService, tm)
	accountHandler := newAccountsHandler(accountService, tm)
	transactionsHandler := newTransactionsHandler(transactionsService, tm)

	mux.HandleFunc("POST /v1/users", userHandler.handleCreateUser)
	mux.HandleFunc("POST /v1/auth/login", userHandler.handleLogin)

	protected := http.NewServeMux()
	// add protected routes
	// User routes
	protected.HandleFunc("GET /v1/users/{userId}", userHandler.handleGetUserByID)
	// Account routes
	protected.HandleFunc("POST /v1/accounts", accountHandler.handleCreateAccount)
	protected.HandleFunc("GET /v1/accounts/{accountNumber}", accountHandler.handleGetAccountByNumber)

	// Transactions Routes
	protected.HandleFunc("POST /v1/accounts/{accountNumber}/transactions", transactionsHandler.handleCreateTransactions)

	mux.Handle("/v1/users/{userId}", middleware.Auth(tm)(protected))
	mux.Handle("/v1/accounts", middleware.Auth(tm)(protected))
	mux.Handle("/v1/accounts/{accountNumber}", middleware.Auth(tm)(protected))
	mux.Handle("/v1/accounts/{accountNumber}/transactions", middleware.Auth(tm)(protected))

	return middleware.Chain(mux,
		middleware.Logging,
		middleware.RecoverPanic,
	)
}
