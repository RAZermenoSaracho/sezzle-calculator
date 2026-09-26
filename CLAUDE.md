# Sezzle Calculator Challenge

## Project Context

This repository is a take-home challenge for the Sezzle Software Engineer II (LATAM) hiring process.

The goal is to build a small, production-quality full-stack calculator while keeping the implementation intentionally simple, readable, testable, and easy to review.

The assignment should demonstrate sound engineering judgment rather than maximize features.

## Core Requirements

The application consists of:

- React + TypeScript frontend
- Vite
- Tailwind CSS
- Go backend
- REST API
- Unit tests for frontend and backend
- Documentation
- Docker support
- GitHub Actions CI
- GitHub Actions CD to a self-hosted server

Supported calculator functionality:

- Addition
- Subtraction
- Multiplication
- Division
- Parentheses
- Unary positive and negative numbers
- Exponentiation
- Square root
- Percentage

The calculator must respect standard mathematical precedence.

Examples:

```text
1 + 2
1 + 2 * 3
(1 + 2) * 3
-5 + 2
2 ^ 3
sqrt(16)
50%
1 + 2 - 4 * 3 * (5 - 1) ^ 0.5 + sqrt(16)
```

Invalid expressions must return a controlled error rather than panic or produce undefined behavior.

## Engineering Priorities

Prioritize, in this order:

1. Correctness
2. Simplicity
3. Readability
4. Testability
5. Maintainability
6. User experience
7. Additional features

Do not over-engineer this assignment.

Avoid unnecessary abstractions, frameworks, dependencies, patterns, configuration, and infrastructure.

Prefer the Go standard library when practical.

## Language Rules

All source code, identifiers, comments, documentation, commit messages, API fields, test descriptions, and developer-facing text must be written in English.

Do not add comments that merely repeat what the code already communicates.

Comments should explain non-obvious decisions, constraints, algorithms, or edge cases.

## Repository Structure

Keep frontend and backend clearly separated.

Target structure:

```text
.
├── backend/
├── frontend/
├── openspec/
├── .github/
│   └── workflows/
├── CLAUDE.md
├── PROMPTS.md
└── README.md
```

Do not introduce a monorepo framework.

## Frontend

Technology:

- React
- TypeScript
- Vite
- Tailwind CSS

Keep the frontend deliberately small.

### UI

The application uses one permanent dark theme. Do not implement light mode or a theme switcher.

The calculator should be centered and minimal.

The primary interface consists of:

1. Expression input
2. Thin horizontal divider
3. Result or error output

The input should not look like a traditional four-sided form field. Avoid a visible box or border around all four sides.

The horizontal divider visually separates the expression from its result.

Example:

```text
1 + 2 * 3
────────────
7
```

Results update automatically while the user types. The user must not need to press Enter or click a Calculate button.

Use a short debounce where appropriate to avoid unnecessary requests while preserving the feeling of immediate feedback.

The UI must be usable on desktop and mobile screen sizes.

Avoid unnecessary visual elements, cards, gradients, animations, menus, navigation, or decorative components.

Use Tailwind utilities instead of custom CSS whenever practical.

### Frontend Responsibilities

The frontend:

- captures the expression
- sends it to the backend
- displays the result
- displays controlled errors
- handles loading/request lifecycle safely
- avoids showing stale responses when requests complete out of order

The frontend must not duplicate the backend expression evaluator.

The backend is the source of truth for calculation semantics.

## Backend

Use Go.

Keep HTTP transport logic separate from calculator/evaluation logic.

A suggested conceptual separation is:

```text
HTTP request
    ↓
handler
    ↓
calculator/evaluator
    ↓
parser
    ↓
result
```

The exact package layout may be adjusted if a simpler idiomatic structure is preferable.

### API

Prefer one expression-oriented endpoint rather than separate endpoints for every arithmetic operator.

Suggested contract:

```http
POST /api/calculate
Content-Type: application/json
```

Request:

```json
{
  "expression": "1 + 2 * 3"
}
```

Successful response:

```json
{
  "result": 7
}
```

Invalid expression:

```json
{
  "error": "invalid expression"
}
```

The API must use appropriate HTTP status codes and JSON content types.

Do not expose Go panic messages, parser internals, stack traces, or implementation details to clients.

A lightweight health endpoint may be added if useful for Docker/deployment verification.

## Expression Evaluator

Do not use `eval`, JavaScript execution, shell execution, or external expression-evaluation services.

Implement the expression evaluation in Go.

The parser must correctly handle mathematical precedence rather than evaluating strictly from left to right.

At minimum, support:

```text
()
unary + and -
sqrt(...)
%
^
* /
+ -
```

Define and document precedence and associativity explicitly.

Exponentiation should behave as exponentiation rather than multiplication.

`sqrt(x)` must reject values outside the supported real-number domain.

Division by zero must return a controlled error.

Malformed expressions must return a controlled syntax error.

Whitespace between tokens should be ignored.

Support decimal numbers.

Percentage should use postfix semantics:

```text
50% = 0.5
```

Do not introduce contextual calculator-style percentage behavior such as interpreting `100 + 10%` as `110` unless an approved OpenSpec change explicitly defines it.

Internally, prefer a small lexer/tokenizer plus a precedence-aware parser such as recursive descent or Pratt parsing.

Choose the simplest implementation that remains clear and testable.

## Error Handling

Expected errors include:

- empty expression
- malformed syntax
- unsupported characters
- unexpected tokens
- missing parentheses
- division by zero
- invalid square root
- invalid numeric values

