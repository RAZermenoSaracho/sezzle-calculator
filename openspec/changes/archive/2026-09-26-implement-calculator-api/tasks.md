# Tasks

## 1. Handler
- [x] 1.1 Create `internal/api` with request/response types, `NewHandler()`, and the injectable evaluator field. Verify `go build ./...`.
- [x] 1.2 Implement `POST /api/calculate` with strict decoding, size limit, content-type check, and sentinel-error mapping per `design.md`. Verify with the spec scenarios via tests.
- [x] 1.3 Implement `GET /api/health`, 405 handling with `Allow`, and the recover wrapper. Verify with tests.
- [x] 1.4 Replace the empty mux in `cmd/server/main.go` with `api.NewHandler()`. Verify `go run ./cmd/server` then `curl -s -XPOST localhost:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"1 + 2 * 3"}'` prints `{"result":7}`.

## 2. Tests
- [x] 2.1 Table-driven `httptest` tests for success, each evaluation error, malformed JSON, missing/non-string field, unknown field, trailing data, wrong method, wrong content type, oversized body. Verify `go test ./internal/api`.
- [x] 2.2 Test the 500 path with an injected failing and a panicking evaluator. Verify no internal text in the body.

## 3. Quality gate
- [x] 3.1 Run `gofmt -l .`, `go vet ./...`, `go test ./...`, `go build ./...` in `backend/`. Verify all succeed and no calculation logic exists in `internal/api`.
- [x] 3.2 Run the curl examples for success, syntax error, and division by zero against the live server; record them for the README change. Verify status codes match the table.
