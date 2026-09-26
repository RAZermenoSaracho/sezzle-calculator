# Proposal: Bootstrap Full-Stack Foundation

## Why

The repository contains an untouched Vite scaffold and an empty Go module. Every later change needs an agreed layout, working toolchains, and runnable check commands (lint, test, build, vet) before feature code lands. Establishing this once keeps the later changes small and reviewable.

## What Changes

- Fix the repository layout: `backend/`, `frontend/`, `openspec/`, `.github/workflows/` (created later by `add-github-ci-cd`).
- Frontend: clean the Vite React + TypeScript template down to a minimal shell, integrate Tailwind CSS, set a permanent dark base, and add Vitest + React Testing Library with a `test` script.
- Backend: keep the Go module `sezzle-calculator/backend`, add `cmd/server` entrypoint that starts an HTTP server (empty mux, `PORT` env, default 8080) with graceful failure logging, and reserve `internal/` for later packages.
- Vite dev server proxies `/api` to `http://localhost:8080`, so no CORS is ever needed.
- Document the canonical local commands (single source used by later changes and CI).
- Extend `.gitignore` only where needed.
- **No calculator behavior, no API routes, no UI beyond a placeholder heading.**

## Capabilities

### New Capabilities
- `project-foundation`: repository layout, toolchains, dev commands, and the frontend/backend integration contract for local development.

### Modified Capabilities
None.

## Non-goals

No monorepo tooling (no Nx/Turborepo/workspaces), no Makefile, no Prettier/husky, no Docker, no CI, no state or routing libraries.

## Impact

- Touches `frontend/` (config, `src/`, `package.json`), `backend/` (`cmd/server/main.go`), `.gitignore`.
- Adds dev dependencies: `tailwindcss`, `@tailwindcss/vite`, `vitest`, `jsdom`, `@testing-library/react`, `@testing-library/jest-dom`, `@testing-library/user-event`.
- Dependency for all later changes.
