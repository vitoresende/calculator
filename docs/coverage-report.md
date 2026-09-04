# Unit Tests and Coverage Report

Comprehensive report detailing the unit test suites, execution pass/fail status, and statement/branch coverage statistics for both the **Go Backend Microservice** and the **React.js Frontend Application**.

> 📂 **Tool-Generated Report Artifacts (Generated Locally)**:
> - 🟢 **Go Backend**:
>   - [Interactive HTML Coverage Report (go tool cover)](reports/backend-coverage.html)
>   - [Statement Coverage by Function (go tool cover -func)](reports/backend-coverage.txt)
>   - [Verbose Test Execution Pass/Fail Log (go test -v)](reports/backend-test-results.txt)
> - 🔵 **React Frontend**:
>   - [Interactive HTML Coverage Dashboard (@vitest/coverage-v8)](reports/frontend-coverage/index.html)
>   - [Verbose Test Execution Pass/Fail Log (vitest --reporter=verbose)](reports/frontend-test-results.txt)
>
> *(Note: The `reports/` folder contains ephemeral build artifacts generated on local test runs and is excluded from Git tracking via `.gitignore` to prevent repository noise. All metrics and function breakdown tables below are permanently documented in this file).*

---

## 1. Executive Summary

| Layer | Framework / Tool | Test Cases / Assertions | Status | Statement Coverage | Critical Module Coverage | Execution Report |
| :--- | :--- | :---: | :---: | :---: | :---: | :--- |
| **Backend** | Go Standard Library (`testing`, `httptest`) | 34 Suites / 100 Assertions | **100% PASS** | **96.6% - 98.0%** *(Domain & Transport)* | `calculator`: **98.0%**<br>`transport/http`: **96.6%** | [backend-test-results.txt](reports/backend-test-results.txt) |
| **Frontend** | Vitest + React Testing Library + V8 | 4 Suites / 29 Test Cases | **100% PASS** | **93.4%** *(Full App Core)* | UI Components: **98.8%**<br>API Client: **100.0%**<br>`useCalculator`: **85.8%** | [frontend-test-results.txt](reports/frontend-test-results.txt) |
| **Total** | | **63 Suites / 129 Tests** | **ALL PASSING (0 Failures)** | | | |

---

## 2. Backend Unit Tests & Statement Coverage (Go)

The Go backend uses clean architecture, isolating pure mathematical domain calculation (`internal/calculator`) from HTTP transport and validation (`internal/transport/http`).

### 2.1. Coverage by Package & Function

```
ok   vitoresende/calculator/backend/internal/calculator        coverage: 98.0% of statements
ok   vitoresende/calculator/backend/internal/transport/http    coverage: 96.6% of statements
ok   vitoresende/calculator/backend/cmd/api                    coverage: 20.6% of statements (buildServer: 100.0%)
```

