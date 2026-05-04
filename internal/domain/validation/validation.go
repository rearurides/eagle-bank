package validation

import (
	"net/mail"
	"regexp"
	"strings"
)

var e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

type Validator struct {
	items []ValidationItem
}

func NewValidator() *Validator {
	return &Validator{}
}

// Required checks if the provided value is not empty. If it is empty, it adds a validation error for the specified field.
func (v *Validator) Required(field, value, message string) *Validator {
	if value == "" {
		v.items = append(v.items, ValidationItem{
			Field:   field,
			Message: message,
			Type:    MISSING_VALUE,
		})
	}
	return v
}

// ValidEmail basic checks if the provided value is a valid email address.
// It uses the net/mail package for parsing and also checks for common invalid characters.
func (v *Validator) ValidEmail(field, value, message string) *Validator {
	if strings.ContainsAny(value, "<>") {
		v.items = append(v.items, ValidationItem{
			Field:   field,
			Message: message,
			Type:    INVALID_FORMAT,
		})
		return v
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		v.items = append(v.items, ValidationItem{
			Field:   field,
			Message: message,
			Type:    INVALID_FORMAT,
		})
	}
	return v
}

// ValidPhoneNumber performs basic phone number validation for E.164 format
func (v *Validator) ValidPhoneNumber(field, value, message string) *Validator {
	if !e164Regex.MatchString(value) {
		v.items = append(v.items, ValidationItem{
			Field:   field,
			Message: message,
			Type:    INVALID_FORMAT,
		})
	}
	return v
}

// Add allows adding a custom validation error for a specific field with a message.
func (v *Validator) Add(field, message string) *Validator {
	v.items = append(v.items, ValidationItem{
		Field:   field,
		Message: message,
		Type:    INVALID_FORMAT,
	})
	return v
}

// HasErrors checks if there are any validation errors collected in the validator.
func (v *Validator) HasErrors() bool {
	return len(v.items) > 0
}

// ToError converts the collected validation errors into a ValidationError with a general message.
func (v *Validator) ToError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
		Items:   v.items,
	}
}
