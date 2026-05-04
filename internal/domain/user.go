package domain

import (
	"time"

	"github.com/rearurides/eagle-bank/internal/domain/validation"
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
) (*User, *validation.ValidationError) {
	v := validation.NewValidator().
		Required("name", name, "name is required").
		Required("email", email, "email is required").
		Required("phoneNumber", phone, "Please enter your phone number including country code, for example +447911123456").
		Required("address.line1", addr.Line1, "address line 1 is required").
		Required("address.town", addr.Town, "town is required").
		Required("address.county", addr.County, "county is required").
		Required("address.postcode", addr.PostCode, "postcode is required")

	if v.HasErrors() {
		return nil, v.ToError("invalid user")
	}

	v.ValidEmail("email", email, "Please enter a valid email address").
		ValidPhoneNumber("phoneNumber", phone, "Your phone number doesn't look right, please include your country code, for example +447911123456")

	if v.HasErrors() {
		return nil, v.ToError("invalid user")
	}

	now := time.Now().UTC()

	return &User{
		ID:          id,
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
		Addr:        addr,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
