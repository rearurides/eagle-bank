package validation

type ValueType string

const (
	INVALID_FORMAT = "INVALID_FORMAT"
	MISSING_VALUE  = "MISSING_REQUIRED_FIELD"
	CONFLICT       = "CONFLICT"
)

// ValidationItem represents a single validation error for a specific field.
type ValidationItem struct {
	Field   string    `json:"field"`
	Message string    `json:"message"`
	Type    ValueType `json:"type"`
}

// ValidationError represents a collection of validation errors for an entity.
type ValidationError struct {
	Message string           `json:"message"`
	Items   []ValidationItem `json:"items"`
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return e.Message
}
