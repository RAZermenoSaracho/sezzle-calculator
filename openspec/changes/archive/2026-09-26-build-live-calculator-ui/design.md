# Design: Live Calculator UI

## Layout

```text
1 + 2 * 3          <- <textarea rows=1>, transparent, no border, large text, autofocus,
                      wraps and grows vertically with its content
────────────       <- <hr>/border-t, thin, muted
7                  <- <output aria-live="polite"> result, or error in a muted red
```

Container: `min-h-screen` flex, centered, `max-w-xl w-full px-4`; input uses `text-2xl sm:text-4xl`, `bg-transparent outline-none`, a visible focus indication via the divider color (not a box). The textarea height is set from `scrollHeight` in a layout effect (works in every browser, no dependency). Enter is ignored and pasted newlines are replaced by spaces, so the expression stays one logical line. `autoComplete="off"`, `spellCheck={false}`, `aria-label="Expression"`. Works at 360 px width and up. Palette: black background with Matrix-style green (`green-400`) for expression and result, `green-900` divider brightening to `green-500` on focus, `green-900` placeholder, and `red-400` for errors so they never read as results. No glow, animation, or decoration. Tailwind utilities only; `index.css` contains just the Tailwind import.

## Files

- `src/calculate.ts`: `calculate(expression, signal): Promise<{result: number} | {error: string}>` doing `fetch('/api/calculate', ...)`. HTTP 200 → result; any 4xx/5xx with `{error}` → that message; network failure or unparseable body → `"Could not reach the server"`.
- `src/components/Calculator.tsx`: state (`expression`, `output`) and one `useEffect` implementing the lifecycle. No custom hook unless it clarifies tests.
- `src/App.tsx`: thin root that renders `Calculator`.
- Tests live in `__tests__` directories next to the code they test (`src/__tests__/calculate.test.ts`, `src/components/__tests__/Calculator.test.tsx`); shared setup stays in `src/test/setup.ts`.

## Request lifecycle

Effect keyed on `expression`:
1. Trimmed-empty input → clear output, issue no request.
2. Otherwise `setTimeout(300)`; on fire, create an `AbortController` and call `calculate`.
3. Effect cleanup clears the timer and aborts the in-flight request, so a response for an older expression can never overwrite a newer one; aborted requests are ignored silently.
4. While a request is pending the previous output stays visible (no flicker, no spinner).

Result display is presentation only: `Number(result.toPrecision(12)).toString()` to hide binary float noise (`0.1 + 0.2` → `0.3`). Very large/small values may render in exponent form; acceptable. This is formatting, not calculation, and is the only numeric handling in the frontend.

Error text is shown exactly as returned by the API (`error` field), keeping API and UI consistent.

## Testing (Vitest + Testing Library, `fetch` mocked, fake timers)

- typing then advancing 300 ms triggers exactly one request with the final text (debounce);
- successful result shown; error message shown;
- clearing the input removes output and sends no request;
- out-of-order: first request resolves after second → only the newer result is shown;
- network failure shows the fallback message.
