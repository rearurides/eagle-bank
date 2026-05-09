package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rearurides/eagle-bank/internal/service"
	"github.com/rearurides/eagle-bank/pkg/token"
)

type Server struct {
	httpServer *http.Server
}

func New(
	addr string,
	userService *service.UserService,
	accountService *service.AccountsService,
	transactionService *service.TransactionsService,
	tm *token.Manager,
) *Server {
	router := NewRouter(
		userService,
		accountService,
		transactionService,
		tm,
	)

	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second, // Protects against Slowloris attacks
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("server starting on:  %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case <-quit:
		log.Println("shutting down sever...")
	}

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped cleanly")

	return nil
}
