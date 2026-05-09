package main

import (
	"log"

	"github.com/rearurides/eagle-bank/config"
	"github.com/rearurides/eagle-bank/internal/repository"
	"github.com/rearurides/eagle-bank/internal/server"
	"github.com/rearurides/eagle-bank/internal/service"
	"github.com/rearurides/eagle-bank/pkg/db"
	"github.com/rearurides/eagle-bank/pkg/token"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database connection
	database, err := db.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := db.RunMigrations(database, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize token manager
	tm := token.NewManager(cfg.JWTSecret, cfg.JWTExpiration)

	// Repositories
	userRepo := repository.NewUserRepo(database)
	accountsRepo := repository.NewAccountsRepo(database)
	transactionsRepo := repository.NewTransactionsRepo(database)

	// Services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountsService(accountsRepo)
	transactionService := service.NewTransactionsService(transactionsRepo, accountsRepo)

	s := server.New(":"+cfg.Port, userService, accountService, transactionService, tm)
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
