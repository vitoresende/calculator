package calculator_test

import (
	"errors"
	"math"
	"testing"

	"vitoresende/calculator/backend/internal/calculator"
)

const floatTolerance = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

func TestAdd_PositiveAndNegativeNumbers(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive integers", 2, 3, 5},
		{"negative and positive", -4, 9, 5},
		{"two negatives", -10, -15, -25},
		{"add zero", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Add(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestAdd_FloatingPointPrecision(t *testing.T) {
	svc := calculator.NewService()

	res, err := svc.Add(0.1, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(res, 0.3) {
		t.Errorf("expected 0.3 without float artifacts, got %v", res)
	}
}

func TestAdd_Overflow(t *testing.T) {
	svc := calculator.NewService()

	_, err := svc.Add(math.MaxFloat64, math.MaxFloat64)
	if err == nil {
		t.Fatalf("expected overflow error, got nil")
	}
	if !errors.Is(err, calculator.ErrOverflow) {
		t.Errorf("expected ErrOverflow, got %v", err)
	}
}

func TestSubtract_SignFlipping(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"subtract negative", 5, -3, 8},
		{"negative minus negative", -7, -12, 5},
		{"negative minus positive", -4, 6, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Subtract(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestSubtract_ZeroIdentity(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"x minus zero", 42, 0, 42},
		{"zero minus x", 0, 42, -42},
		{"zero minus zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Subtract(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestMultiply_ZeroMultiplication(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name string
		a, b float64
	}{
		{"number by zero", 42, 0},
		{"zero by negative", 0, -100},
		{"zero by zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Multiply(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != 0.0 {
				t.Errorf("expected 0, got %f", res)
			}
		})
	}
}

func TestMultiply_NegativeSignRules(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"negative by negative", -6, -7, 42},
		{"negative by positive", -5, 8, -40},
		{"positive by negative", 3, -4, -12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Multiply(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestMultiply_Underflow(t *testing.T) {
	svc := calculator.NewService()

	res, err := svc.Multiply(1e-200, 1e-200)
	if err != nil {
		t.Fatalf("unexpected error during underflow: %v", err)
	}
	if res != 0.0 {
		t.Errorf("expected underflow to evaluate to 0.0, got %e", res)
	}
}

func TestDivide_StandardDecimals(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"fraction 1/3", 1, 3, 1.0 / 3.0},
		{"clean decimal 7/2", 7, 2, 3.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Divide(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestDivide_ByZero_ReturnsError(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name string
		a, b float64
	}{
		{"positive by zero", 10, 0},
		{"negative by zero", -5, 0},
		{"zero by zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Divide(tt.a, tt.b)
			if err == nil {
				t.Fatalf("expected ErrDivisionByZero, got nil")
			}
			if !errors.Is(err, calculator.ErrDivisionByZero) {
				t.Errorf("expected ErrDivisionByZero, got %v", err)
			}
		})
	}
}

func TestDivide_ZeroDividedByNonZero(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name string
		a, b float64
	}{
		{"zero by positive", 0, 5},
		{"zero by negative", 0, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Divide(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != 0.0 {
				t.Errorf("expected 0.0, got %f", res)
			}
		})
	}
}

func TestPow_ZeroExponent(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name string
		base float64
	}{
		{"positive base", 5},
		{"negative base", -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Pow(tt.base, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, 1.0) {
				t.Errorf("expected 1.0, got %f", res)
			}
		})
	}
}

func TestPow_ZeroBaseAndZeroExponent(t *testing.T) {
	svc := calculator.NewService()

	res, err := svc.Pow(0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(res, 1.0) {
		t.Errorf("expected 0^0 to evaluate to 1.0, got %f", res)
	}
}

func TestPow_NegativeExponent(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		base     float64
		exp      float64
		expected float64
	}{
		{"2^-2", 2, -2, 0.25},
		{"10^-3", 10, -3, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Pow(tt.base, tt.exp)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestPow_FractionalExponentNegativeBase(t *testing.T) {
	svc := calculator.NewService()

	_, err := svc.Pow(-4, 0.5)
	if err == nil {
		t.Fatalf("expected ErrInvalidDomain for fractional exponent with negative base, got nil")
	}
	if !errors.Is(err, calculator.ErrInvalidDomain) {
		t.Errorf("expected ErrInvalidDomain, got %v", err)
	}
}

func TestSqrt_PositiveValues(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		val      float64
		expected float64
	}{
		{"root 9", 9, 3},
		{"root 0.25", 0.25, 0.5},
		{"root 2", 2, math.Sqrt(2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Sqrt(tt.val)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestSqrt_Zero(t *testing.T) {
	svc := calculator.NewService()

	res, err := svc.Sqrt(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != 0.0 {
		t.Errorf("expected 0.0, got %f", res)
	}
}

func TestSqrt_NegativeValue_ReturnsError(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name string
		val  float64
	}{
		{"negative integer", -9},
		{"negative small decimal", -0.0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Sqrt(tt.val)
			if err == nil {
				t.Fatalf("expected ErrNegativeSquareRoot, got nil")
			}
			if !errors.Is(err, calculator.ErrNegativeSquareRoot) {
				t.Errorf("expected ErrNegativeSquareRoot, got %v", err)
			}
		})
	}
}

func TestPercentage_DirectValue(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		val      float64
		expected float64
	}{
		{"50 percent", 50, 0.5},
		{"100 percent", 100, 1.0},
		{"0 percent", 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Percentage(tt.val)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestPercentage_RelativeAddSubtract(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		base     float64
		pct      float64
		op       string
		expected float64
	}{
		{"100 + 10%", 100, 10, "+", 110},
		{"200 - 25%", 200, 25, "-", 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.PercentageRelative(tt.base, tt.pct, tt.op)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(res, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, res)
			}
		})
	}
}

func TestPercentageRelative_MultiplicationDivisionAndInvalid(t *testing.T) {
	svc := calculator.NewService()

	// Multiply relative: 100 * 50% = 50
	resMulWord, err := svc.PercentageRelative(100, 50, "multiply")
	if err != nil || !almostEqual(resMulWord, 50) {
		t.Errorf("expected 50 for multiply, got %f (err: %v)", resMulWord, err)
	}

	resMulSym, err := svc.PercentageRelative(100, 50, "*")
	if err != nil || !almostEqual(resMulSym, 50) {
		t.Errorf("expected 50 for *, got %f (err: %v)", resMulSym, err)
	}

	// Divide relative: 100 / 50% = 200
	resDivWord, err := svc.PercentageRelative(100, 50, "divide")
	if err != nil || !almostEqual(resDivWord, 200) {
		t.Errorf("expected 200 for divide, got %f (err: %v)", resDivWord, err)
	}

	resDivSym, err := svc.PercentageRelative(100, 50, "/")
	if err != nil || !almostEqual(resDivSym, 200) {
		t.Errorf("expected 200 for /, got %f (err: %v)", resDivSym, err)
	}

	// Invalid operation
	_, err = svc.PercentageRelative(100, 50, "invalid_op")
	if !errors.Is(err, calculator.ErrInvalidOperation) {
		t.Errorf("expected ErrInvalidOperation, got %v", err)
	}
}

func TestOverflow_AllOperations(t *testing.T) {
	svc := calculator.NewService()

	t.Run("add overflow", func(t *testing.T) {
		_, err := svc.Add(math.MaxFloat64, math.MaxFloat64)
		if !errors.Is(err, calculator.ErrOverflow) {
			t.Errorf("expected ErrOverflow, got %v", err)
		}
	})

	t.Run("subtract overflow", func(t *testing.T) {
		_, err := svc.Subtract(-math.MaxFloat64, math.MaxFloat64)
		if !errors.Is(err, calculator.ErrOverflow) {
			t.Errorf("expected ErrOverflow, got %v", err)
		}
	})

	t.Run("multiply overflow", func(t *testing.T) {
		_, err := svc.Multiply(math.MaxFloat64, 2)
		if !errors.Is(err, calculator.ErrOverflow) {
			t.Errorf("expected ErrOverflow, got %v", err)
		}
	})

	t.Run("divide overflow", func(t *testing.T) {
		_, err := svc.Divide(math.MaxFloat64, 0.5)
		if !errors.Is(err, calculator.ErrOverflow) {
			t.Errorf("expected ErrOverflow, got %v", err)
		}
	})

	t.Run("pow overflow", func(t *testing.T) {
		_, err := svc.Pow(1e200, 2)
		if !errors.Is(err, calculator.ErrOverflow) {
			t.Errorf("expected ErrOverflow, got %v", err)
		}
	})
}

func TestDomain_NaNAndEdgePrecision(t *testing.T) {
	svc := calculator.NewService()

	t.Run("pow NaN returns ErrInvalidDomain", func(t *testing.T) {
		_, err := svc.Pow(math.NaN(), 2)
		if !errors.Is(err, calculator.ErrInvalidDomain) {
			t.Errorf("expected ErrInvalidDomain for NaN base, got %v", err)
		}
	})

	t.Run("sqrt NaN returns ErrInvalidDomain", func(t *testing.T) {
		_, err := svc.Sqrt(math.NaN())
		if !errors.Is(err, calculator.ErrInvalidDomain) {
			t.Errorf("expected ErrInvalidDomain for NaN sqrt, got %v", err)
		}
	})

	t.Run("normalize sub-epsilon values to exact 0.0", func(t *testing.T) {
		res, err := svc.Add(1e-13, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != 0.0 {
			t.Errorf("expected 0.0 for sub-epsilon value, got %f", res)
		}
	})
}
