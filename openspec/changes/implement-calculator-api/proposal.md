# Proposal: Implement Calculator API

## Why

The frontend must talk to the backend, which is the single source of truth for calculation. A small, stable REST contract decouples the UI from the engine and makes both independently testable.

## What Changes

- Add package `backend/internal/api` with an `http.Handler` (`NewHandler()`), `POST /api/calculate`, and `GET /api/health`.
- Define JSON request/response types, request validation, status codes, and error mapping from `calculator` sentinel errors to stable client messages.
- Wire the handler into `cmd/server` (replacing the empty mux from the foundation).
- Add `httptest`-based table-driven handler tests.
- Handlers only decode, call `calculator.Evaluate`, map errors, and encode; no calculation logic.

## Capabilities

### New Capabilities
- `calculator-api`: HTTP contract for evaluating expressions and checking service health.

### Modified Capabilities
None.

## Non-goals

CORS (the Vite proxy and same-origin production serving make it unnecessary), authentication, rate limiting, versioned routes, logging middleware frameworks, static file serving (owned by `add-quality-documentation-and-docker`).

## Impact

- New: `backend/internal/api/*.go`, tests; modified: `backend/cmd/server/main.go`.
- Depends on: `implement-go-expression-engine`. Consumed by: `build-live-calculator-ui`, Docker/CD health checks.
