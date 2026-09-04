# Agent & Workspace Customizations Changelog

All notable changes to the Antigravity agent configuration, workspace rules, and skills will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.3] - 2026-09-03
### Added
- Increased React frontend test coverage to **93.42%** statements across application source code (**98.76%** UI components, **100%** API client, **100%** App root).
- Added `frontend/src/tests/api.test.ts` testing `calculate`, `CalculatorApiError`, and all backend HTTP error code mappings (422, 502, network failure).
- Added `frontend/src/tests/App.test.tsx` testing `<App />` root rendering.
- Expanded `Calculator.test.tsx` and `useCalculator.test.ts` to test all physical keyboard shortcuts (`+`, `-`, `*`, `/`, `^`, `%`, `Backspace`, `Enter`, `Escape`), unary button interactions, and repeated equals state transitions.
- Generated updated `docs/reports/frontend-coverage/index.html` and `docs/reports/frontend-test-results.txt` (29 tests, 100% PASS).

## [1.4.2] - 2026-09-03
### Added
- Increased Go backend test coverage to **98.0%** in domain engine (`internal/calculator`) and **96.6%** in HTTP transport (`internal/transport/http`).
- Added verbose test execution reports for all 100 backend assertions (`docs/reports/backend-test-results.txt`) and 15 frontend tests (`docs/reports/frontend-test-results.txt`) with pass/fail matrix in `docs/coverage-report.md`.
- Added tests for arithmetic overflow boundaries, relative percentage variations, NaN domain guards, and HTTP 500 unmapped error handling.
- Factored server creation in `cmd/api/main.go` and added `main_test.go` (100% `buildServer` coverage).

## [1.4.1] - 2026-09-03
### Added
- Comprehensive test coverage report in `docs/coverage-report.md` detailing statement coverage statistics for Go backend (78.4% - 93.3%) and React frontend (94.1% UI components).
- Coverage summary table and links directly in `README.md`.
- Integrated `@vitest/coverage-v8` in frontend with `npm run test:coverage` script.
- Configured `.gitignore` for Go and Vitest coverage artifacts.
### Changed
- Preserved last entered user input on calculator display during calculation errors, restricting error text strictly to the alert banner.

## [1.4.0] - 2026-09-03
### Added
- Deployment documentation in `docs/deploy.md` covering Cloud Build trigger setup, substitution variables, required IAM roles, and local CLI deployment workflows.
- Dynamic backend URL resolution in `cloudbuild.yaml` preventing hardcoded personal project data in version control.
- Added deployment guide link in `README.md`.
- Excluded personal environment files (`.env.*`) and `.firebaserc` in `.gitignore`.

## [1.3.0] - 2026-09-03
### Added
- Complete Go 1.23 backend microservice:
  - Clean Architecture with `cmd/api`, `internal/calculator`, and `internal/transport/http`.
  - Comprehensive operations (+, -, *, /, ^, sqrt, %).
  - Float precision normalization and IEEE-754 edge-case handling.
  - Interactive Swagger UI documentation at `/docs` backed by embedded OpenAPI 3.0 specification (`openapi.json`).
  - Complete table-driven test suite (23 unit tests).
  - Multi-stage Dockerfile (`golang:1.23-alpine` -> `alpine:3.20`).
- Complete React.js + TypeScript frontend:
  - Vite setup with strict TypeScript (`strict: true`, zero plain `.js`/`.jsx`).
  - Tailwind CSS mobile-first responsive layout (Keypad, Display, ErrorBanner).
  - Deterministic state machine via `useReducer` and `useCalculator` hook.
  - Physical keyboard navigation support (`keydown`).
  - Vitest + React Testing Library test suite (15 unit tests).
  - Multi-stage Dockerfile with Nginx serving static assets and proxying `/api/`.
- Full-Stack Docker Compose orchestration (`docker-compose.yml`) with automated health checks.
- Comprehensive root `.env.example` covering GCP, Cloud Run, and Firebase configurations.
- Production-grade `README.md` documentation.

## [1.2.0] - 2026-09-03
### Added
- React.js frontend specifications:
  - Tooling: Vite with TypeScript (**strictly mandatory**: plain `.js`/`.jsx` is prohibited, `strict: true`, `noImplicitAny: true`).
  - Styling: Tailwind CSS for responsive mobile-first layouts (CSS Grid).
  - State Architecture: `useReducer` state machine / custom hook decoupling UI from business logic.
  - Accessibility (A11y): Physical keyboard support (`keydown`), ARIA labels, and `aria-live` screen-reader display.
  - Testing: Vitest with `@testing-library/react` and `@testing-library/user-event`.
- Unit test specification document in `docs/unit-tests-spec.md` detailing critical test cases for Go backend and React frontend.
- Updated `GEMINI.md`, `code-analyzer`, and `test-runner` skills with frontend standards.

## [1.1.0] - 2026-09-03
### Added
- Go (Golang) REST API specifications: Clean Architecture (`cmd/api`, `internal/calculator`, `internal/transport/http`), pure domain with sentinel errors, float edge-case handling, and strict JSON validation.
- Deployment operations and cloud architecture:
  - Docker multi-stage build specifications for Go backend.
  - Google Cloud Run deployment configuration.
  - Firebase Hosting configuration (`firebase.json`).
  - Google Cloud Build CI/CD pipeline (`cloudbuild.yaml`) to automate tests and deployments upon git commit.
- Workspace skill `.agents/skills/deploy-ops/SKILL.md`.
- Go test standards (`go test -v -race ./...`) in `test-runner` skill.

## [1.0.0] - 2026-09-03
### Added
- Project guidelines in `GEMINI.md` enforcing English-only code, Clean Architecture, and strict repository isolation.
- Core Antigravity workspace skills in `.agents/skills/`:
  - `code-analyzer`: Standards for Clean Code, SOLID principles, import hygiene, and static typing.
  - `git-ops`: Conventional Commits in English and repository scope isolation.
  - `test-runner`: Test execution rules enforcing `poetry run pytest` and TDD practices.
- Project root `.gitignore` configured for Python and Node.js/TypeScript environments.
