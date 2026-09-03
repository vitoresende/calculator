# Unit Testing Specification: Backend (Go) & Frontend (React + TypeScript)

This document specifies the critical unit tests required for the Calculator system, covering the **Go Core Engine & REST API** on the backend and the **React.js + TypeScript UI & State Machine** on the frontend.

---

## 1. Test Execution & Tooling

| Scope | Tooling | Execution Command |
| :--- | :--- | :--- |
| **Backend (Go)** | Go Standard `testing`, `net/http/httptest` | `go test -v -race ./...` |
| **Frontend (React)** | Vitest, `@testing-library/react`, `@testing-library/user-event` | `npm test` |

---

## 2. Backend Unit Tests (Go - Core Engine & API)

The Go backend must utilize **Table-Driven Tests** (`tests := []struct{ name string, ... }`) to ensure comprehensive coverage across normal operations, boundary values, and mathematical edge cases.

### 2.1. Addition (`+`)

- **`TestAdd_PositiveAndNegativeNumbers`**
  - **Objective**: Verify standard arithmetic addition for mixed signs.
  - **Test Cases**:
    - `2 + 3 = 5`
    - `-4 + 9 = 5`
    - `-10 + (-15) = -25`
  - **Expected Outcome**: Exact integer-equivalent floating-point results.

- **`TestAdd_FloatingPointPrecision`**
  - **Objective**: Prevent standard IEEE-754 precision artifacts.
  - **Test Cases**:
    - `0.1 + 0.2` must evaluate to `0.3` within an epsilon tolerance (e.g., `1e-9`) or via fixed-precision math (`math/big` / decimal).
  - **Expected Outcome**: No raw artifact representations like `0.30000000000000004` exposed to consumers.

- **`TestAdd_Overflow`**
  - **Objective**: Verify system behavior when results exceed `math.MaxFloat64`.
  - **Test Cases**:
    - `math.MaxFloat64 + math.MaxFloat64`
  - **Expected Outcome**: Controlled behavior (explicit domain error `ErrOverflow` or controlled `+Inf` with detection via `math.IsInf`).

---

### 2.2. Subtraction (`-`)

- **`TestSubtract_SignFlipping`**
  - **Objective**: Verify subtraction when subtracting negative operands.
  - **Test Cases**:
    - `5 - (-3) = 8`
    - `-7 - (-12) = 5`
    - `-4 - 6 = -10`
  - **Expected Outcome**: Correct algebraic signs.

- **`TestSubtract_ZeroIdentity`**
  - **Objective**: Verify the identity property of subtraction.
  - **Test Cases**:
    - `x - 0 = x`
    - `0 - x = -x`
    - `0 - 0 = 0`
  - **Expected Outcome**: Identical original values without inversion errors.

---

### 2.3. Multiplication (`*`)

- **`TestMultiply_ZeroMultiplication`**
  - **Objective**: Multiplying any valid operand by zero must yield zero.
  - **Test Cases**:
    - `42 * 0 = 0`
    - `0 * (-100) = 0`
    - `0 * 0 = 0`
  - **Expected Outcome**: `0` (handling potential `-0.0` normalization to `0.0`).

- **`TestMultiply_NegativeSignRules`**
  - **Objective**: Verify sign propagation rules.
  - **Test Cases**:
    - `(-6) * (-7) = 42`
    - `(-5) * 8 = -40`
  - **Expected Outcome**: Positive outcome for even negative counts; negative for odd.

- **`TestMultiply_Underflow`**
  - **Objective**: Multiplying very small values approaching zero without unexpected panic or corruption.
  - **Test Cases**:
    - `1e-200 * 1e-200`
  - **Expected Outcome**: Evaluates to `0.0` gracefully without panic.

---

### 2.4. Division (`/`)

- **`TestDivide_StandardDecimals`**
  - **Objective**: Verify non-terminating and standard decimal divisions.
  - **Test Cases**:
    - `1 / 3 = 0.3333333333333333`
    - `7 / 2 = 3.5`
  - **Expected Outcome**: Consistent precision representation up to standard float limits.

- **`TestDivide_ByZero_ReturnsError`**
  - **Objective**: Dividing by zero must return an explicit sentinel error rather than propagating `+Inf` or causing a runtime panic.
  - **Test Cases**:
    - `10 / 0`
    - `-5 / 0`
    - `0 / 0`
  - **Expected Outcome**: Returns `ErrDivisionByZero`. When exposed via HTTP handler, maps to HTTP `422 Unprocessable Entity` or `400 Bad Request`.

