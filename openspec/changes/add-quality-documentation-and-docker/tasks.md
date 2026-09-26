# Tasks

## 1. Static serving
- [ ] 1.1 Add `STATIC_DIR` handling in `cmd/server/main.go` (mux composition kept testable, e.g. `newHandler(staticDir string)`). Verify with `httptest` tests for enabled, disabled, and API-unaffected cases.
- [ ] 1.2 Run `gofmt -l .`, `go vet ./...`, `go test ./...` in `backend/`. Verify success.

## 2. Docker
- [ ] 2.1 Write the multi-stage `Dockerfile` and `.dockerignore` per `design.md`. Verify by review that Node/Go versions match `package.json`/`go.mod`, the runtime is non-root, and `STATIC_DIR`/`PORT` are set.
- [ ] 2.2 If Docker is available locally, `docker build` and `docker run` and `curl` `/`, `/api/health`, and `/api/calculate`; otherwise record what was validated statically (`npm run build`, `go build`, tests) and what needs a Docker-capable environment. Do not touch the server or VM.

## 3. Coverage
- [ ] 3.1 Add `@vitest/coverage-v8`; verify `npm test -- --run --coverage` prints a summary and `go test -cover ./...` prints package coverage.

## 4. README
- [ ] 4.1 Write `README.md` with all `CLAUDE.md`-required sections, sourcing syntax and error tables from the change 2 and 3 specs. Verify each required section exists (checklist against `CLAUDE.md`).
- [ ] 4.2 Execute every documented local command (setup, run, tests, curl). Verify each works as written.

## 5. Quality gate
- [ ] 5.1 Run all backend and frontend checks from `CLAUDE.md`. Verify success and that no server-side or CI/CD files were added.