| Package / File | Function / Component | Coverage (%) | Description |
| :--- | :--- | :---: | :--- |
| **`internal/calculator`** | | **98.0%** | **Calculation Engine Core** |
| `service.go` | `NewService` | **100.0%** | Service constructor |
| `service.go` | `Add` | **100.0%** | Mixed signs, precision handling, arithmetic overflow detection |
| `service.go` | `Subtract` | **100.0%** | Sign flipping, identity laws, overflow boundary |
| `service.go` | `Multiply` | **100.0%** | Zero rules, negative signs, underflow, overflow boundary |
| `service.go` | `Divide` | **100.0%** | Decimals, zero divisor -> `ErrDivisionByZero`, overflow |
| `service.go` | `Pow` | **100.0%** | Real powers, zero exponent, negative base fractional domain, NaN, overflow |
| `service.go` | `Sqrt` | **100.0%** | Positive roots, zero, negative -> `ErrNegativeSquareRoot`, NaN domain check |
| `service.go` | `Percentage` | **100.0%** | Direct percentage evaluation |
| `service.go` | `PercentageRelative` | **100.0%** | Contextual additive/subtractive/multiplicative/divisive percentage |
| `service.go` | `normalizeFloat` | **88.9%** | IEEE-754 epsilon tolerance, negative zero normalization, sub-epsilon zeroing |
| **`internal/transport/http`**| | **96.6%** | **REST Transport & Error Mapping** |
| `request.go` | `Validate` | **100.0%** | Strict DTO structure and required field validation |
| `handler.go` | `NewHandler` | **100.0%** | Handler constructor with fallback default origins (`*`) |
| `handler.go` | `Routes` | **100.0%** | ServeMux routing registration (Go 1.22+) |
| `handler.go` | `handleSwaggerUI` | **100.0%** | Interactive Swagger UI documentation at `/docs` |
| `handler.go` | `handleOpenAPISpec`| **100.0%** | Embedded OpenAPI 3.0.3 specification JSON |
| `handler.go` | `handleHealth` | **100.0%** | Health check endpoints (`/health`, `/api/v1/health`) |
| `handler.go` | `handleCalculate` | **93.8%** | Payload decoding, domain execution, domain error mapping (400, 422, 500) |
| `handler.go` | `writeJSON` | **100.0%** | JSON envelope serialization |
| `handler.go` | `writeError` | **100.0%** | Standardized error envelope serialization |
| `handler.go` | `corsMiddleware` | **100.0%** | CORS headers and preflight (OPTIONS) negotiation |
| **`cmd/api`** | | | **Server Entrypoint & Factory** |
| `main.go` | `buildServer` | **100.0%** | Server initialization with environment fallback and route configuration |

---

## 3. Frontend Unit Tests & Statement Coverage (React + TypeScript)

```
 % Coverage report from v8
-------------------|---------|----------|---------|---------|-------------------------------------------------
File               | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s                               
-------------------|---------|----------|---------|---------|-------------------------------------------------
All files          |   93.42 |    89.11 |   83.87 |   93.42 |                                                 
 src/App.tsx       |  100.00 |   100.00 |  100.00 |  100.00 |                                                 
 src/components    |   98.76 |    96.15 |   78.26 |   98.76 |                                                 
  Button.tsx       |  100.00 |   100.00 |  100.00 |  100.00 |                                                 
  Calculator.tsx   |   95.34 |    91.66 |  100.00 |   95.34 | 23-24,48-49                                     
  Display.tsx      |  100.00 |   100.00 |  100.00 |  100.00 |                                                 
  ErrorBanner.tsx  |  100.00 |   100.00 |  100.00 |  100.00 |                                                 
  Keypad.tsx       |  100.00 |   100.00 |   72.22 |  100.00 |                                                 
 src/hooks         |   85.76 |    82.27 |  100.00 |   85.76 |                                                 
  useCalculator.ts |   85.76 |    82.27 |  100.00 |   85.76 | 14-15,68-69,182-188,191,255-258,264-282,302-308 
 src/services      |  100.00 |   100.00 |  100.00 |  100.00 |                                                 
  api.ts           |  100.00 |   100.00 |  100.00 |  100.00 |                                                 
-------------------|---------|----------|---------|---------|-------------------------------------------------
```

