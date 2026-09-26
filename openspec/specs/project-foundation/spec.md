# project-foundation Specification

## Purpose
TBD - created by archiving change bootstrap-fullstack-foundation. Update Purpose after archive.

## Requirements

### Requirement: Repository layout
The repository SHALL keep the frontend in `frontend/` and the backend in `backend/`, each with its own dependency manifest, and SHALL NOT introduce a monorepo framework.

#### Scenario: Independent toolchains
- **WHEN** a developer runs frontend commands inside `frontend/` and backend commands inside `backend/`
- **THEN** each toolchain works without requiring the other to be installed or built

### Requirement: Frontend toolchain
The frontend SHALL be a React + TypeScript + Vite application styled with Tailwind CSS utilities, using a permanent dark base with no theme switcher, and SHALL provide `dev`, `build`, `lint`, and `test` npm scripts.

#### Scenario: Frontend checks pass on the bare foundation
- **WHEN** `npm run lint`, `npm test -- --run`, and `npm run build` are run in `frontend/`
- **THEN** all three succeed with no errors

#### Scenario: Tailwind is active
- **WHEN** the placeholder page is rendered in the dev server
- **THEN** Tailwind utility classes on the page take effect and the page has a dark background

### Requirement: Backend foundation
The backend SHALL be a Go module with a `cmd/server` entrypoint that listens on the port from the `PORT` environment variable (default 8080), and SHALL build and vet cleanly using only the standard library.

#### Scenario: Backend checks pass on the bare foundation
- **WHEN** `gofmt -l .`, `go vet ./...`, `go test ./...`, and `go build ./...` are run in `backend/`
- **THEN** `gofmt` prints nothing and every other command exits successfully

#### Scenario: Server starts
- **WHEN** `go run ./cmd/server` is executed
- **THEN** the process listens on port 8080 (or `PORT`) without panicking

### Requirement: Local development integration
The Vite dev server SHALL proxy requests under `/api` to the backend at `http://localhost:8080`, so that the browser always uses same-origin relative URLs and the backend requires no CORS configuration.

#### Scenario: Proxy configured
- **WHEN** the frontend dev server receives a request for `/api/anything`
- **THEN** it forwards the request to `http://localhost:8080/api/anything`

### Requirement: Foundation contains no product behavior
This change SHALL NOT implement calculator evaluation, API routes, or calculator UI.

#### Scenario: Placeholder only
- **WHEN** the foundation is complete
- **THEN** the UI is a placeholder and the server exposes no calculator routes
