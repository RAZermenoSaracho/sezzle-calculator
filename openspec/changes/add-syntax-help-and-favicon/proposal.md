# Proposal: Syntax Help and Favicon

## Why

Users cannot tell what the calculator accepts, and expressions such as `----5` or `54+++++6666` look like mistakes even though they are valid. A short in-page explanation makes the grammar discoverable and honest. The page also still shows the untouched Vite template favicon.

## What Changes

- Add a compact, secondary "Supported syntax" section below the calculator listing every supported operation, plus a short "Why are `++++` and `----` allowed?" explanation that matches the parser.
- Replace the template favicon with a minimal calculator-style SVG.
- Frontend only: no change to the parser, API, or calculator behavior.

## Capabilities

### New Capabilities
None.

### Modified Capabilities
- `live-calculator-ui`: adds requirements for the syntax help section and the favicon.

## Non-goals

Any backend or grammar change (repeated unary signs stay valid), interactive help, examples that fill the input, icons libraries, or new dependencies.

## Impact

- New: `frontend/src/components/SyntaxHelp.tsx` and its test; modified: `App.tsx`, `Calculator.tsx` (layout wrapper moves to `App`), `public/favicon.svg`, tests.
- Depends on: archived `build-live-calculator-ui`.
