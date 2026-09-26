# Tasks

## 1. API client
- [ ] 1.1 Implement `src/calculate.ts` per `design.md` with abort support. Verify with `calculate.test.ts` covering 200, 400 with error, non-JSON body, and network failure (mocked `fetch`).

## 2. Interface
- [ ] 2.1 Build the centered layout in `App.tsx` with Tailwind utilities only (borderless input, divider, output with `aria-live`). Verify visually at 360 px and 1280 px in `npm run dev`.
- [ ] 2.2 Implement the debounce/abort effect and result formatting. Verify manually against the running backend: `1 + 2 * 3` → `7`, `1 / 0` → `division by zero`, `sqrt(16)` → `4`.
- [ ] 2.3 Remove leftover template/placeholder code and the foundation smoke test. Verify `npm run build` is clean.

## 3. Tests
- [ ] 3.1 `App.test.tsx` with fake timers covering debounce, success, API error, network error, empty input, and out-of-order responses. Verify `npm test -- --run` passes.

## 4. Quality gate
- [ ] 4.1 Run `npm run lint`, `npm test -- --run`, `npm run build`. Verify all succeed.
- [ ] 4.2 Start backend and frontend locally and report URLs (`http://localhost:5173`, API `http://localhost:8080`) and manual checks: live update while typing, rapid typing, parentheses, `50%`, error cases, mobile width, no Calculate button.
