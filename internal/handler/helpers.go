package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const TimeLayout = "2006-01-02T15:04:05.000Z"

// errorResponse represents the structure of error messages returned to the client for bad requests.
type errorResponse struct {
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
func writeError(w http.ResponseWriter, status int, v errorResponse) {
	writeJSON(w, status, v)
}

type Money float32

func (m Money) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 8)
	buf = fmt.Appendf(buf, "%.2f", m)

	return buf, nil
}
