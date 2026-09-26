# Design: Go Expression Engine

## Grammar (authoritative)

```text
expression := additive EOF
additive   := term   { ('+' | '-') term }
term       := unary  { ('*' | '/') unary }
unary      := ('+' | '-') unary | power
power      := postfix [ '^' unary ]
postfix    := primary { '%' }
primary    := NUMBER | '(' additive ')' | 'sqrt' '(' additive ')'
NUMBER     := DIGIT+ [ '.' DIGIT+ ] | '.' DIGIT+
```

## Precedence and associativity (highest to lowest)

| Level | Operator | Form | Associativity |
|---|---|---|---|
| 1 | `%` | postfix | left (`50%%` = 0.005) |
| 2 | `^` | binary | **right** (`2^3^2` = 512) |
| 3 | unary `+` `-` | prefix | right |
| 4 | `*` `/` | binary | left |
| 5 | `+` `-` | binary | left |

Consequences, all covered by tests:

- `-2^2` = -4 (unary is below `^`); `2^-2` = 0.25 (the exponent is a `unary`, so a sign is allowed there); `2^-3^2` = 2^(-(3^2)).
- `--5` = 5; `+-5` = -5.
- `2^50%` = 2^0.5 and `50%^2` = 0.25 (`%` binds tightest); `-50%` = -0.5.
- `100 + 10%` = 100.1 (postfix only, no contextual behavior).
- `sqrt(16)` = 4; `sqrt(2)^2` parses as `(sqrt(2))^2`; `sqrt` is a primary and requires parentheses.

## Lexing rules

- Whitespace (space, tab, newline) between tokens is ignored; whitespace inside a number (`1 2`) yields two tokens and therefore a syntax error.
- Numbers: `12`, `1.5`, `.5` valid; `5.`, `1.2.3`, `1e3` invalid. No thousands separators.
- Identifier `sqrt` (lowercase) is the only allowed word. Any other character or word is a syntax error (`ErrSyntax`).
- Work on bytes/runes without regexp.

## Evaluation and errors

Parser evaluates as it parses, returning `float64`. Use `strconv.ParseFloat` on the number lexeme; this cannot fail for valid lexemes except overflow, which maps to `ErrInvalidNumber`.

Exported sentinel errors (`errors.Is`), each possibly wrapped with an internal position for tests/logging:

| Error | Trigger |
|---|---|
| `ErrEmpty` | empty or whitespace-only input |
| `ErrSyntax` | unexpected/unsupported character, unexpected token, missing `(` or `)`, trailing tokens, `sqrt` without parentheses |
| `ErrDivisionByZero` | `x / 0`, and `0 ^ negative` |
| `ErrInvalidSqrt` | `sqrt` of a negative number |
| `ErrInvalidNumber` | any result or intermediate that is NaN or ±Inf (overflow, `(-8)^0.5`) |
| `ErrTooComplex` | input longer than 1000 bytes or nesting deeper than 100 levels; prevents stack exhaustion |

`0 ^ negative` is checked before calling `math.Pow`, because `Pow(0, -1)` returns +Inf and would otherwise be reported as `ErrInvalidNumber`. The 1000-byte length limit and the 100-level depth limit are both enforced; deep nesting within 1000 bytes is what the depth check catches.

Error `Error()` strings are for developers; the HTTP layer maps by `errors.Is`, never by string. `-0` results are normalized to `0`.

## Structure

- `calculator.go`: `Evaluate`, sentinel errors.
- `lexer.go`: tokens (`number`, `+ - * / ^ % ( )`, `sqrt`, `eof`).
- `parser.go`: recursive descent implementing the grammar; one method per grammar rule; depth counter for `ErrTooComplex`.
- Tests: `calculator_test.go` table-driven with `{name, input, want, wantErr}`; float comparison with a small epsilon.

Recursive descent was chosen over Pratt because the grammar has one postfix and one right-associative operator and maps one-to-one to the rules above, which is easy for a reviewer to verify.
