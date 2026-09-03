---
name: code-analyzer
description: >-
  Use this skill when analyzing, reviewing, or writing code to ensure strict English conventions,
  SOLID principles, Clean Architecture, Go idiomatic practices, static typing, and robust REST APIs.
---

# Code Analyzer & Architecture Specialist

This skill provides guidelines and checklists to ensure high-quality, readable, and idiomatic code across the entire codebase.

## 1. Strict Language Standard
- All source code, docstrings, variable names, functions, classes, interfaces, and inline comments MUST be written exclusively in English.

## 2. Go Backend Architecture & REST API Standards
- **Standard Go Project Layout**:
  - `cmd/api/main.go`: Server initialization, configuration, dependency wiring.
  - `internal/calculator/`: Pure domain logic, isolated from transport protocols, defining domain errors.
  - `internal/transport/http/`: HTTP routing, handlers, request validation, JSON envelopes.
- **Error Handling**:
  - Use sentinel errors (`errors.New`) for known business failures (e.g. `ErrDivisionByZero`).
  - Wrap unexpected errors using `fmt.Errorf("...: %w", err)`.
  - Check errors using `errors.Is` and `errors.As`.
- **JSON Input Validation**:
  - Always decode with `dec.DisallowUnknownFields()`.
  - Return clear, standardized error payloads (e.g. RFC 7807 Problem Details or structured error envelope).
  - Map business validation errors to appropriate HTTP status codes:
    - `400 Bad Request`: Malformed JSON or syntax errors.
    - `422 Unprocessable Entity`: Semantic domain errors (e.g., division by zero).
    - `200 OK`: Successful calculation.
- **Float Precision & Edge Cases**:
  - Guard against division by zero before execution.
  - Anticipate floating-point rounding characteristics (IEEE 754).
  - Check for infinite or NaN values in results using `math.IsInf` and `math.IsNaN`.

## 3. SOLID & Clean Architecture
- **Single Responsibility Principle (SRP)**: Each module, class, or function must have one clear responsibility.
- **Open/Closed Principle (OCP)**: Calculator operations must be extensible without modifying core interfaces.
- **Liskov Substitution Principle (LSP)**: Implementations must honor interface contracts.
- **Interface Segregation Principle (ISP)**: Prefer small, focused interfaces in Go (e.g., 1-3 methods).
- **Dependency Inversion Principle (DIP)**: Handlers depend on service interfaces, not concrete structs.

## 4. Code Organization & Clean Code
- Avoid magic numbers and strings; define constants or enums.
- Avoid deep nesting (maximum 3 levels); use early exits and guard clauses.
- Limit function length to concise, cohesive operations (typically < 50 lines).

## 5. Module & Import Conventions
- Place all imports at the top of the file (module level / package declaration level).
- **NEVER** place imports inside functions, methods, or middle of files.
- Group imports logically:
  1. Standard library
  2. Third-party packages
  3. Internal project packages

## 6. React.js & Frontend Architecture
- **Framework & Tooling**: React.js with Vite and TypeScript.
- **Mandatory TypeScript Rule**:
  - **MANDATORY**: ALL frontend code (components, custom hooks, services, state reducers, utility helpers, and tests) MUST be written in TypeScript (`.ts` and `.tsx`).
  - Plain JavaScript (`.js` and `.jsx`) is **strictly prohibited**.
  - TypeScript compiler options must enforce `strict: true` and `noImplicitAny: true`.
- **State Management & Decoupling**:
  - Keep presentation components pure and lightweight.
  - Encapsulate calculation state and transitions in a `useReducer` state machine or dedicated custom hook (`useCalculator`).
  - Do not mix HTTP/fetch calls directly inside presentational components; use an isolated API service module.
- **Accessibility (A11y)**:
  - Add explicit `aria-label` attributes to calculator buttons.
  - Implement keyboard navigation (`keydown` listener for numbers, operators, Enter, Backspace, Escape).
  - Use `aria-live="polite"` on the display area to notify assistive technologies of calculation results.
- **Validation & Error Boundaries**:
  - Provide immediate client-side feedback for invalid input sequences.
  - Wrap UI in React Error Boundaries to prevent full application crashes on unhandled errors.

## 7. Static Typing
- **Go**: Strict type safety, explicit error return values (never ignore errors with `_` unless explicitly documented).
- **TypeScript**: Strict mode enabled (`strict: true`, `noImplicitAny: true`); avoid `any`. Define strong DTO interfaces matching backend payloads.