| Component / Hook / Service | Statements | Branches | Functions | Lines | Key Scenarios Tested |
| :--- | :---: | :---: | :---: | :---: | :--- |
| **`App.tsx`** | **100%** | **100%** | **100%** | **100%** | Application layout, header, footer, full container integration |
| **`Button.tsx`** | **100%** | **100%** | **100%** | **100%** | Button variants (`number`, `operator`, `function`, `action`, `equal`), accessible ARIA attributes, disabled states |
| **`Display.tsx`** | **100%** | **100%** | **100%** | **100%** | Screen reader polite announcements (`aria-live`), formula expression display |
| **`ErrorBanner.tsx`** | **100%** | **100%** | **100%** | **100%** | Alert role, Lucide alert icon, error message rendering |
| **`Keypad.tsx`** | **100%** | **100%** | 72.2% | **100%** | CSS Grid responsive keypad, all button triggers, disabled states |
| **`Calculator.tsx`** | **95.3%** | 91.7% | **100%** | **95.3%** | All physical keyboard shortcuts (`0-9`, `+`, `-`, `*`, `/`, `^`, `%`, `Enter`, `Escape`, `Backspace`, `.`) |
| **`useCalculator.ts`** | **85.8%** | 82.3% | **100%** | **85.8%** | Deterministic state machine, operator overwriting, multi-decimal prevention, display persistence on error, repeated equals |
| **`api.ts`** | **100%** | **100%** | **100%** | **100%** | Successful REST calls, all HTTP error status mappings (`DIVISION_BY_ZERO`, `NEGATIVE_SQUARE_ROOT`, `INVALID_DOMAIN`, `ARITHMETIC_OVERFLOW`), non-JSON fallbacks |

---

## 4. Test Execution Pass/Fail Status Report

### 4.1. Backend Test Suites Execution Results
All backend test suites were executed with Go's race detector enabled. Every single test passed successfully:

