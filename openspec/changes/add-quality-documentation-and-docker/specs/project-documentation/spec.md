## ADDED Requirements

### Requirement: Reviewer-oriented README
`README.md` SHALL let a reviewer with no prior context understand, run, test, and evaluate the project, and SHALL include: project overview, architecture, supported operations with precedence and associativity, prerequisites, frontend setup, backend setup, running the full application, running tests (including coverage), API contract, curl examples, error behavior, design decisions, assumptions, Docker instructions, CI/CD overview, AI tooling disclosure, and a reference to the prompts in `PROMPTS.md`.

#### Scenario: Required sections present
- **WHEN** a reviewer reads `README.md`
- **THEN** every section listed above is present and concise

#### Scenario: Commands work as written
- **WHEN** the documented setup, run, test, and curl commands are executed
- **THEN** they succeed as described on a machine with the stated prerequisites

### Requirement: Documentation matches behavior
The README's syntax, precedence, status codes, and error messages SHALL match the specs of `implement-go-expression-engine` and `implement-calculator-api`.

#### Scenario: Consistent error table
- **WHEN** comparing the README error table to the API behavior
- **THEN** statuses and `error` strings are identical

### Requirement: Coverage reporting
The repository SHALL document commands producing backend and frontend coverage reports, without enforcing numeric thresholds.

#### Scenario: Coverage commands
- **WHEN** `go test -cover ./...` and `npm test -- --run --coverage` are executed
- **THEN** each prints a coverage summary and exits successfully

### Requirement: AI disclosure
The README SHALL disclose the AI tooling used during development and point to `PROMPTS.md`.

#### Scenario: Disclosure present
- **WHEN** a reviewer looks for AI usage information
- **THEN** the README names the tools and references `PROMPTS.md`
