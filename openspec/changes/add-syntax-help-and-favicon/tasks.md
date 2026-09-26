# Tasks

## 1. Favicon
- [x] 1.1 Replace `frontend/public/favicon.svg` with the calculator-style mark. Verify the file is valid SVG (`npm run build` copies it to `dist`) and `index.html` links `/favicon.svg`.

## 2. Syntax help
- [x] 2.1 Move page layout from `Calculator` to `App`; keep calculator behavior identical. Verify existing Calculator tests still pass unchanged.
- [x] 2.2 Add `SyntaxHelp` with the supported syntax list and the repeated-sign explanation from `design.md`. Verify against the running API that every example in the text evaluates as stated.
- [x] 2.3 Render `SyntaxHelp` below the calculator in `App`. Verify visually on desktop and phone width.

## 3. Tests and checks
- [x] 3.1 Add `SyntaxHelp` and `App` tests for the operations and explanation claims. Verify `npm test -- --run` passes.
- [x] 3.2 Run `npm run lint`, `npm test -- --run`, `npm run build`, and `openspec validate --strict`. Verify all succeed and no backend files changed (`git diff --stat -- backend` is empty).
