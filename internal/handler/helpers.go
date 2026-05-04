package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rearurides/eagle-bank/internal/domain/validation"
)

// BadRequestMessage represents the structure of error messages returned to the client for bad requests.
type BadRequestMessage struct {
	Message string                      `json:"message"`
	Details []validation.ValidationItem `json:"details,omitempty"`
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
func writeError(w http.ResponseWriter, status int, v BadRequestMessage) {
	writeJSON(w, status, v)
}
