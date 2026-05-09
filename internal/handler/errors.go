package handler

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInternalSever      = errors.New("internal server error")
)
