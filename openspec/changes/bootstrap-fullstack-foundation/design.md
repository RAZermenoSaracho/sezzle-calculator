# Design: Bootstrap Full-Stack Foundation

## Layout

```text
backend/
  go.mod                  module sezzle-calculator/backend
  cmd/server/main.go      process entrypoint only (env, listen, log)
  internal/               populated by later changes:
    calculator/           (change 2) expression engine, no net/http import
    api/                  (change 3) HTTP handlers and JSON types
frontend/
  src/                    App.tsx, main.tsx, index.css (Tailwind import)
  vite.config.ts          react + tailwind plugins, /api proxy, vitest config
```

`internal/` packages are created by the change that owns them; this change creates only `cmd/server`.

## Decisions

- **Tailwind v4 via `@tailwindcss/vite`**: no `tailwind.config.js`, no PostCSS config; `index.css` is `@import "tailwindcss";`. This is the least configuration. Template CSS (`App.css`, demo assets) is deleted.
- **Dark only**: `index.html` sets a dark background class on `<body>`; there is no theme switch and no `prefers-color-scheme` branching.
- **Vitest** (same Vite config, jsdom environment, `globals` off, jest-dom setup file `src/test/setup.ts`). A single smoke test renders the placeholder to prove the pipeline; change 4 replaces it.
- **Proxy instead of CORS**: `server.proxy['/api'] -> http://localhost:8080`. The backend therefore never sets CORS headers. In production the same origin serves both (change 5).
- **Backend entrypoint**: `main` reads `PORT` (default `8080`), builds an `http.Server` around a mux passed in from a `newHandler()` function (returns an empty mux for now), sets `ReadHeaderTimeout`, and logs fatal errors. Change 3 replaces the mux contents; it does not restructure `main`.
- **Go version**: `go.mod` currently declares `go 1.27.1` matching the installed toolchain; leave as is and let CI use `go-version-file`.
- **No Makefile**: commands are short enough to list in the README; adding a Makefile is extra surface.

## Canonical commands

| Purpose | Command |
|---|---|
| Frontend dev | `cd frontend && npm run dev` (http://localhost:5173) |
| Frontend check | `npm run lint`, `npm test -- --run`, `npm run build` |
| Backend run | `cd backend && go run ./cmd/server` (http://localhost:8080) |
| Backend check | `gofmt -l .` (must print nothing), `go vet ./...`, `go test ./...`, `go build ./...` |

CI (change 6) and the README (change 5) reference exactly these.