- **`TestDivide_ZeroDividedByNonZero`**
  - **Objective**: Zero in the numerator with non-zero denominator.
  - **Test Cases**:
    - `0 / 5 = 0`
    - `0 / (-10) = 0`
  - **Expected Outcome**: `0.0` without errors.

---

### 2.5. Exponentiation (`^` / `pow`)

- **`TestPow_ZeroExponent`**
  - **Objective**: Any non-zero base raised to power 0 equals 1.
  - **Test Cases**:
    - `5 ^ 0 = 1`
    - `(-3) ^ 0 = 1`
  - **Expected Outcome**: `1.0`.

- **`TestPow_ZeroBaseAndZeroExponent`**
  - **Objective**: Clarify system agreement on the mathematically indeterminate `0 ^ 0`.
  - **Test Cases**:
    - `0 ^ 0`
  - **Expected Outcome**: Standard library evaluation to `1` or explicit domain validation decision.

- **`TestPow_NegativeExponent`**
  - **Objective**: Inverted fractional results for negative exponents.
  - **Test Cases**:
    - `2 ^ (-2) = 0.25`
    - `10 ^ (-3) = 0.001`
  - **Expected Outcome**: Exact fractional equivalents.

- **`TestPow_FractionalExponentNegativeBase`**
  - **Objective**: Prevent silent complex/imaginary number corruptions when real numbers are expected.
  - **Test Cases**:
    - `(-4) ^ 0.5`
  - **Expected Outcome**: Explicit domain error `ErrInvalidDomain` rather than silent `NaN`.

---

### 2.6. Square Root (`√`)

- **`TestSqrt_PositiveValues`**
  - **Objective**: Standard roots and fractional inputs.
  - **Test Cases**:
    - `√9 = 3`
    - `√0.25 = 0.5`
    - `√2 ≈ 1.41421356237`
  - **Expected Outcome**: Accurate root computation.

- **`TestSqrt_Zero`**
  - **Objective**: Square root of zero.
  - **Test Cases**:
    - `√0 = 0`
  - **Expected Outcome**: `0.0`.

- **`TestSqrt_NegativeValue_ReturnsError`**
  - **Objective**: Square root of negative values in the real number domain.
  - **Test Cases**:
    - `√(-9)`
    - `√(-0.0001)`
  - **Expected Outcome**: Returns explicit `ErrNegativeSquareRoot` domain error.

---

### 2.7. Percentage (`%`)

- **`TestPercentage_DirectValue`**
  - **Objective**: Direct unary percentage transformation.
  - **Test Cases**:
    - `50% = 0.5`
    - `100% = 1.0`
    - `0% = 0.0`
  - **Expected Outcome**: Divides operand by 100.

- **`TestPercentage_RelativeAddSubtract`**
  - **Objective**: Context-aware percentage evaluation in composite expressions.
  - **Test Cases**:
    - `100 + 10% = 110`
    - `200 - 25% = 150`
  - **Expected Outcome**: Percentage scaled relative to the primary operand.

---

### 2.8. API & Parser Layer (Go HTTP Handlers)

- **`TestParse_MalformedExpression`**
  - **Objective**: Reject invalid mathematical syntax gracefully.
  - **Test Cases**:
    - Consecutive decimals: `5..2 + 1`
    - Missing operators: `5 5 + 2`
    - Unbalanced parentheses: `(5 + 3 * 2`
  - **Expected Outcome**: HTTP `400 Bad Request` with descriptive error payload:
    ```json
    {
      "error": {
        "code": "MALFORMED_EXPRESSION",
        "message": "Unbalanced parentheses at position 0"
      }
    }
    ```

- **`TestAPI_PayloadValidation`**
  - **Objective**: Validate HTTP boundary decoding, unknown fields, and types.
  - **Test Cases**:
    - Empty request body
    - Malformed JSON syntax (e.g., trailing comma)
    - Unknown fields with `DisallowUnknownFields()` (e.g., `{"operation": "add", "extra": true}`)
    - Non-numeric operand types (e.g., `{"a": "invalid", "b": 2}`)
  - **Expected Outcome**: HTTP `400 Bad Request` with validation diagnostics.

---

## 3. Frontend Unit Tests (React.js + TypeScript)

Frontend unit tests focus on **UI interactions**, **the state machine (`useReducer`)**, and **accessibility** using Vitest and React Testing Library.

### 3.1. Input & Display Management

- **`rendersInitialDisplayWithZero`**
  - **Test**: Mount the `<Calculator />` component.
  - **Assertion**: Display element contains `"0"` and is accessible via `aria-live`.

