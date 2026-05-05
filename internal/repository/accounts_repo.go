package repository

import (
	"database/sql"
	"fmt"

	"github.com/rearurides/eagle-bank/internal/domain"
)

type AccountsRepository struct {
	db *sql.DB
}

func NewAccountsRepo(db *sql.DB) *AccountsRepository {
	return &AccountsRepository{db: db}
}

func (r *AccountsRepository) Create(account *domain.Account) error {
	res, err := r.db.Exec(
		`INSERT INTO accounts (sort_code, user_id, name, account_type, currency, minor_unit, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		account.SortCode, account.UserID, account.Name,
		account.AccountType, account.Currency, account.MinorUnit,
		account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		// TODO: check for unique constraint violation
		return fmt.Errorf("accounts.Create insert: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("accounts.Create last insert id: %w", err)
	}

	err = r.db.QueryRow(`SELECT account_number FROM accounts WHERE id = ?`, id).
		Scan(&account.AccountNumber)
	if err != nil {
		return fmt.Errorf("failed to retrieve created account: %w", err)
	}

	return nil
}
