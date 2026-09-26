## ADDED Requirements

### Requirement: CI workflow triggers
A GitHub Actions workflow SHALL run on pull requests and on pushes to `main`, with read-only repository permissions.

#### Scenario: Pull request
- **WHEN** a pull request is opened or updated
- **THEN** the CI workflow runs

### Requirement: Frontend verification
CI SHALL install frontend dependencies with `npm ci`, then run lint, tests (`npm test -- --run`), and build, failing on any error.

#### Scenario: Frontend failure
- **WHEN** a frontend test fails
- **THEN** the frontend job fails and the workflow is red

### Requirement: Backend verification
CI SHALL check Go formatting with `gofmt -l` (failing when it prints anything), then run `go vet ./...`, `go test ./...`, and `go build ./...`, using the Go version declared in `backend/go.mod`.

#### Scenario: Unformatted code
- **WHEN** a Go file is not gofmt-formatted
- **THEN** the backend job fails at the formatting check

### Requirement: Container build check
CI SHALL build the Docker image without publishing it.

#### Scenario: Broken Dockerfile
- **WHEN** the Dockerfile fails to build
- **THEN** the docker job fails

### Requirement: Deterministic and minimal
The workflow SHALL use only first-party `actions/*` actions pinned to major versions, no secrets, and no matrix builds.

#### Scenario: No secrets required
- **WHEN** CI runs on a fork pull request
- **THEN** it completes without needing any repository secret