- **`preventsMultipleConsecutiveDecimals`**
  - **Test**: Simulate user clicks: `5`, `.`, `.`, `3`.
  - **Assertion**: Display shows `"5.3"` (subsequent dots are ignored within the same operand).

- **`handlesLeadingZerosCorrectly`**
  - **Test**: Simulate clicking `0`, `0`, `0`, `4`.
  - **Assertion**: Display shows `"4"`.
  - **Test Variation**: Simulate `0`, `.`, `4` -> display shows `"0.4"`.

- **`preventsMultipleOperators`**
  - **Test**: Simulate entering `5`, `+`, `*`, `2`.
  - **Assertion**: The second operator `*` overrides the initial operator `+`. The final calculation computes `5 * 2 = 10`.

---

### 3.2. State & Chain Operations

- **`chainsMultipleOperations`**
  - **Test**: Simulate consecutive chained calculations: `5`, `+`, `2`, `*`, `3`, `=`.
  - **Assertion**: Correctly computes following the defined evaluation model (immediate execution: `(5 + 2) * 3 = 21` or operator precedence: `5 + (2 * 3) = 11`).

- **`handlesRepeatedEquals`**
  - **Test**: Enter `5`, `+`, `2`, `=`, `=`, `=`.
  - **Assertion**:
    - First `=`: `5 + 2 = 7`
    - Second `=`: `7 + 2 = 9`
    - Third `=`: `9 + 2 = 11`

- **`clearsStateOnAllClear_AC`**
  - **Test**: Enter `5`, `+`, `10`, then click `AC`.
  - **Assertion**: Resets display to `"0"`, clears current buffer, pending operations, and calculation history.

- **`clearsCurrentEntryOnClearEntry_CE`**
  - **Test**: Enter `15`, `+`, `99`, then click `CE`, then enter `5`, `=`.
  - **Assertion**: Only `"99"` is discarded; computes `15 + 5 = 20`.

---

### 3.3. Error & Boundary Rendering

- **`displaysFriendlyDivisionByZeroMessage`**
  - **Test**: Enter `10`, `/`, `0`, `=`. Mock backend response returning HTTP `422` / `400` with code `DIVISION_BY_ZERO`.
  - **Assertion**: Visor displays `"Cannot divide by zero"` or `"Error"`, and `<ErrorBanner />` is rendered.

- **`disablesOperatorsDuringErrorState`**
  - **Test**: While in an error state, click `+`, `-`, `*`, `/`.
  - **Assertion**: Operators are ignored. The only valid escape actions are clicking `AC` or typing a new numeric digit.

- **`formatsLargeNumbersForDisplay`**
  - **Test**: Compute a number exceeding the visible character width (e.g., `100000000000000` or `1 / 3`).
  - **Assertion**: Output is formatted using scientific notation (e.g., `1.0e14`) or fixed precision rounding to prevent container overflow.

---

## 4. Test Implementation Patterns

### 4.1. Go Table-Driven Test Pattern Example
```go
package calculator_test

import (
    "errors"
    "math"
    "testing"

    "vitoresende/calculator/internal/calculator"
)

func TestDivide(t *testing.T) {
    tests := []struct {
        name        string
        a, b        float64
        expected    float64
        expectedErr error
    }{
        {
            name:        "standard division",
            a:           10,
            b:           2,
            expected:    5,
            expectedErr: nil,
        },
        {
            name:        "division by zero",
            a:           5,
            b:           0,
            expected:    0,
            expectedErr: calculator.ErrDivisionByZero,
        },
        {
            name:        "zero divided by number",
            a:           0,
            b:           5,
            expected:    0,
            expectedErr: nil,
        },
    }

    svc := calculator.NewService()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            res, err := svc.Divide(tt.a, tt.b)

            if tt.expectedErr != nil {
                if !errors.Is(err, tt.expectedErr) {
                    t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if math.Abs(res-tt.expected) > 1e-9 {
                t.Errorf("expected %f, got %f", tt.expected, res)
            }
        })
    }
}
```

### 4.2. React (Vitest + Testing Library) Pattern Example
```tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect } from 'vitest';
import { Calculator } from '../components/Calculator';

describe('Calculator Component Unit Tests', () => {
  it('prevents multiple consecutive decimals', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^5$/ }));
    await user.click(screen.getByRole('button', { name: /decimal point/i }));
    await user.click(screen.getByRole('button', { name: /decimal point/i }));
    await user.click(screen.getByRole('button', { name: /^3$/ }));

    const display = screen.getByTestId('calculator-display');
    expect(display).toHaveTextContent('5.3');
  });
});
```
