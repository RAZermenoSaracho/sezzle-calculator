# Proposal: Implement Go Expression Engine

## Why

Correct calculation is the top project priority. The evaluator must exist, be exhaustively tested, and be independent of HTTP before any transport or UI depends on it.

## What Changes

- Add package `backend/internal/calculator` exposing `Evaluate(expression string) (float64, error)`.
- Hand-written lexer plus recursive-descent parser that evaluates directly (no AST, no `eval`, no third-party library).
- Explicit grammar, precedence, and associativity (see design).
- Sentinel errors that later layers map to stable messages.
- Table-driven unit tests covering every semantic below.
- **The package MUST NOT import `net/http` or any JSON/transport concern.**

## Capabilities

### New Capabilities
- `expression-evaluation`: parsing and evaluating arithmetic expression strings with defined semantics and errors.

### Modified Capabilities
None.

## Non-goals

Variables, functions other than `sqrt`, implicit multiplication (`2(3)`), scientific notation, contextual percentage (`100 + 10%` is `100.1`, never `110`), arbitrary-precision arithmetic, an AST or public tokenizer API.

## Impact

- New: `backend/internal/calculator/{lexer,parser,errors}.go` (or a single file if small) and `*_test.go`.
- Depends on: `bootstrap-fullstack-foundation`. Consumed by: `implement-calculator-api`.
