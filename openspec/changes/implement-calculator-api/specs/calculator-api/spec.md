## ADDED Requirements

### Requirement: Calculate endpoint
The backend SHALL expose `POST /api/calculate` accepting `{"expression": string}` and returning `{"result": number}` with status 200 when the expression evaluates successfully, delegating all evaluation to `internal/calculator`.

#### Scenario: Successful calculation
- **WHEN** a client posts `{"expression": "1 + 2 * 3"}`
- **THEN** the response is 200 with `Content-Type: application/json` and body `{"result":7}`

#### Scenario: Reference expression
- **WHEN** a client posts `{"expression": "sqrt(16) + 50%"}`
- **THEN** the response is 200 with `{"result":4.5}`

### Requirement: Evaluation error responses
The endpoint SHALL respond with status 400 and body `{"error": <message>}` for evaluation failures, using stable messages: `empty expression`, `invalid expression`, `division by zero`, `invalid square root`, `invalid number`, `expression too complex`.

#### Scenario: Syntax error
- **WHEN** a client posts `{"expression": "1 +"}`
- **THEN** the response is 400 with `{"error":"invalid expression"}`

#### Scenario: Division by zero
- **WHEN** a client posts `{"expression": "1 / 0"}`
- **THEN** the response is 400 with `{"error":"division by zero"}`

#### Scenario: Invalid square root
- **WHEN** a client posts `{"expression": "sqrt(-4)"}`
- **THEN** the response is 400 with `{"error":"invalid square root"}`

#### Scenario: Empty expression
- **WHEN** a client posts `{"expression": "  "}`
- **THEN** the response is 400 with `{"error":"empty expression"}`

### Requirement: Request validation
The endpoint SHALL reject malformed JSON, a missing or non-string `expression`, unknown fields, and trailing data with 400 `{"error":"invalid request"}`; bodies over 4 KiB with 413; non-JSON content types with 415; and non-POST methods with 405 and an `Allow: POST` header.

#### Scenario: Malformed JSON
- **WHEN** a client posts `{"expression":`
- **THEN** the response is 400 with `{"error":"invalid request"}`

#### Scenario: Missing field
- **WHEN** a client posts `{}`
- **THEN** the response is 400 with `{"error":"invalid request"}`

#### Scenario: Wrong method
- **WHEN** a client sends `GET /api/calculate`
- **THEN** the response is 405 with `Allow: POST`

#### Scenario: Oversized body
- **WHEN** a client posts a body larger than 4 KiB
- **THEN** the response is 413

### Requirement: No internal detail leakage
Responses SHALL NOT include Go panic messages, stack traces, parser positions, or wrapped internal error text; unexpected failures SHALL return 500 `{"error":"internal error"}`.

#### Scenario: Unexpected failure
- **WHEN** the evaluator returns an unrecognized error or panics
- **THEN** the response is 500 with `{"error":"internal error"}` and no other detail

### Requirement: Health endpoint
The backend SHALL expose `GET /api/health` returning 200 `{"status":"ok"}`.

#### Scenario: Health check
- **WHEN** a client sends `GET /api/health`
- **THEN** the response is 200 with `{"status":"ok"}`

### Requirement: Transport separated from domain
The `internal/api` package SHALL depend on `internal/calculator`, and `internal/calculator` SHALL NOT depend on `internal/api` or `net/http`.

#### Scenario: Dependency direction
- **WHEN** listing dependencies of `internal/calculator`
- **THEN** neither `net/http` nor `internal/api` appears
