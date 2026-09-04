package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vitoresende/calculator/backend/internal/calculator"
	transportHttp "vitoresende/calculator/backend/internal/transport/http"
)

func setupTestServer() http.Handler {
	svc := calculator.NewService()
	handler := transportHttp.NewHandler(svc, "*")
	return handler.Routes()
}

func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp transportHttp.HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
}

func TestCalculateEndpoint_Success(t *testing.T) {
	server := setupTestServer()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedResult float64
	}{
		{
			name:           "valid addition",
			payload:        `{"operation": "add", "a": 10, "b": 25}`,
			expectedStatus: http.StatusOK,
			expectedResult: 35,
		},
		{
			name:           "valid subtraction",
			payload:        `{"operation": "subtract", "a": 50, "b": 20}`,
			expectedStatus: http.StatusOK,
			expectedResult: 30,
		},
		{
			name:           "valid multiplication",
			payload:        `{"operation": "multiply", "a": 6, "b": 7}`,
			expectedStatus: http.StatusOK,
			expectedResult: 42,
		},
		{
			name:           "valid division",
			payload:        `{"operation": "divide", "a": 20, "b": 4}`,
			expectedStatus: http.StatusOK,
			expectedResult: 5,
		},
		{
			name:           "valid exponentiation",
			payload:        `{"operation": "pow", "a": 2, "b": 8}`,
			expectedStatus: http.StatusOK,
			expectedResult: 256,
		},
		{
			name:           "valid square root",
			payload:        `{"operation": "sqrt", "a": 144}`,
			expectedStatus: http.StatusOK,
			expectedResult: 12,
		},
		{
			name:           "valid percentage",
			payload:        `{"operation": "percentage", "a": 75}`,
			expectedStatus: http.StatusOK,
			expectedResult: 0.75,
		},
		{
			name:           "valid symbol alias +",
			payload:        `{"operation": "+", "a": 3, "b": 4}`,
			expectedStatus: http.StatusOK,
			expectedResult: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			var resp transportHttp.CalculateResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Result != tt.expectedResult {
				t.Errorf("expected result %f, got %f", tt.expectedResult, resp.Result)
			}
		})
	}
}

func TestAPI_PayloadValidation(t *testing.T) {
	server := setupTestServer()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "empty body",
			payload:        "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "EMPTY_BODY",
		},
		{
			name:           "malformed json syntax",
			payload:        `{"operation": "add", "a": 10,}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MALFORMED_JSON",
		},
		{
			name:           "disallow unknown fields",
			payload:        `{"operation": "add", "a": 10, "b": 5, "extra": "forbidden"}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "UNKNOWN_FIELD",
		},
		{
			name:           "non-numeric operand type",
			payload:        `{"operation": "add", "a": "ten", "b": 5}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_FIELD_TYPE",
		},
		{
			name:           "missing operation",
			payload:        `{"a": 10, "b": 5}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "missing operand a",
			payload:        `{"operation": "add", "b": 5}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "missing operand b for binary operation",
			payload:        `{"operation": "divide", "a": 10}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "unsupported operation",
			payload:        `{"operation": "magic", "a": 10, "b": 5}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			var errResp transportHttp.ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}

			if errResp.Error.Code != tt.expectedCode {
				t.Errorf("expected error code %q, got %q", tt.expectedCode, errResp.Error.Code)
			}
		})
	}
}

func TestCalculateEndpoint_DomainErrors(t *testing.T) {
	server := setupTestServer()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "division by zero returns 422",
			payload:        `{"operation": "divide", "a": 10, "b": 0}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedCode:   "DIVISION_BY_ZERO",
		},
		{
			name:           "square root of negative returns 422",
			payload:        `{"operation": "sqrt", "a": -16}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedCode:   "NEGATIVE_SQUARE_ROOT",
		},
		{
			name:           "fractional exponent of negative base returns 422",
			payload:        `{"operation": "pow", "a": -4, "b": 0.5}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedCode:   "INVALID_DOMAIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			var errResp transportHttp.ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}

			if errResp.Error.Code != tt.expectedCode {
				t.Errorf("expected error code %q, got %q", tt.expectedCode, errResp.Error.Code)
			}
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204 for OPTIONS, got %d", rec.Code)
	}

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin '*', got %q", origin)
	}
}

func TestSwaggerUIEndpoint(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /docs, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type text/html, got %q", contentType)
	}

	if !strings.Contains(rec.Body.String(), "SwaggerUIBundle") {
		t.Errorf("expected SwaggerUIBundle script in /docs response")
	}
}

func TestOpenAPISpecEndpoint(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /docs/openapi.json, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var spec map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&spec); err != nil {
		t.Fatalf("failed to decode openapi.json: %v", err)
	}

	if spec["openapi"] != "3.0.3" {
		t.Errorf("expected openapi 3.0.3, got %v", spec["openapi"])
	}
}

func TestNewHandler_DefaultOrigins(t *testing.T) {
	svc := calculator.NewService()
	h := transportHttp.NewHandler(svc, "")
	if h == nil {
		t.Fatal("expected handler not to be nil")
	}
}

func TestCalculateEndpoint_ArithmeticOverflow(t *testing.T) {
	server := setupTestServer()

	payload := `{"operation": "pow", "a": 1e200, "b": 2}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", rec.Code)
	}

	var errResp transportHttp.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error.Code != "ARITHMETIC_OVERFLOW" {
		t.Errorf("expected ARITHMETIC_OVERFLOW, got %q", errResp.Error.Code)
	}
}

type mockFailingService struct {
	calculator.CalculatorService
}

func (m *mockFailingService) Add(a, b float64) (float64, error) {
	return 0, errors.New("unexpected unmapped domain failure")
}

func TestCalculateEndpoint_InternalError(t *testing.T) {
	mockSvc := &mockFailingService{}
	h := transportHttp.NewHandler(mockSvc, "*")
	server := h.Routes()

	payload := `{"operation": "add", "a": 1, "b": 2}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	var errResp transportHttp.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error.Code != "INTERNAL_ERROR" {
		t.Errorf("expected INTERNAL_ERROR, got %q", errResp.Error.Code)
	}
}
