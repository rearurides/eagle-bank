package validation

import (
	"testing"
)

func TestValidator_Required(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		message   string
		hasErrors bool
	}{
		{"valid value", "name", "John", "name is required", false},
		{"empty value", "name", "", "name is required", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator().Required(tt.field, tt.value, tt.message)
			if v.HasErrors() != tt.hasErrors {
				t.Errorf("expected HasErrors to be %v, got %v", tt.hasErrors, v.HasErrors())
			}
			if tt.hasErrors {
				if len(v.items) != 1 {
					t.Errorf("expected 1 validation item, got %d", len(v.items))
				}
				if v.items[0].Field != tt.field {
					t.Errorf("expected field %s, got %s", tt.field, v.items[0].Field)
				}
				if v.items[0].Message != tt.message {
					t.Errorf("expected message %q, got %q", tt.message, v.items[0].Message)
				}
				if v.items[0].Type != MISSING_VALUE {
					t.Errorf("expected type %s, got %s", MISSING_VALUE, v.items[0].Type)
				}
			}
		})
	}
}

func TestValidator_ValidEmail(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		message   string
		hasErrors bool
	}{
		{"valid email", "email", "john@example.com", "email is invalid", false},
		{"invalid email", "email", "Barry Gibbs <bg@example.com>", "email is invalid", true},
		{"invalid email", "email", "bg@@example.com", "email is invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator().ValidEmail(tt.field, tt.value, tt.message)
			if v.HasErrors() != tt.hasErrors {
				t.Errorf("expected HasErrors to be %v, got %v", tt.hasErrors, v.HasErrors())
			}
			if tt.hasErrors {
				if len(v.items) != 1 {
					t.Errorf("expected 1 validation item, got %d", len(v.items))
				}
				if v.items[0].Field != tt.field {
					t.Errorf("expected field %s, got %s", tt.field, v.items[0].Field)
				}
				if v.items[0].Message != tt.message {
					t.Errorf("expected message %q, got %q", tt.message, v.items[0].Message)
				}
				if v.items[0].Type != INVALID_FORMAT {
					t.Errorf("expected type %s, got %s", INVALID_FORMAT, v.items[0].Type)
				}
			}
		})
	}
}

func TestValidator_ValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		message   string
		hasErrors bool
	}{
		{"valid phone number", "phoneNumber", "+447911123456", "phone number is invalid", false},
		{"valid phone number", "phoneNumber", "+44791112345612345", "phone number is invalid", true},
		{"invalid phone number", "phoneNumber", "07911123456", "phone number is invalid", true},
		{"invalid phone number", "phoneNumber", "+44 7911 123456", "phone number is invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator().ValidPhoneNumber(tt.field, tt.value, tt.message)
			if v.HasErrors() != tt.hasErrors {
				t.Errorf("expected HasErrors to be %v, got %v", tt.hasErrors, v.HasErrors())
			}
			if tt.hasErrors {
				if len(v.items) != 1 {
					t.Errorf("expected 1 validation item, got %d", len(v.items))
				}
				if v.items[0].Field != tt.field {
					t.Errorf("expected field %s, got %s", tt.field, v.items[0].Field)
				}
				if v.items[0].Message != tt.message {
					t.Errorf("expected message %q, got %q", tt.message, v.items[0].Message)
				}
				if v.items[0].Type != INVALID_FORMAT {
					t.Errorf("expected type %s, got %s", INVALID_FORMAT, v.items[0].Type)
				}
			}
		})
	}
}

func TestValidator_Add(t *testing.T) {
	v := NewValidator().Add("field1", "error message 1").Add("field2", "error message 2")
	if !v.HasErrors() {
		t.Errorf("expected errors, got none")
	}
	if len(v.items) != 2 {
		t.Errorf("expected 2 validation items, got %d", len(v.items))
	}
}

func TestValidator_HasErrors(t *testing.T) {
	v := NewValidator()
	if v.HasErrors() {
		t.Errorf("expected no errors, got some")
	}
	v.Add("field1", "error message 1")
	if !v.HasErrors() {
		t.Errorf("expected errors, got none")
	}
}

func TestValidator_ToError(t *testing.T) {
	v := NewValidator().Add("field1", "error message 1")
	err := v.ToError("validation failed")
	if err.Message != "validation failed" {
		t.Errorf("expected error message 'validation failed', got %q", err.Message)
	}
	if len(err.Items) != 1 {
		t.Errorf("expected 1 validation item, got %d", len(err.Items))
	}
	if err.Items[0].Field != "field1" {
		t.Errorf("expected field 'field1', got %s", err.Items[0].Field)
	}
	if err.Items[0].Message != "error message 1" {
		t.Errorf("expected message 'error message 1', got %q", err.Items[0].Message)
	}
	if err.Items[0].Type != INVALID_FORMAT {
		t.Errorf("expected type %s, got %s", INVALID_FORMAT, err.Items[0].Type)
	}
}
