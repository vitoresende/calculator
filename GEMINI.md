# Project Guidelines & Agent Directives

## 1. Repository Boundary
- All changes, file creations, and edits must be strictly confined to this repository (`vitoresende/calculator`).
- Never modify or create files outside this directory unless explicitly instructed.

## 2. Language Standard (Strict English)
- **ALL** source code, docstrings, variable names, functions, classes, interfaces, and inline comments MUST be written exclusively in English.
- Commit messages, documentation, and pull request notes must be in English.

## 3. Technology Stack & Architecture
- **Backend**:
  - Language: **Go (Golang)** with modern standard library (`net/http` Go 1.22+).
  - Architecture: Clean Architecture / Standard Go Project Layout:
    - `cmd/api/main.go`: Application entrypoint and dependency injection.
    - `internal/calculator/`: Pure domain calculation engine, zero HTTP dependencies, sentinel errors.
    - `internal/transport/http/`: REST HTTP handlers, strict JSON payload decoding, validation, standardized responses.
  - Calculations & Edge Cases: Explicit handling of division by zero, float precision, and invalid inputs.
- **Frontend**:
  - Framework: **React.js** with **Vite** and **TypeScript** (**STRICTLY MANDATORY**: All components, hooks, utilities, and tests MUST be written in TypeScript (`.ts`/`.tsx`). Plain JavaScript (`.js`/`.jsx`) is strictly prohibited. Strict mode (`strict: true`) must be enabled).
  - Styling: **Tailwind CSS** for responsive design (mobile-first layout, CSS Grid).
  - Architecture: Separation of concerns (UI components, `useReducer` state machine / custom hooks, typed API client).
  - Accessibility & UX: Physical keyboard navigation (`keydown`), ARIA labels, `aria-live` for screen readers, and clear error banners.
  - Hosting: Designed for deployment to **Firebase Hosting**.
- **Infrastructure & CI/CD**:
  - **Docker**: Multi-stage minimal container for the Go backend (e.g. `gcr.io/distroless/static` or Alpine).
  - **Cloud Run**: Fully managed serverless container runtime for the Go REST API.
  - **Firebase Hosting**: Fast global CDN for the frontend web application.
  - **Cloud Build (`cloudbuild.yaml`)**: Automated pipeline to test, build, and deploy both frontend and backend upon commit.

## 4. Code Quality & Idiomatic Design
- Write clean, readable, and idiomatic code adhering to SOLID principles and Clean Architecture.
- Maintain high cohesion, low coupling, and clear separation of concerns (frontend and backend).
- Ensure strict static typing (Go static types, TypeScript strict mode).
- Avoid magic numbers and strings; use constants, enums, or sentinel errors.
- Use explicit error handling and guard clauses.

## 5. Module & Import Conventions
- **Imports**: ALWAYS place imports at the top of the file (module level / package level).
- **NEVER** add imports in the middle of code, inside methods, or at the start of functions.
- In Go: group imports (standard library first, third-party libraries second, internal packages third).

## 6. Testing & Verification
- **Go Backend**: Run tests using `go test -v -race ./...`. Use table-driven tests and `net/http/httptest`.
- **Python (if applicable)**: ALWAYS run tests using `poetry run pytest`. NEVER run `pytest` directly without `poetry`.
- **Frontend**: Run tests using `npm test` / Vitest.
- Maintain high test coverage for business logic and edge cases.

## 7. Antigravity Skills & Versioning
- Project skills are defined in [`.agents/skills/`](file:///.agents/skills).
- Consult [`.agents/CHANGELOG.md`](file:///.agents/CHANGELOG.md) for version history of rules and agent configurations.
