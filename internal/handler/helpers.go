package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const TimeLayout = "2006-01-02T15:04:05.000Z"

// badRequestErrorResponse represents the structure of error messages returned to the client for bad requests.
type badRequestErrorResponse struct {
	Message string           `json:"message"`
	Details []validationItem `json:"details,omitempty"`
}

// writeJSON is a helper function to write a JSON response with the given status code and value.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// writeError is a helper function to write an error response in JSON format.
func writeError(w http.ResponseWriter, status int, v badRequestErrorResponse) {
	writeJSON(w, status, v)
}

type validationItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type"`
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

type Money float32

func (m Money) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 8)
	buf = fmt.Appendf(buf, "%.2f", m)

	return buf, nil
}