Errors presented through the API should be stable and understandable.

Do not panic for user input.

## Testing

Both layers require automated tests.

### Backend

Tests should cover the evaluator independently from HTTP.

Include representative tests for:

- operator precedence
- parentheses
- unary operators
- decimals
- exponentiation
- square root
- percentage
- whitespace
- division by zero
- malformed expressions
- invalid characters

HTTP handler tests should cover the API contract and important error cases.

Prefer table-driven Go tests.

### Frontend

Use a lightweight standard React/Vite testing setup.

Tests should focus on meaningful behavior rather than implementation details.

Cover at least:

- expression entry
- successful calculation
- error display
- live/debounced calculation behavior

Do not pursue arbitrary coverage percentages by testing trivial implementation details.

## Formatting and Quality

Before considering implementation complete:

Backend:

```bash
go fmt ./...
go vet ./...
go test ./...
```

Frontend:

```bash
npm run lint
npm test -- --run
npm run build
```

Use the actual scripts configured by the repository if they differ.

There must be no knowingly ignored compiler, linter, or test failures.

## Docker

Provide Docker support for the completed application.

Keep the runtime architecture simple.

Prefer a production arrangement where the built frontend and backend can be deployed together without requiring separate public services unless there is a strong reason otherwise.

Do not add Kubernetes, Terraform, Docker Swarm, or similar infrastructure.

## CI

GitHub Actions CI should run on pull requests and relevant pushes.

CI should verify at minimum:

- frontend install
- frontend lint
- frontend tests
- frontend build
- backend formatting check
- backend vet
- backend tests
- backend build

Keep the workflow understandable and deterministic.

## CD and Server Constraints

Deployment targets Ricardo's existing Mac server.

Important constraints:

- Do not assume Docker exists on the macOS host.
- Do not install Docker Desktop or Colima.
- Container workloads run inside Vagrant VMs.
- The existing shared Docker VM is intended for ephemeral testing and lives at `~/scripts/docker-test-vm`.
- Do not modify or reuse the production n8n container without explicit operator approval.
- The server already uses repository-specific GitHub Actions self-hosted runners.
- Server deployment work must be treated separately from local application development.

Any operation requiring changes to the server, GitHub runner configuration, Vagrant configuration, Cloudflare configuration, DNS, or production processes requires explicit operator approval before execution.

During local development, design the deployment artifacts but do not mutate the server.

## Secrets

Never commit:

- credentials
- API tokens
- GitHub tokens
- SSH private keys
- Cloudflare tokens
- `.env` files containing secrets

Use environment variables or GitHub Actions secrets where credentials are eventually required.

## README

The final README must be useful to a reviewer who has never seen the repository.

It must include:

- project overview
- architecture
- supported operations
- local prerequisites
- frontend setup
- backend setup
- how to run the full application
- how to run tests
- API contract
- curl examples
- error behavior
- design decisions
- assumptions
- Docker instructions
- CI/CD overview
- AI tooling disclosure
- prompts used during development

Keep it concise enough for a take-home assignment.

## OpenSpec Workflow

This repository uses OpenSpec for planned changes.

Do not implement substantial functionality outside an approved OpenSpec change.

The workflow is:

```text
proposal
→ review
→ /opsx:apply
→ automated verification
→ operator/manual verification
→ closure
→ archive
→ commit
→ push
→ next change
```

Only one implementation change should normally be active at a time.

### Proposal Phase

Proposal work may:

- inspect the repository
- research existing code
- create OpenSpec artifacts
- identify dependencies
- define acceptance criteria

Proposal work must not implement the application.

When creating multiple proposals, independent repository inspection may be delegated to subagents to reduce elapsed time, but the final proposals must be mutually consistent.

### Apply Phase

When `/opsx:apply` is requested:

1. Read the complete change.
2. Inspect the current repository state.
3. Implement only that change.
4. Run its relevant automated verification.
5. Update OpenSpec task state accurately.
6. Stop before archive, commit, or push.
7. Report:
   - files changed
   - important design decisions
   - tests/checks run
   - results
   - remaining manual verification

Do not archive, commit, or push automatically after implementation.

### Manual Verification Gate

Ricardo is the operator.

After an apply phase, wait for explicit confirmation that the change is accepted.

Do not interpret successful automated tests as permission to close the change.

### Closure

Only after explicit operator approval:

1. Recheck repository state.
2. Run appropriate final verification.
3. Mark the approved change complete.
4. Archive it according to OpenSpec conventions.
5. Commit only the files belonging to the approved change.
6. Use a concise conventional commit message.
7. Push the current branch.
8. Report the resulting commit and repository state.

Never silently include unrelated modifications.

## Git Discipline

Before modifying files, inspect:

```bash
git status
git diff
```

Do not overwrite unrelated user work.

Do not use destructive Git commands to clean unrelated changes.

Keep commits scoped to one approved OpenSpec change whenever practical.

Do not commit or push until explicitly authorized through the closure workflow.

## Scope Control

This is a hiring challenge, not a long-term SaaS product.

Do not add features such as:

- authentication
- user accounts
- databases
- calculation history persistence
- analytics
- localization
- themes
- scientific calculator functions beyond the defined scope
- WebSockets
- GraphQL
- Redux or another global state framework
- microservice orchestration
- Kubernetes
- observability platforms

If a proposed implementation significantly increases complexity without directly improving the assignment requirements, prefer the simpler design.