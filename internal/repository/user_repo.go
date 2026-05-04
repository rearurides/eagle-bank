package repository

import (
	"database/sql"

	"github.com/rearurides/eagle-bank/internal/domain"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

// Create inserts a new user into the database. It returns an error if the email already exists or if there is a database error.
func (r *userRepo) Create(user *domain.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, name, email, phone_number, line_1, line_2, line_3, town, county, postcode, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.PhoneNumber,
		user.Addr.Line1, user.Addr.Line2, user.Addr.Line3,
		user.Addr.Town, user.Addr.County, user.Addr.PostCode,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return domain.ErrEmailAlreadyExists
		}

		return err
	}
	return nil
}
