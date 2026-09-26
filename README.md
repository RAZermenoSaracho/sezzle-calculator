# Sezzle Calculator

A small full-stack calculator. The user types an expression, and the result updates live. All parsing and arithmetic happen in a Go backend; the React frontend only sends the text and displays the answer.

This is a take-home challenge for the Sezzle Software Engineer II (LATAM) hiring process.

## Architecture

```text
Browser (React + TypeScript + Vite + Tailwind)
   │  POST /api/calculate  {"expression": "1 + 2 * 3"}
   ▼
Go server
   ├── internal/api          HTTP handlers, JSON, status codes (no math)
   └── internal/calculator   lexer + recursive-descent evaluator (no HTTP)
```

```text
backend/
  cmd/server/           process entrypoint (PORT, STATIC_DIR)
  internal/calculator/  expression engine, table-driven tests, fuzz test
  internal/api/         REST handlers and tests
frontend/
  src/components/       Calculator component
  src/calculate.ts      API client
openspec/               change proposals, specs, and design decisions
```

In production, one Go process serves both `/api` and the built frontend on one port, so no CORS or reverse proxy is needed. In development, the Vite dev server proxies `/api` to the backend.

## Supported operations

| Feature | Syntax | Example |
|---|---|---|
| Add, subtract | `+` `-` | `1 + 2 - 4` |
| Multiply, divide | `*` `/` | `6 / 3 * 2` |
| Exponent | `^` | `2 ^ 3` |
| Square root | `sqrt(x)` | `sqrt(16)` |
| Percentage (postfix) | `x%` | `50%` is `0.5` |
| Parentheses | `( )` | `(1 + 2) * 3` |
| Unary sign | `+x` `-x` | `-5 + 2` |
| Decimals | `1.5`, `.5` | `0.1 + 0.2` |

Whitespace between tokens is ignored.

### Precedence and associativity (highest to lowest)

| Level | Operator | Associativity |
|---|---|---|
| 1 | `%` (postfix) | left (`50%%` is `0.005`) |
| 2 | `^` | **right** (`2 ^ 3 ^ 2` is `512`) |
| 3 | unary `+` `-` | right |
| 4 | `*` `/` | left |
| 5 | `+` `-` | left |

Consequences: `-2 ^ 2` is `-4`, `2 ^ -2` is `0.25`, `100 + 10%` is `100.1` (percent is plain division by 100, never contextual), and `1 + 2 - 4 * 3 * (5 - 1) ^ 0.5 + sqrt(16)` is `-17`.

## Prerequisites

- Go 1.27.1 or newer (see `backend/go.mod`)
- Node.js 24 and npm
- Docker (optional, only for the container)

## Setup and running

### Backend

```bash
cd backend
go run ./cmd/server        # listens on :8080 (override with PORT)
```

### Frontend

```bash
cd frontend
npm ci
npm run dev                # http://localhost:5173, proxies /api to :8080
```

### Full application (development)

Run the two commands above in separate terminals and open http://localhost:5173.

### Full application (single process)

```bash
cd frontend && npm ci && npm run build && cd ..
cd backend && STATIC_DIR=../frontend/dist go run ./cmd/server
# http://localhost:8080
```

## Tests

```bash
# Backend (formatting check, vet, tests, build)
cd backend
test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go build ./...

# Frontend (lint, tests, build)
cd frontend
npm run lint && npm test -- --run && npm run build
```

Coverage reports (no thresholds are enforced):

```bash
cd backend  && go test -cover ./...
cd frontend && npm test -- --run --coverage
```

Optional fuzzing of the evaluator: `cd backend && go test -run='^$' -fuzz=FuzzEvaluate -fuzztime=10s ./internal/calculator`.

## API

### `POST /api/calculate`

Request (`Content-Type: application/json`):

```json
{ "expression": "1 + 2 * 3" }
```

Success (`200`):

```json
{ "result": 7 }
```

Failure (`400` and others): `{ "error": "<message>" }`

### `GET /api/health`

Returns `200` with `{"status":"ok"}`.

### curl examples

```bash
curl -s -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' -d '{"expression":"1 + 2 * 3"}'
# {"result":7}

curl -s -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"1 + 2 - 4 * 3 * (5 - 1) ^ 0.5 + sqrt(16)"}'
# {"result":-17}

curl -s -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' -d '{"expression":"1 / 0"}'
# {"error":"division by zero"}

curl -s localhost:8080/api/health
# {"status":"ok"}
```

