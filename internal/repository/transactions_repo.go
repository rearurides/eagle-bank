package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rearurides/eagle-bank/internal/domain"
)

type transactionsRepo struct {
	db *sql.DB
}

func NewTransactionsRepo(db *sql.DB) *transactionsRepo {
	return &transactionsRepo{db: db}
}

func (r *transactionsRepo) Deposit(tan *domain.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s begin tx: %w", tan.TransactionType, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	res, err := tx.Exec(`
		UPDATE accounts
		SET balance = balance + ?,
			updated_at = ?
		WHERE id = ?
		AND balance + ? <= 1000000
		`, tan.Amount, time.Now().UTC(), tan.AccountID, tan.Amount)
	if err != nil {

		return fmt.Errorf("deposit update balance: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deposit rows effected: %w", err)
	}
	if rows == 0 {
		return domain.ErrAccountNotFound
	}

	_, err = tx.Exec(`
			INSERT INTO transactions (id, account_id, ammount, type, reference, currency, minor_unit, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		tan.TransactionID, tan.AccountID, tan.Amount, tan.TransactionType,
		tan.Reference, tan.Currency, tan.MinorUnit,
		tan.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("deposit insert transaction: %w", err)
	}

	return nil
}

func (r *transactionsRepo) Withdraw(tan *domain.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s begin tx: %w", tan.TransactionType, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	res, err := tx.Exec(`
		UPDATE accounts
		SET balance = balance - ?,
			updated_at = ?
		WHERE id = ?
		AND balance - ? >= 0
		`, tan.Amount, time.Now().UTC(), tan.AccountID, tan.Amount)
	if err != nil {
		return fmt.Errorf("deposit update balance: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deposit rows effected: %w", err)
	}
	if rows == 0 {
		return domain.ErrAccountNotFound
	}

	_, err = tx.Exec(`
			INSERT INTO transactions (id, account_id, ammount, type, reference, currency, minor_unit, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		tan.TransactionID, tan.AccountID, tan.Amount, tan.TransactionType,
		tan.Reference, tan.Currency, tan.MinorUnit,
		tan.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("withdrawal insert transaction: %w", err)
	}

	return nil
}
