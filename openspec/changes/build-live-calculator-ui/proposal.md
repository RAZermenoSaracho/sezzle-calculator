# Proposal: Build Live Calculator UI

## Why

The calculator experience is the user-facing deliverable. It must feel immediate, stay minimal, and rely entirely on the backend for calculation semantics.

## What Changes

- Replace the placeholder with a centered, dark-only calculator: an expression input with no four-sided border, a thin horizontal divider, and a result/error line below it.
- Calculate automatically while typing using a 300 ms debounce; no button, no Enter requirement.
- Add a small API client for `POST /api/calculate` (contract from `implement-calculator-api`).
- Guard against stale and out-of-order responses.
- Replace the foundation smoke test with behavior tests.

## Capabilities

### New Capabilities
- `live-calculator-ui`: the live-updating calculator interface and its request lifecycle.

### Modified Capabilities
None.

## Non-goals

Local expression evaluation or validation, history, keypad buttons, themes, animations, routing, state libraries, spinners, extra components beyond what clarity requires.

## Impact

- Modified: `frontend/src/App.tsx`, `index.css`; new: `frontend/src/calculate.ts`, `App.test.tsx`, `calculate.test.ts`.
- Depends on: `implement-calculator-api` (and the Vite proxy from the foundation).