| Test Suite / Function | Category | Subtests / Assertions | Status | Duration |
| :--- | :--- | :--- | :---: | :---: |
| `TestBuildServer_DefaultAndCustomConfig` | API Server Factory | Default port, custom port, routes verification | **✅ PASS** | < 0.01s |
| `TestAdd_PositiveAndNegativeNumbers` | Arithmetic (Add) | Positive integers, negative & positive, two negatives, zero | **✅ PASS** | < 0.01s |
| `TestAdd_FloatingPointPrecision` | IEEE-754 Precision | Precision handling `0.1 + 0.2 = 0.3` | **✅ PASS** | < 0.01s |
| `TestAdd_Overflow` | Edge Case | Exceeding `math.MaxFloat64` yields `ErrOverflow` | **✅ PASS** | < 0.01s |
| `TestSubtract_SignFlipping` | Arithmetic (Subtract) | Sign flipping with negative operands | **✅ PASS** | < 0.01s |
| `TestSubtract_ZeroIdentity` | Identity Laws | `x - 0 = x`, `0 - x = -x`, `0 - 0 = 0` | **✅ PASS** | < 0.01s |
| `TestMultiply_ZeroMultiplication` | Arithmetic (Multiply) | Number by zero, zero by negative, zero by zero | **✅ PASS** | < 0.01s |
| `TestMultiply_NegativeSignRules` | Arithmetic (Multiply) | Negative by negative, negative by positive | **✅ PASS** | < 0.01s |
| `TestMultiply_Underflow` | Float Boundaries | Multiplying tiny subnormal values | **✅ PASS** | < 0.01s |
| `TestDivide_StandardDecimals` | Arithmetic (Divide) | Fractional values (1/3), clean decimals (7/2) | **✅ PASS** | < 0.01s |
| `TestDivide_ByZero_ReturnsError` | Sentinel Errors | Zero divisor returns `ErrDivisionByZero` | **✅ PASS** | < 0.01s |
| `TestDivide_ZeroDividedByNonZero` | Arithmetic (Divide) | Zero divided by positive, zero divided by negative | **✅ PASS** | < 0.01s |
| `TestPow_ZeroExponent` | Exponentiation | Base to power 0 equals 1 | **✅ PASS** | < 0.01s |
| `TestPow_ZeroBaseAndZeroExponent` | Exponentiation | `0^0` mathematical convention | **✅ PASS** | < 0.01s |
| `TestPow_NegativeExponent` | Exponentiation | `2^-2 = 0.25`, `10^-3 = 0.001` | **✅ PASS** | < 0.01s |
| `TestPow_FractionalExponentNegativeBase` | Domain Errors | Negative base with decimal exponent returns `ErrInvalidDomain` | **✅ PASS** | < 0.01s |
| `TestSqrt_PositiveValues` | Square Root | Exact squares, decimal roots, irrational roots | **✅ PASS** | < 0.01s |
| `TestSqrt_Zero` | Square Root | `sqrt(0) = 0` | **✅ PASS** | < 0.01s |
| `TestSqrt_NegativeValue_ReturnsError` | Sentinel Errors | Negative input returns `ErrNegativeSquareRoot` | **✅ PASS** | < 0.01s |
| `TestPercentage_DirectValue` | Percentage | 50% = 0.5, 100% = 1.0, 0% = 0.0 | **✅ PASS** | < 0.01s |
| `TestPercentage_RelativeAddSubtract` | Relative Percentage | `100 + 10% = 110`, `200 - 25% = 150` | **✅ PASS** | < 0.01s |
| `TestPercentageRelative_MultiplicationDivisionAndInvalid` | Relative Percentage | Multiply/Divide percentage and unsupported operation rejection | **✅ PASS** | < 0.01s |
| `TestOverflow_AllOperations` | Overflow Boundaries | Add, Subtract, Multiply, Divide, Pow overflow tests | **✅ PASS** | < 0.01s |
| `TestDomain_NaNAndEdgePrecision` | Edge Domain Checks | NaN handling in Pow/Sqrt, sub-epsilon precision zeroing | **✅ PASS** | < 0.01s |
| `TestHealthEndpoint` | REST HTTP API | `GET /health` returns 200 OK | **✅ PASS** | < 0.01s |
| `TestCalculateEndpoint_Success` | REST HTTP API | Add, Subtract, Multiply, Divide, Pow, Sqrt, Percentage endpoints | **✅ PASS** | < 0.01s |
| `TestAPI_PayloadValidation` | REST Security | Disallow unknown fields, non-numeric operands, malformed JSON | **✅ PASS** | < 0.01s |
| `TestCalculateEndpoint_DomainErrors` | REST HTTP Mapping | Maps domain errors to HTTP 422 JSON envelopes | **✅ PASS** | < 0.01s |
| `TestCalculateEndpoint_ArithmeticOverflow` | REST HTTP Mapping | Maps arithmetic overflow to HTTP 422 `ARITHMETIC_OVERFLOW` | **✅ PASS** | < 0.01s |
| `TestCalculateEndpoint_InternalError` | REST HTTP Mapping | Maps unexpected failures to HTTP 500 `INTERNAL_ERROR` | **✅ PASS** | < 0.01s |
| `TestCORSHeaders` | REST Security | CORS headers and preflight (OPTIONS) response | **✅ PASS** | < 0.01s |
| `TestSwaggerUIEndpoint` | Documentation API | `GET /docs` serves interactive Swagger UI HTML | **✅ PASS** | < 0.01s |
| `TestOpenAPISpecEndpoint` | Documentation API | `GET /docs/openapi.json` serves OpenAPI 3.0 specification | **✅ PASS** | < 0.01s |
| `TestNewHandler_DefaultOrigins` | Transport Config | Default origin fallback | **✅ PASS** | < 0.01s |

---

### 4.2. Frontend Test Suites Execution Results
All frontend test cases were executed via Vitest. Every single test passed successfully:

