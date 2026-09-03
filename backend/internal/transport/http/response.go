package http

// CalculateResponse represents a successful calculation response payload.
type CalculateResponse struct {
	Result    float64   `json:"result"`
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
}

// ErrorDetail contains machine-readable code and human-readable message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents a standardized error envelope.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// HealthResponse represents service health status.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
