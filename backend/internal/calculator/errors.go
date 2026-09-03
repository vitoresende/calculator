package calculator

import "errors"

var (
	// ErrDivisionByZero indicates an attempt to divide any value by zero.
	ErrDivisionByZero = errors.New("division by zero is not allowed")

	// ErrNegativeSquareRoot indicates an attempt to compute the square root of a negative number.
	ErrNegativeSquareRoot = errors.New("square root of negative number is not allowed in real domain")

	// ErrInvalidDomain indicates a mathematical operation outside valid real domain.
	ErrInvalidDomain = errors.New("operation is undefined in real domain")

	// ErrOverflow indicates an operation resulting in a value exceeding maximum representable float.
	ErrOverflow = errors.New("arithmetic operation resulted in float overflow")

	// ErrInvalidOperation indicates an unsupported or unknown calculator operation.
	ErrInvalidOperation = errors.New("unsupported calculation operation")
)
