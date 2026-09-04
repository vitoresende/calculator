# Modern Full-Stack Cloud Calculator

A production-grade, full-stack calculator application built with a **Go (Golang)** microservice backend and a **React.js (TypeScript)** frontend. Designed with Clean Architecture, strict static typing, comprehensive unit testing, and ready for containerized deployment to **Google Cloud Run** and **Firebase Hosting**.

> 🚀 **Looking to deploy?** See the step-by-step [Deployment Guide](docs/deploy.md) for automated Cloud Build trigger configuration, substitution variables, and local CLI deployment instructions.

---

## Architecture Overview

```
                      +-----------------------------+
                      |   React.js + TypeScript     |
                      |   (Vite + Tailwind CSS)     |
                      |   Hosted on Firebase CDN    |
                      +--------------+--------------+
                                     |
                                     | JSON REST API (/api/v1/calculate)
                                     v
                      +-----------------------------+
                      |    Go 1.23 REST Backend     |
                      |    Hosted on Cloud Run      |
                      +--------------+--------------+
                                     |
                +--------------------+--------------------+
                |                                         |
                v                                         v
     [internal/transport/http]                 [internal/calculator]
     - net/http (Go 1.22+ routing)             - Pure Domain Engine
     - Strict JSON Decoding                    - Float Precision Handling
     - Envelope & Error Mapping                - Sentinel Domain Errors
```

### Core Technologies
- **Backend**: Go 1.23, Standard Library (`net/http`, `log/slog`), Clean Architecture, Docker multi-stage build.
- **Frontend**: React 18, TypeScript (Strict Mode), Vite, Tailwind CSS, Vitest, React Testing Library.
- **DevOps & Cloud**: Docker & Docker Compose, Google Cloud Run, Firebase Hosting, Google Cloud Build (`cloudbuild.yaml`).

---

## Features & Supported Operations

| Operation | Symbol | Type | Backend Endpoint Behavior |
| :--- | :--- | :--- | :--- |
| **Addition** | `+` | Binary | Handles mixed signs, float precision (`0.1 + 0.2 = 0.3`), overflow |
| **Subtraction** | `−` | Binary | Sign flipping, identity laws |
| **Multiplication** | `×` | Binary | Zero laws, underflow protection |
| **Division** | `÷` | Binary | Explicit `ErrDivisionByZero` -> HTTP 422 |
| **Exponentiation** | `^` | Binary | Real powers, zero exponent, domain checks on negative base |
| **Square Root** | `√` | Unary | Negative input validation -> HTTP 422 |
| **Percentage** | `%` | Unary | Direct and contextual percentage evaluation |

---

## Getting Started

### Prerequisites
- **Docker & Docker Compose** (Recommended for zero-dependency execution)
- **Go 1.22+** (For local backend development)
- **Node.js 20+** and **npm** via `nvm` (For local frontend development)

---

### Quickstart with Docker Compose (Recommended)

To run the entire full-stack application (frontend + backend with reverse proxy) locally:

```bash
# 1. Clone the repository and navigate to root
cd /path/to/calculator

# 2. (Optional) Copy environment configuration
cp .env.example .env

# 3. Build and launch containers
docker compose up --build
```

