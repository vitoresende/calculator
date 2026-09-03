package http

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"vitoresende/calculator/backend/internal/calculator"
)

//go:embed openapi.json
var openAPISpec []byte

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Calculator API - Swagger UI Documentation</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; background: #f8fafc; }
    .swagger-ui .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: "/docs/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

// Handler manages routing and HTTP requests for the calculator API.
type Handler struct {
	calculatorService calculator.CalculatorService
	allowedOrigins    string
}

// NewHandler initializes a new Handler instance.
func NewHandler(svc calculator.CalculatorService, allowedOrigins string) *Handler {
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}
	return &Handler{
		calculatorService: svc,
		allowedOrigins:    allowedOrigins,
	}
}

// Routes configures the HTTP router with all API endpoints and middlewares.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.handleHealth)
	mux.HandleFunc("GET /api/v1/health", h.handleHealth)
	mux.HandleFunc("POST /api/v1/calculate", h.handleCalculate)

	// Swagger UI and OpenAPI documentation endpoints
	mux.HandleFunc("GET /docs", h.handleSwaggerUI)
	mux.HandleFunc("GET /docs/", h.handleSwaggerUI)
	mux.HandleFunc("GET /docs/openapi.json", h.handleOpenAPISpec)
	mux.HandleFunc("GET /openapi.json", h.handleOpenAPISpec)

	return h.corsMiddleware(mux)
}

func (h *Handler) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerUIHTML))
}

func (h *Handler) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "calculator-backend",
	})
}

func (h *Handler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	// Limit request body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req CalculateRequest
	if err := dec.Decode(&req); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			h.writeError(w, http.StatusBadRequest, "MALFORMED_JSON", fmt.Sprintf("Request body contains malformed JSON at position %d", syntaxErr.Offset))
		case errors.As(err, &unmarshalTypeErr):
			h.writeError(w, http.StatusBadRequest, "INVALID_FIELD_TYPE", fmt.Sprintf("Field %q must be a valid number", unmarshalTypeErr.Field))
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			h.writeError(w, http.StatusBadRequest, "UNKNOWN_FIELD", fmt.Sprintf("Request contains unknown field %s", fieldName))
		case errors.Is(err, io.EOF):
			h.writeError(w, http.StatusBadRequest, "EMPTY_BODY", "Request body cannot be empty")
		default:
			h.writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		}
		return
	}

	// Validate semantics
	if err := req.Validate(); err != nil {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	op := strings.ToLower(strings.TrimSpace(req.Operation))
	var result float64
	var err error
	var operands []float64

	switch op {
	case "add", "+":
		operands = []float64{*req.A, *req.B}
		result, err = h.calculatorService.Add(*req.A, *req.B)
	case "subtract", "-":
		operands = []float64{*req.A, *req.B}
		result, err = h.calculatorService.Subtract(*req.A, *req.B)
	case "multiply", "*":
		operands = []float64{*req.A, *req.B}
		result, err = h.calculatorService.Multiply(*req.A, *req.B)
	case "divide", "/":
		operands = []float64{*req.A, *req.B}
		result, err = h.calculatorService.Divide(*req.A, *req.B)
	case "pow", "^":
		operands = []float64{*req.A, *req.B}
		result, err = h.calculatorService.Pow(*req.A, *req.B)
	case "sqrt":
		operands = []float64{*req.A}
		result, err = h.calculatorService.Sqrt(*req.A)
	case "percentage", "%":
		operands = []float64{*req.A}
		result, err = h.calculatorService.Percentage(*req.A)
	default:
		h.writeError(w, http.StatusBadRequest, "INVALID_OPERATION", "Operation is not supported")
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, calculator.ErrDivisionByZero):
			h.writeError(w, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", err.Error())
		case errors.Is(err, calculator.ErrNegativeSquareRoot):
			h.writeError(w, http.StatusUnprocessableEntity, "NEGATIVE_SQUARE_ROOT", err.Error())
		case errors.Is(err, calculator.ErrInvalidDomain):
			h.writeError(w, http.StatusUnprocessableEntity, "INVALID_DOMAIN", err.Error())
		case errors.Is(err, calculator.ErrOverflow):
			h.writeError(w, http.StatusUnprocessableEntity, "ARITHMETIC_OVERFLOW", err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred during calculation")
		}
		return
	}

	h.writeJSON(w, http.StatusOK, CalculateResponse{
		Result:    result,
		Operation: op,
		Operands:  operands,
	})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, code, message string) {
	h.writeJSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func (h *Handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", h.allowedOrigins)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