| Test Suite / Specification | Tested Behavior | Status | Duration |
| :--- | :--- | :---: | :---: |
| `useCalculator.test.ts` | `rendersInitialDisplayWithZero`: Initial state starts with 0 and clean buffers | **✅ PASS** | 12ms |
| `useCalculator.test.ts` | `preventsMultipleConsecutiveDecimals`: Pressing decimal point repeatedly retains single point | **✅ PASS** | 8ms |
| `useCalculator.test.ts` | `handlesLeadingZerosCorrectly`: Leading zero prevention (`0005` -> `5`) | **✅ PASS** | 6ms |
| `useCalculator.test.ts` | `preventsMultipleOperators`: Consecutive operators overwrite previous (`5 + * 2` -> `5 * 2`) | **✅ PASS** | 7ms |
| `useCalculator.test.ts` | `clearsStateOnAllClear_AC`: Resets inputs, memory, and operator state | **✅ PASS** | 6ms |
| `useCalculator.test.ts` | `clearsCurrentEntryOnClearEntry_CE`: Resets current buffer without losing pending calculation | **✅ PASS** | 7ms |
| `useCalculator.test.ts` | `disablesOperatorsDuringErrorState`: Operators locked during error; visor preserves typed digit | **✅ PASS** | 9ms |
| `useCalculator.test.ts` | `formatsLargeNumbersForDisplay`: Large numbers gracefully formatted via scientific notation | **✅ PASS** | 5ms |
| `useCalculator.test.ts` | `covers getOperatorSymbol for all operation types and default`: Operation symbols (+, −, ×, ÷, ^, √, %) | **✅ PASS** | 4ms |
| `useCalculator.test.ts` | `handles repeated equals evaluation with last evaluated operand and operation`: Consecutive `=` repeats op | **✅ PASS** | 5ms |
| `useCalculator.test.ts` | `recovers from error state on input decimal`: Inputting decimal clears error buffer | **✅ PASS** | 5ms |
| `useCalculator.test.ts` | `handles CLEAR_ENTRY in error state`: CE clears error banner and resets buffer | **✅ PASS** | 4ms |
| `api.test.ts` | `executes successful calculation request`: POST `/calculate` JSON payload & response mapping | **✅ PASS** | 15ms |
| `api.test.ts` | `maps DIVISION_BY_ZERO code to friendly user message`: Error mapping -> 'Cannot divide by zero' | **✅ PASS** | 8ms |
| `api.test.ts` | `maps NEGATIVE_SQUARE_ROOT code to friendly user message`: Error mapping -> 'Invalid input for square root' | **✅ PASS** | 7ms |
| `api.test.ts` | `maps INVALID_DOMAIN code to friendly user message`: Error mapping -> 'Domain error' | **✅ PASS** | 7ms |
| `api.test.ts` | `maps ARITHMETIC_OVERFLOW code to friendly user message`: Error mapping -> 'Overflow error' | **✅ PASS** | 7ms |
| `api.test.ts` | `handles non-JSON error response from backend gracefully`: Non-JSON HTTP 502 error fallback | **✅ PASS** | 6ms |
| `App.test.tsx` | `renders application header, calculator interface, and footer`: Full app rendering | **✅ PASS** | 42ms |
| `Calculator.test.tsx` | `rendersInitialDisplayWithZero and has accessible display`: Visor renders with `aria-live` | **✅ PASS** | 65ms |
| `Calculator.test.tsx` | `handles user digit clicks and decimal entry`: Button click event handling | **✅ PASS** | 82ms |
| `Calculator.test.tsx` | `preventsMultipleConsecutiveDecimals`: UI prevents duplicate decimal points | **✅ PASS** | 44ms |
| `Calculator.test.tsx` | `executes addition calculation successfully`: End-to-end user addition flow | **✅ PASS** | 110ms |
| `Calculator.test.tsx` | `displaysFriendlyDivisionByZeroMessage`: Div-by-zero shows banner alert & preserves typed digit | **✅ PASS** | 95ms |
| `Calculator.test.tsx` | `clears state on AC button press`: Clear all button resets UI components | **✅ PASS** | 38ms |
| `Calculator.test.tsx` | `supports physical keyboard input`: Handles desktop physical keyboard strokes | **✅ PASS** | 85ms |
| `Calculator.test.tsx` | `supports all operator and editing keyboard shortcuts`: +, -, *, /, ^, %, Backspace, Enter | **✅ PASS** | 92ms |
| `Calculator.test.tsx` | `executes unary square root operation on button click`: Square root button flow | **✅ PASS** | 78ms |
| `Calculator.test.tsx` | `executes direct percentage operation on button click`: Percentage button flow | **✅ PASS** | 74ms |
