# Tasks

## 1. Engine
- [x] 1.1 Create `internal/calculator` with `Evaluate` and the sentinel errors. Verify `go build ./...` and that the package does not import `net/http` (`go list -deps ./internal/calculator | grep net/http` prints nothing).
- [x] 1.2 Implement the lexer (numbers, operators, parentheses, `sqrt`, whitespace, unsupported characters). Verify with lexer-level table cases in tests.
- [x] 1.3 Implement the recursive-descent parser exactly per `design.md` grammar, including depth and length limits. Verify against the spec scenarios.
- [x] 1.4 Normalize `-0`, reject NaN/Inf, map division by zero and negative sqrt. Verify with error-case table tests.

## 2. Tests
- [x] 2.1 Table-driven tests for precedence, associativity, parentheses, unary, decimals, whitespace, exponent, sqrt, percent, and the reference expression. Verify `go test ./internal/calculator` passes.
- [x] 2.2 Table-driven error tests for empty, malformed, invalid characters, division by zero, invalid sqrt, non-finite, too complex, using `errors.Is`. Verify pass.
- [x] 2.3 Add a `FuzzEvaluate` or randomized no-panic test with a short seed corpus. Verify `go test` passes and `go test -run=^$ -fuzz=FuzzEvaluate -fuzztime=5s ./internal/calculator` finds no panic.

## 3. Quality gate
- [x] 3.1 Run `gofmt -l .`, `go vet ./...`, `go test ./...`, `go build ./...` in `backend/`. Verify all succeed and no HTTP code was added.