- **Frontend UI**: Open [http://localhost:3001](http://localhost:3001) (or your configured `FRONTEND_PORT`)
- **Backend Health Check**: Open [http://localhost:8080/health](http://localhost:8080/health)

---

### Running Locally without Docker

#### 1. Backend (Go)
```bash
cd backend

# Run the REST API server
go run cmd/api/main.go
```
The backend will start on port `8080`.

#### 2. Frontend (React + TypeScript)
```bash
cd frontend

# Load Node 20 via nvm
source ~/.nvm/nvm.sh && nvm use 20

# Install dependencies
npm install

# Start Vite development server
npm run dev
```
The frontend will start on [http://localhost:5173](http://localhost:5173) with automatic proxying of `/api` requests to `http://localhost:8080`.

---

## Unit Tests & Coverage Report

> 📊 **Interactive Tool-Generated Reports & Test Execution Results**:
> - 📖 **Full Specification & Consolidated Metrics**: See [Comprehensive Test & Coverage Report](docs/coverage-report.md) and [Unit Tests Specification](docs/unit-tests-spec.md)
> - 🟢 **Backend (Go)**: [HTML Coverage Dashboard](docs/reports/backend-coverage.html) | [Test Log](docs/reports/backend-test-results.txt)
> - 🔵 **Frontend (Vitest)**: [HTML Coverage Dashboard](docs/reports/frontend-coverage/index.html) | [Test Log](docs/reports/frontend-test-results.txt)
>
> *(Note: Raw HTML and log artifacts in `docs/reports/` are generated locally on test runs and excluded from Git tracking via `.gitignore` to maintain repository hygiene).*

| Layer | Framework & Generator Tool | Tests / Assertions | Status | Statement Coverage | Tool-Generated Reports |
| :--- | :--- | :---: | :---: | :---: | :--- |
| **Backend** | Go `testing` + `go tool cover` | 34 Suites / 100 Tests | **100% PASS** | **98.0%** Domain / **96.6%** Transport | 📄 [Coverage HTML](docs/reports/backend-coverage.html)<br>📋 [Test Results Log](docs/reports/backend-test-results.txt) |
| **Frontend** | Vitest + `@vitest/coverage-v8` | 4 Suites / 29 Tests | **100% PASS** | **93.4%** App Core / **98.8%** Components | 📄 [Coverage HTML](docs/reports/frontend-coverage/index.html)<br>📋 [Test Results Log](docs/reports/frontend-test-results.txt) |
| **Total** | | **63 Suites / 129 Tests** | **ALL PASSING (0 Failures)** | | 📊 [Full Analysis & Matrix](docs/coverage-report.md) |

### Backend Unit Tests (Go)
Comprehensive table-driven unit tests covering all edge cases, overflow, floating-point precision, and HTTP handlers:

```bash
# Using local Go installation:
cd backend
go test -v -race ./...

# Or using Docker without installing Go on your host:
docker run --rm -v "$PWD/backend:/app" -w /app golang:1.23-alpine go test -v ./...
```

Generate interactive HTML coverage report:
```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Frontend Unit Tests (Vitest + React Testing Library)
Tests covering state machine reducers, display limits, keyboard input, and error boundaries:

```bash
cd frontend
source ~/.nvm/nvm.sh && nvm use 20

# Run all tests once:
npm test

# Generate coverage report:
npm run test:coverage
```

---

## REST API Documentation

### Interactive Swagger UI: `/docs`
- **Swagger UI**: Access the interactive API documentation and test endpoints directly at **[http://localhost:8080/docs](http://localhost:8080/docs)** (or `/docs` on your Cloud Run deployment).
- **OpenAPI 3.0 JSON Spec**: Available at `GET /docs/openapi.json`.

### Base URL: `/api/v1`

### 1. Health Check
- **Endpoint**: `GET /health` or `GET /api/v1/health`
- **Response (`200 OK`)**:
  ```json
  {
    "status": "ok",
    "service": "calculator-backend"
  }
  ```

---

### 2. Calculation Endpoint
- **Endpoint**: `POST /api/v1/calculate`
- **Content-Type**: `application/json`

#### Examples:

#### Addition (`+`)
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 0.1, "b": 0.2}'
```
**Response (`200 OK`)**:
```json
{
  "result": 0.3,
  "operation": "add",
  "operands": [0.1, 0.2]
}
```

#### Division (`/`) - Success
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 50, "b": 4}'
```
**Response (`200 OK`)**:
```json
{
  "result": 12.5,
  "operation": "divide",
  "operands": [50, 4]
}
```

#### Division by Zero - Edge Case (`422 Unprocessable Entity`)
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 10, "b": 0}'
```
**Response (`422 Unprocessable Entity`)**:
```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "division by zero is not allowed"
  }
}
```

#### Square Root (`sqrt`) - Unary
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "sqrt", "a": 144}'
```
**Response (`200 OK`)**:
```json
{
  "result": 12,
  "operation": "sqrt",
  "operands": [144]
}
```

#### Square Root of Negative Number (`422 Unprocessable Entity`)
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "sqrt", "a": -25}'
```
**Response (`422 Unprocessable Entity`)**:
```json
{
  "error": {
    "code": "NEGATIVE_SQUARE_ROOT",
    "message": "square root of negative number is not allowed in real domain"
  }
}
```

#### Exponentiation (`pow`)
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "pow", "a": 2, "b": 10}'
```
**Response (`200 OK`)**:
```json
{
  "result": 1024,
  "operation": "pow",
  "operands": [2, 10]
}
```

#### Malformed JSON or Extra Fields (`400 Bad Request`)
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 10, "b": 5, "extra": "rejected"}'
```
**Response (`400 Bad Request`)**:
```json
{
  "error": {
    "code": "UNKNOWN_FIELD",
    "message": "Request contains unknown field \"extra\""
  }
}
```

---

## Design Decisions & Technical Rationale

### 1. Backend Domain Isolation (Clean Architecture)
- The domain business logic in `internal/calculator` has **zero external or HTTP dependencies**. It operates exclusively on pure numbers and returns typed sentinel errors (`ErrDivisionByZero`, `ErrNegativeSquareRoot`, `ErrInvalidDomain`).
- The transport layer in `internal/transport/http` translates HTTP concerns (JSON decoding, status codes, CORS headers) to and from the domain layer.

### 2. Float Normalization & Precision Protection
- In standard IEEE-754 arithmetic, operations like `0.1 + 0.2` yield `0.30000000000000004`. The Go engine implements a precision normalizer with an epsilon threshold of `1e-12`, ensuring exact decimal representation while avoiding arithmetic distortion.
- IEEE-754 negative zero (`-0.0`) is automatically normalized to `0.0`.

### 3. React State Machine via `useReducer`
- Rather than scattering multiple `useState` hooks that invite race conditions, the calculator's operational lifecycle is managed via a deterministic `useReducer` state machine.
- Immediate execution chaining (e.g. `5 + 2 * 3 =`) evaluates intermediate values transparently.
- Consecutive operators overwrite the previous operator instead of crashing (e.g. `5 + * 2` evaluates as `5 * 2`).

### 4. Accessibility & UX First
- Visor uses `aria-live="polite"` so screen readers announce results dynamically.
- Global `keydown` listeners enable seamless physical keyboard usage on desktop computers.
- Mobile-first responsive CSS Grid layout styled with Tailwind CSS.

---

## Cloud Deployment Guide

> 📖 **Comprehensive Walkthrough**: For full details on IAM permissions, Cloud Build trigger substitutions, and manual CLI deployments, consult the [Deployment Guide](docs/deploy.md).

### Deploying with Google Cloud Build (`cloudbuild.yaml`)
This repository contains an automated `cloudbuild.yaml` pipeline that triggers on Git commits:
1. Runs backend unit tests (`go test -v -race ./...`).
2. Builds and pushes the Go multi-stage Docker container.
3. Deploys the Go container to **Google Cloud Run** with serverless auto-scaling.
4. Builds the React TypeScript frontend (`npm run build`).
5. Deploys static assets to **Firebase Hosting** CDN.

To manually trigger Cloud Build from your local terminal:
```bash
gcloud builds submit --config=cloudbuild.yaml .
```
