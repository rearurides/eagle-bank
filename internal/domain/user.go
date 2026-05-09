package domain

import (
	"time"
)

type Addr struct {
	Line1    string
	Line2    string
	Line3    string
	Town     string
	County   string
	PostCode string
}

type User struct {
	ID          string
	Name        string
	Addr        Addr
	PhoneNumber string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUser(
	id string,
	name string,
	email string,
	phone string,
	addr Addr,
) *User {
	now := time.Now().UTC()

	return &User{
		ID:          id,
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
		Addr:        addr,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
