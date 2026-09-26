# Design: Documentation and Docker

## Topology

One container, one port:

```text
browser ──► :8080 Go server ──► /api/*      handlers (change 3)
                          └──► /*          static frontend from STATIC_DIR
```

Same origin, so no CORS and no proxy. The frontend's relative `/api` URLs work unchanged.

## Dockerfile (multi-stage)

1. `node:24-alpine`: `npm ci`, `npm run build` in `frontend/`.
2. `golang:<version from go.mod>-alpine`: `CGO_ENABLED=0 go build -o /server ./cmd/server`.
3. Runtime: `gcr.io/distroless/static:nonroot` (or `alpine` if a shell/wget healthcheck is preferred); copy `/server` and `frontend/dist` to `/app/static`; `ENV PORT=8080 STATIC_DIR=/app/static`; `EXPOSE 8080`; run as non-root.

Health: `GET /api/health`. With distroless there is no shell; health verification is done externally (`curl`), and the Dockerfile defines no `HEALTHCHECK`. If the operator needs one, the runtime base can change to alpine; decide at apply time and record in README.

## Static serving

In `main.go`: if `STATIC_DIR` is non-empty, mount `http.FileServer(http.Dir(dir))` at `/` alongside the API handler on the same mux. `/api/` continues to belong to the API. No SPA fallback is needed (one page, no client routing). Test with `httptest` and a temp directory: `/` returns `index.html`; `/api/health` still returns JSON; with `STATIC_DIR` unset `/` is 404.

## README plan

Concise, ordered exactly as `CLAUDE.md` lists. Content sources are the specs of changes 2 and 3 (grammar/precedence table, error messages/status codes) so documentation does not drift; curl examples are copied from the verified commands in change 3. Design decisions to state: backend as source of truth, recursive descent, `float64`, postfix percent, 400 for evaluation errors, proxy vs CORS, single-image deployment. Assumptions: `^` right-associative, `-2^2 = -4`, no implicit multiplication, float precision limits. AI disclosure: Claude Code with OpenSpec, prompts in `PROMPTS.md`.

## Verification honesty

Docker is not installed on the development macOS host. The apply phase must state which parts were validated statically (Dockerfile syntax review, `.dockerignore`, `go build` and `npm run build` succeeding, static-serving tests) versus which require a Docker-capable environment (`docker build`, `docker run`, `curl` against the container). Executing the image in `~/scripts/docker-test-vm` counts as server-side work and needs Ricardo's approval.
