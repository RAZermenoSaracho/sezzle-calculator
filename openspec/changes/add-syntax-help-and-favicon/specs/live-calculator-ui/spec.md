## ADDED Requirements

### Requirement: Supported syntax help
The UI SHALL show a compact, visually secondary section below the calculator listing the supported syntax: `+` and `-`, `*` and `/`, `^`, `sqrt(x)`, postfix `x%`, parentheses, unary `+x` and `-x`, and decimal numbers.

#### Scenario: Help is visible with the calculator
- **WHEN** the page loads
- **THEN** the calculator and the supported-syntax section are both rendered, with the section below the calculator

### Requirement: Repeated unary sign explanation
The UI SHALL explain that unary `+` and `-` can be chained and are applied repeatedly, that an even number of `-` signs cancels out, and that the signs must be followed by a value, so `++++5` is valid while `++++` alone is not. The explanation SHALL match the backend grammar.

#### Scenario: Explanation content
- **WHEN** a user reads the explanation
- **THEN** it gives a chained-sign example that evaluates as stated and states that a sign-only string is not a complete expression

### Requirement: Application favicon
The application SHALL provide an SVG favicon at `/favicon.svg`, referenced from the HTML, that suits a calculator and uses no external dependency.

#### Scenario: Favicon referenced
- **WHEN** the built page is loaded
- **THEN** its `<head>` links `/favicon.svg` as the icon and the file is served as SVG
