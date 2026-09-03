---
name: test-runner
description: >-
  Use this skill when creating, configuring, or executing tests for Go backend,
  frontend applications, or Python services.
---

# Test Runner & Testing Standards

This skill outlines testing methodologies, test runner commands, and quality requirements.

## 1. Go Backend Testing
- **Command**: Run tests with race detection and verbosity:
  ```bash
  go test -v -race ./...
  ```
- **Methodology**:
  - Use **Table-Driven Tests** (`tests := []struct{ name string, ... }`) for all operations and calculation edge cases.
  - Test HTTP handlers using `net/http/httptest` (`httptest.NewRecorder()` and `httptest.NewRequest()`).
  - Cover edge cases: division by zero, negative operands, float overflow, invalid JSON bodies, extra fields.

## 2. Frontend Testing
- **Command**: Run tests using `npm test` or `npm run test:ui`.
- **Framework**: Vitest with `@testing-library/react` and `@testing-library/user-event`.
- **Structure & Practices**:
  - **Unit Tests**: Test the `useCalculator` state machine / reducer logic with pure inputs/actions.
  - **Interaction Tests**: Simulate realistic user actions using `userEvent.click` and `userEvent.keyboard`.
  - **Error & Edge Case Tests**: Test UI error feedback on server 400/422 responses (e.g. division by zero banner).
  - **Accessibility Assertions**: Ensure button labels and display ARIA attributes are tested.

## 3. Python Testing (Global Rule Compliance)
- **Mandatory Command**: ALWAYS run Python tests using `poetry run pytest`.
- **Strict Prohibition**: NEVER run `pytest` directly without `poetry`.
