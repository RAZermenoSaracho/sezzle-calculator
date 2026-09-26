# Design: Calculator API

## Contract (authoritative; the UI change consumes exactly this)

`POST /api/calculate`, `Content-Type: application/json`

Request: `{"expression": "1 + 2 * 3"}`

| Outcome | Status | Body |
|---|---|---|
| Success | 200 | `{"result": 7}` |
| Empty/whitespace expression | 400 | `{"error": "empty expression"}` |
| Syntax error | 400 | `{"error": "invalid expression"}` |
| Division by zero | 400 | `{"error": "division by zero"}` |
| Negative sqrt | 400 | `{"error": "invalid square root"}` |
| NaN / overflow | 400 | `{"error": "invalid number"}` |
| Too long / too deep | 400 | `{"error": "expression too complex"}` |
| Malformed JSON, missing `expression`, non-string `expression`, unknown fields, trailing data | 400 | `{"error": "invalid request"}` |
| Body over 4 KiB | 413 | `{"error": "request too large"}` |
| Wrong method on known path | 405 | `{"error": "method not allowed"}` with `Allow: POST` |
| Non-JSON `Content-Type` | 415 | `{"error": "unsupported media type"}` |
| Unexpected internal failure | 500 | `{"error": "internal error"}` |

`GET /api/health` → 200 `{"status": "ok"}`.

All responses set `Content-Type: application/json`. `result` is a JSON number (float64). Errors are always `{"error": string}`, so the UI can show `error` verbatim. Messages are stable and never contain parser internals, positions, or panic text.

## Decisions

- **400 for all evaluation failures**: they are client input errors; one code keeps the frontend trivial. The `error` string distinguishes them.
- **Error mapping** is a single `switch` using `errors.Is` on the calculator sentinels; unknown errors become 500 `internal error`.
- **Strict decoding**: `http.MaxBytesReader` (4 KiB), `json.Decoder.DisallowUnknownFields`, `expression` decoded as `*string` to distinguish missing from empty, and a second `Decode` must return `io.EOF`.
- **Panic safety**: a small `recover` wrapper returns 500 `internal error` without leaking details. The engine does not panic; this is defense in depth.
- **Routing**: Go 1.22+ `ServeMux` method patterns (`POST /api/calculate`, `GET /api/health`). A wrong method yields 405 with `Allow` from the mux; unknown paths keep the mux default 404. A small wrapper renders those mux-generated 405s as the JSON body above only if it stays trivial; otherwise the plain-text mux response is acceptable for 405/404 and the contract table applies to handler-produced responses.
- **Package boundary**: `api` imports `calculator`; `calculator` never imports `api`. Handler takes the evaluator as a `func(string) (float64, error)` field defaulting to `calculator.Evaluate` so tests can inject a failing evaluator to cover the 500 path.

## Local development

Frontend on :5173 proxies `/api` to backend :8080 (configured in the foundation). No CORS headers are added.

## Wiring

`main.go` keeps the foundation's `newHandler()` function, which now returns `api.NewHandler()`. Change 5 extends it to `newHandler(staticDir string)`.