### Error behavior

Evaluation errors are client input errors, so they all return `400` with a stable message. Responses never include panic text, stack traces, or parser internals.

| Status | `error` | Cause |
|---|---|---|
| 400 | `empty expression` | empty or whitespace-only input |
| 400 | `invalid expression` | syntax error, unsupported character, missing parenthesis |
| 400 | `division by zero` | `x / 0`, `0 ^ -1` |
| 400 | `invalid square root` | `sqrt` of a negative number |
| 400 | `invalid number` | overflow or NaN, such as `10 ^ 1000` or `(-8) ^ 0.5` |
| 400 | `expression too complex` | over 1000 characters or nested over 100 levels |
| 400 | `invalid request` | malformed JSON, missing or non-string `expression`, unknown fields, trailing data |
| 405 | `method not allowed` | wrong method (with an `Allow` header) |
| 413 | `request too large` | body over 4 KiB |
| 415 | `unsupported media type` | content type is not `application/json` |
| 500 | `internal error` | unexpected failure |

The UI shows the `error` text as returned, or `Could not reach the server` if the request fails.

## Design decisions

- **Backend is the source of truth.** The frontend never evaluates expressions, so there is one set of semantics.
- **Hand-written lexer and recursive-descent parser** that evaluates as it parses. The grammar maps one-to-one to the precedence table, which keeps it easy to review. No `eval`, no third-party evaluator.
- **`float64` arithmetic.** Results such as `0.1 + 0.2` carry binary rounding, so the UI displays 12 significant digits. The API returns the raw value.
- **Postfix percent only.** `50%` is `0.5`; `100 + 10%` is `100.1`, not `110`.
- **Controlled failures.** Every failure mode has a sentinel error mapped to a stable message. A recover wrapper is defense in depth. Input length and nesting depth are limited so hostile input cannot exhaust the stack.
- **Live UI without races.** A 300 ms debounce reduces requests, each keystroke aborts the previous request, and late replies are ignored.
- **One endpoint.** A single expression endpoint replaces per-operator endpoints.
- **Same-origin deployment.** The Vite proxy in development and static serving in production mean no CORS.
- **OpenSpec.** Each piece of work was planned and specified in `openspec/` before implementation; archived changes hold the history.

## Assumptions

- `^` is right-associative and binds tighter than unary minus (`-2 ^ 2 = -4`).
- No implicit multiplication (`2(3)` is invalid), no scientific notation (`1e3` is invalid), and `sqrt` is the only function.
- Numbers must be `12`, `1.5`, or `.5`; `5.` is invalid.
- Double-precision floating point is sufficient; there is no arbitrary-precision mode.
- The frontend supports current evergreen browsers.

## Docker

One multi-stage `Dockerfile` builds the frontend and a static Go binary, then runs both in a small non-root distroless image on port 8080.

```bash
docker build -t sezzle-calculator .
docker run --rm -p 8080:8080 sezzle-calculator
# open http://localhost:8080, or: curl localhost:8080/api/health
```

Environment variables: `PORT` (default `8080`) and `STATIC_DIR` (set to `/app/static` in the image; unset means API only).

## CI/CD

- **CI** (`.github/workflows/ci.yml`) runs on pull requests and pushes to `main`, on GitHub-hosted runners with no secrets:
  - frontend: `npm ci`, lint, tests, build
  - backend: `gofmt` check, `go vet`, tests, build
  - docker: builds the image (not published)
- **CD** (`.github/workflows/deploy.yml`) runs on the repository's self-hosted runner after CI succeeds on `main`. It builds the commit inside the shared Docker VM through Vagrant, replaces the `sezzle-calculator` container, and checks `/api/health`. It is disabled until the repository variable `DEPLOY_ENABLED` is `true`.

Runner, GitHub Environment, variables, VM, and Cloudflare setup are operator tasks; see [`deploy/README.md`](deploy/README.md).

## AI tooling disclosure

This project was developed with **Claude Code** (Anthropic's CLI agent) using the **OpenSpec** spec-driven workflow: changes were proposed, reviewed, implemented, manually verified by the author, and only then archived and committed. Claude wrote most of the code, tests, and documentation under the author's direction; the author reviewed and approved every change.

The prompts used to drive the work are recorded in [`PROMPTS.md`](PROMPTS.md), and the project rules given to the agent are in [`CLAUDE.md`](CLAUDE.md).
