package server

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rearurides/eagle-bank/internal/handler"
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
	// initialise validator
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	mux := http.NewServeMux()

	userHandler := handler.NewUserHandler(userService, validate)
	loginHandler := handler.NewLoginHandler(userService, validate, *tm)
	accountHandler := handler.NewAccountsHandler(accountService, validate)
	// transactionsHandler := NewTransactionsHandler(transactionsService, tm)

	mux.HandleFunc("POST /v1/users", userHandler.HandleCreateUser)
	mux.HandleFunc("POST /v1/auth/login", loginHandler.HandleLogin)

	protected := http.NewServeMux()
	// add protected routes
	// User routes
	//protected.HandleFunc("GET /v1/users/{userId}", userHandler.handleGetUserByID)
	// Account routes
	protected.HandleFunc("POST /v1/accounts", accountHandler.HandleCreateAccount)
	//protected.HandleFunc("GET /v1/accounts/{accountNumber}", accountHandler.handleGetAccountByNumber)

	// Transactions Routes
	//protected.HandleFunc("POST /v1/accounts/{accountNumber}/transactions", transactionsHandler.handleCreateTransactions)

	//mux.Handle("/v1/users/{userId}", middleware.Auth(tm)(protected))
	//mux.Handle("/v1/accounts", middleware.Auth(tm)(protected))
	//mux.Handle("/v1/accounts/{accountNumber}", middleware.Auth(tm)(protected))
	//mux.Handle("/v1/accounts/{accountNumber}/transactions", middleware.Auth(tm)(protected))

	return middleware.Chain(mux,
		middleware.Logging,
		middleware.RecoverPanic,
	)
}
