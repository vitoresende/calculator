package calculator

import (
	"math"
)

const (
	// Epsilon represents tolerance threshold for floating point normalization.
	Epsilon = 1e-12

	// MaxSafePrecisionDecimals defines maximum rounded decimal positions.
	MaxSafePrecisionDecimals = 12
)

// CalculatorService defines business operations for arithmetic calculations.
type CalculatorService interface {
	Add(a, b float64) (float64, error)
	Subtract(a, b float64) (float64, error)
	Multiply(a, b float64) (float64, error)
	Divide(a, b float64) (float64, error)
	Pow(base, exponent float64) (float64, error)
	Sqrt(val float64) (float64, error)
	Percentage(val float64) (float64, error)
	PercentageRelative(base, pct float64, op string) (float64, error)
}

type calculatorService struct{}

// NewService instantiates a new CalculatorService implementation.
func NewService() CalculatorService {
	return &calculatorService{}
}

func (s *calculatorService) Add(a, b float64) (float64, error) {
	res := a + b
	if math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return normalizeFloat(res), nil
}

func (s *calculatorService) Subtract(a, b float64) (float64, error) {
	res := a - b
	if math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return normalizeFloat(res), nil
}

func (s *calculatorService) Multiply(a, b float64) (float64, error) {
	res := a * b
	if math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return normalizeFloat(res), nil
}

func (s *calculatorService) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	res := a / b
	if math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return normalizeFloat(res), nil
}

func (s *calculatorService) Pow(base, exponent float64) (float64, error) {
	if base < 0 && exponent != math.Floor(exponent) {
		return 0, ErrInvalidDomain
	}
	res := math.Pow(base, exponent)
	if math.IsNaN(res) {
		return 0, ErrInvalidDomain
	}
	if math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return normalizeFloat(res), nil
}

func (s *calculatorService) Sqrt(val float64) (float64, error) {
	if val < 0 {
		return 0, ErrNegativeSquareRoot
	}
	res := math.Sqrt(val)
	if math.IsNaN(res) {
		return 0, ErrInvalidDomain
	}
	return normalizeFloat(res), nil
}

func (s *calculatorService) Percentage(val float64) (float64, error) {
	res := val / 100.0
	return normalizeFloat(res), nil
}

func (s *calculatorService) PercentageRelative(base, pct float64, op string) (float64, error) {
	delta := base * (pct / 100.0)
	switch op {
	case "add", "+":
		return s.Add(base, delta)
	case "subtract", "-":
		return s.Subtract(base, delta)
	case "multiply", "*":
		return s.Multiply(base, pct/100.0)
	case "divide", "/":
		return s.Divide(base, pct/100.0)
	default:
		return 0, ErrInvalidOperation
	}
}

// normalizeFloat cleans floating point artifacts, normalizes negative zero,
// and rounds to safe decimal positions within small epsilon tolerances.
func normalizeFloat(val float64) float64 {
	// Eliminate negative zero (-0.0 -> 0.0)
	if val == 0.0 {
		return 0.0
	}

	factor := math.Pow(10, float64(MaxSafePrecisionDecimals))
	rounded := math.Round(val*factor) / factor

	if math.Abs(val-rounded) < Epsilon {
		if rounded == 0.0 {
			return 0.0
		}
		return rounded
	}

	return val
}
