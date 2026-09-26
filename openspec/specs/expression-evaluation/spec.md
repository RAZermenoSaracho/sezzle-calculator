# expression-evaluation Specification

## Purpose
Evaluate arithmetic expression strings in Go, independent of any transport, with explicit precedence, associativity, and controlled errors.

## Requirements

### Requirement: Independent evaluator
The system SHALL provide a Go package `internal/calculator` exposing `Evaluate(expression string) (float64, error)` that does not depend on `net/http` or any transport package, and SHALL NOT use eval-like mechanisms or external evaluators.

#### Scenario: Pure function
- **WHEN** `Evaluate("1 + 2")` is called from a plain Go test
- **THEN** it returns `3` and a nil error without any HTTP setup

### Requirement: Arithmetic operators and precedence
The evaluator SHALL support `+`, `-`, `*`, `/`, `^`, parentheses, unary `+`/`-`, and decimal numbers with precedence (high to low) postfix `%`, `^` (right-associative), unary sign, `*` and `/` (left-associative), `+` and `-` (left-associative).

#### Scenario: Multiplication before addition
- **WHEN** evaluating `1 + 2 * 3`
- **THEN** the result is `7`

#### Scenario: Parentheses override precedence
- **WHEN** evaluating `(1 + 2) * 3`
- **THEN** the result is `9`

#### Scenario: Left associativity
- **WHEN** evaluating `10 - 4 - 3` and `100 / 10 / 5`
- **THEN** the results are `3` and `2`

#### Scenario: Exponentiation is right-associative
- **WHEN** evaluating `2 ^ 3 ^ 2`
- **THEN** the result is `512`

#### Scenario: Unary minus is below exponentiation
- **WHEN** evaluating `-2 ^ 2` and `2 ^ -2`
- **THEN** the results are `-4` and `0.25`

#### Scenario: Unary operators
- **WHEN** evaluating `-5 + 2`, `--5`, and `+-5`
- **THEN** the results are `-3`, `5`, and `-5`

#### Scenario: Decimals
- **WHEN** evaluating `.5 + 1.25`
- **THEN** the result is `1.75`

#### Scenario: Whitespace ignored
- **WHEN** evaluating `  1+   2 *3 `
- **THEN** the result is `7`

#### Scenario: Reference expression
- **WHEN** evaluating `1 + 2 - 4 * 3 * (5 - 1) ^ 0.5 + sqrt(16)`
- **THEN** the result is `-17`

### Requirement: Square root
The evaluator SHALL support `sqrt(expression)` and SHALL reject arguments outside the real-number domain with `ErrInvalidSqrt`.

#### Scenario: Valid square root
- **WHEN** evaluating `sqrt(16)` and `sqrt(2 + 2) * 3`
- **THEN** the results are `4` and `6`

#### Scenario: Negative argument
- **WHEN** evaluating `sqrt(-1)`
- **THEN** the error satisfies `errors.Is(err, ErrInvalidSqrt)`

#### Scenario: Missing parentheses
- **WHEN** evaluating `sqrt 16`
- **THEN** the error satisfies `errors.Is(err, ErrSyntax)`

### Requirement: Postfix percentage
The evaluator SHALL treat `%` as a postfix operator dividing its operand by 100, with no contextual calculator-style behavior.

#### Scenario: Simple percentage
- **WHEN** evaluating `50%`
- **THEN** the result is `0.5`

#### Scenario: Percentage is not contextual
- **WHEN** evaluating `100 + 10%`
- **THEN** the result is `100.1`

#### Scenario: Percentage precedence
- **WHEN** evaluating `-50%`, `50%^2`, and `2^50%`
- **THEN** the results are `-0.5`, `0.25`, and `sqrt(2)`

### Requirement: Controlled arithmetic errors
The evaluator SHALL return sentinel errors instead of panicking or returning non-finite values.

#### Scenario: Division by zero
- **WHEN** evaluating `1 / 0` or `1 / (2 - 2)` or `0 ^ -1`
- **THEN** the error satisfies `errors.Is(err, ErrDivisionByZero)`

#### Scenario: Non-finite result
- **WHEN** evaluating `(-8) ^ 0.5` or `10 ^ 1000`
- **THEN** the error satisfies `errors.Is(err, ErrInvalidNumber)`

### Requirement: Controlled syntax errors
The evaluator SHALL return `ErrEmpty` for empty or whitespace-only input and `ErrSyntax` for malformed input, and SHALL NOT panic for any input.

#### Scenario: Empty input
- **WHEN** evaluating `""` or `"   "`
- **THEN** the error satisfies `errors.Is(err, ErrEmpty)`

#### Scenario: Malformed expressions
- **WHEN** evaluating each of `1 +`, `* 2`, `(1 + 2`, `1 + 2)`, `()`, `1 2`, `5.`, `1e3`, `2(3)`, `abc`, `1 $ 2`
- **THEN** each error satisfies `errors.Is(err, ErrSyntax)`

### Requirement: Input size limits
The evaluator SHALL return `ErrTooComplex` for input longer than 1000 bytes or nested deeper than 100 levels, so that hostile input cannot exhaust the stack.

#### Scenario: Deep nesting
- **WHEN** evaluating 200 nested opening parentheses (within the length limit) or 10,000 parentheses (over it)
- **THEN** the error satisfies `errors.Is(err, ErrTooComplex)` and the process does not crash
