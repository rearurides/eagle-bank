package handler

import (
	"net/http"

	"github.com/rearurides/eagle-bank/internal/handler/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	return middleware.Chain(mux,
		middleware.Logging,
		middleware.RecoverPanic,
	)
}
