# Proposal: Quality, Documentation, and Docker

## Why

Reviewers must be able to understand, run, and verify the project without prior context, and CD needs a single deployable artifact. The working product exists at this point, so packaging and documentation can describe reality.

## What Changes

- Write the final `README.md` with every section required by `CLAUDE.md` (overview, architecture, supported syntax with precedence table, prerequisites, setup, run, tests, API contract, curl examples, error behavior, design decisions, assumptions, Docker, CI/CD overview, AI tooling disclosure, prompts reference to `PROMPTS.md`).
- Add coverage commands (`go test -cover`, `vitest --coverage` with `@vitest/coverage-v8`) and document them; no coverage thresholds.
- Add a single multi-stage `Dockerfile`: build frontend, build Go binary, run both from one small image. The Go server serves the built frontend and `/api` on one port.
- Small backend change: `cmd/server` serves static files from `STATIC_DIR` when set (unset in local dev), plus tests for that behavior.
- Add `.dockerignore`.

## Capabilities

### New Capabilities
- `project-documentation`: reviewer-facing README content requirements.
- `container-packaging`: single-image production packaging and same-origin serving.

### Modified Capabilities
None.

## Non-goals

docker-compose, reverse proxy, nginx, Kubernetes, registries, deployment scripts (owned by `add-github-ci-cd`), any server mutation. Running the container on Ricardo's infrastructure requires his approval and is out of scope.

## Impact

- New: `README.md` (replaces stub), `Dockerfile`, `.dockerignore`; modified: `backend/cmd/server/main.go`, `frontend/package.json` (coverage dep).
- Depends on: `build-live-calculator-ui`. Consumed by: `add-github-ci-cd`.
