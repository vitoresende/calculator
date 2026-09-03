package http

import (
	"errors"
	"strings"
)

var (
	ErrMissingOperation = errors.New("field 'operation' is required")
	ErrMissingOperandA  = errors.New("field 'a' is required")
	ErrMissingOperandB  = errors.New("field 'b' is required for binary operation")
	ErrUnknownOperation = errors.New("unknown or unsupported operation")
)

// CalculateRequest represents the incoming JSON payload for calculation requests.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b,omitempty"`
}

// Validate ensures all required parameters for the specified operation are present.
func (r *CalculateRequest) Validate() error {
	op := strings.ToLower(strings.TrimSpace(r.Operation))
	if op == "" {
		return ErrMissingOperation
	}

	if r.A == nil {
		return ErrMissingOperandA
	}

	switch op {
	case "add", "+", "subtract", "-", "multiply", "*", "divide", "/", "pow", "^":
		if r.B == nil {
			return ErrMissingOperandB
		}
	case "sqrt", "percentage", "%":
		// Unary operations only require operand A
	default:
		return ErrUnknownOperation
	}

	return nil
}
