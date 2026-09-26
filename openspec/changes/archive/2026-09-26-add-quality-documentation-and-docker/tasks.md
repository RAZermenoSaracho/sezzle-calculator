# Tasks

## 1. Static serving
- [x] 1.1 Add `STATIC_DIR` handling in `cmd/server/main.go` (mux composition kept testable, e.g. `newHandler(staticDir string)`). Verify with `httptest` tests for enabled, disabled, and API-unaffected cases.
- [x] 1.2 Run `gofmt -l .`, `go vet ./...`, `go test ./...` in `backend/`. Verify success.

## 2. Docker
- [x] 2.1 Write the multi-stage `Dockerfile` and `.dockerignore` per `design.md`. Verify by review that Node/Go versions match `package.json`/`go.mod`, the runtime is non-root, and `STATIC_DIR`/`PORT` are set.
- [x] 2.2 If Docker is available locally, `docker build` and `docker run` and `curl` `/`, `/api/health`, and `/api/calculate`; otherwise record what was validated statically (`npm run build`, `go build`, tests) and what needs a Docker-capable environment. Do not touch the server or VM.

## 3. Coverage
- [x] 3.1 Add `@vitest/coverage-v8`; verify `npm test -- --run --coverage` prints a summary and `go test -cover ./...` prints package coverage.

## 4. README
- [x] 4.1 Write `README.md` with all `CLAUDE.md`-required sections, sourcing syntax and error tables from the change 2 and 3 specs. Verify each required section exists (checklist against `CLAUDE.md`).
- [x] 4.2 Execute every documented local command (setup, run, tests, curl). Verify each works as written.

## 5. Quality gate
- [x] 5.1 Run all backend and frontend checks from `CLAUDE.md`. Verify success and that no server-side or CI/CD files were added.

## Verification record

Docker is not available on the development Mac, so task 2.2 was first closed with static validation only (Dockerfile and `.dockerignore` review, `npm run build`, static `CGO_ENABLED=0 go build`, and the Go server serving the built `dist` on one port). The container behavior was later verified manually by Ricardo on a Docker-capable Ubuntu environment (Vagrant VM):

- `docker build -t sezzle-calculator:test .` succeeded; image about 16.2 MB on disk (3.69 MB content).
- `docker run -d --name sezzle-calculator -p 8080:8080 sezzle-calculator:test` started; the server listened on 8080.
- `GET /api/health` returned 200 `{"status":"ok"}`.
- `POST /api/calculate` with `1 + 2 - 4 * 3 * (5 - 1) ^ 0.5 + sqrt(16)` returned 200 `{"result":-17}`.
- `GET /` returned 200 with the built frontend HTML.
- The container ran as `nonroot:nonroot` using about 2.5 MiB of memory.
- Port 8080 was forwarded from the VM to the macOS host, and the app loaded through the existing Cloudflare Tunnel with the frontend reaching the backend on the same origin.

Out of scope, recorded for a separate change: an expression such as `54+++++6666` is accepted (repeated unary plus) and evaluates to `6720`. It is unchanged here.
