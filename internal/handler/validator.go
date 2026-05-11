package handler

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v.RegisterValidation("decimal2", validateDecimal2)
	return v
}

// validateDecimal2 validates to 2 decimal places. allows whole numbers
// TODO: restrict to 2 decimal places
func validateDecimal2(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val == math.Round(val*100)/100
}

// formatValidationErrs formats the validation response items
func formatValidationErrs(errs validator.ValidationErrors) []validationItem {
	items := make([]validationItem, 0, len(errs))
	for _, e := range errs {
		items = append(items, formatValidationMessage(e))
	}
	return items
}

// formatValidationMessage formats the validation message
func formatValidationMessage(e validator.FieldError) validationItem {
	switch e.Tag() {
	case "required":
		return validationItem{
			Field:   jsonPath(e),
			Message: "must not be blank",
			Type:    "required",
		}
	case "email":
		return validationItem{
			Field:   jsonPath(e),
			Message: "must be a valid email",
			Type:    "invalid_format",
		}
	case "e164":
		return validationItem{
			Field:   jsonPath(e),
			Message: "must be a valid e.164 phone number",
			Type:    "invalid_format",
		}
	case "oneof":
		opts := strings.ReplaceAll(e.Param(), " ", ", ")
		return validationItem{
			Field:   jsonPath(e),
			Message: fmt.Sprintf("must be one of: %s", opts),
			Type:    "invalid_enum",
		}
	case "decimal2":
		return validationItem{
			Field:   jsonPath(e),
			Message: "must be a valid number with up to 2 decimal places",
			Type:    "invalid_format",
		}
	case "gt":
		return validationItem{
			Field:   jsonPath(e),
			Message: fmt.Sprintf("must be a valid number greater than %s", e.Param()),
			Type:    "invalid_value",
		}
	default:
		return validationItem{
			Field:   jsonPath(e),
			Message: "failed validation",
			Type:    "invalid_value",
		}
	}
}

// jsonPath
func jsonPath(e validator.FieldError) string {
	ns := e.Namespace()
	dotIdx := strings.Index(ns, ".")
	if dotIdx == -1 {
		return e.Field()
	}
	return ns[dotIdx+1:]
}
