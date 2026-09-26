# Tasks

## 1. Frontend shell
- [x] 1.1 Remove template assets and `App.css`; reduce `App.tsx`/`main.tsx` to a placeholder heading. Verify `npm run build` succeeds.
- [x] 1.2 Install `tailwindcss` and `@tailwindcss/vite`; register the plugin in `vite.config.ts`; set `index.css` to `@import "tailwindcss";`; dark body classes in `index.html`. Verify a Tailwind class visibly applies in `npm run dev`.
- [x] 1.3 Add `/api` proxy to `http://localhost:8080` in `vite.config.ts`. Verify by reading the config and, once a backend is running, `curl localhost:5173/api/x` reaches it (404 from Go is acceptable).

## 2. Frontend testing
- [x] 2.1 Install Vitest, jsdom, Testing Library packages; configure jsdom + setup file in `vite.config.ts`; add `"test": "vitest"` script. Verify `npm test -- --run` runs.
- [x] 2.2 Add one smoke test that renders `App`. Verify it passes and `npm run lint` is clean.

## 3. Backend shell
- [x] 3.1 Create `backend/cmd/server/main.go` with `PORT` handling, `newHandler()` returning an empty mux, and a server with `ReadHeaderTimeout`. Verify `go run ./cmd/server` listens on 8080.
- [x] 3.2 Run `gofmt -l .`, `go vet ./...`, `go test ./...`, `go build ./...` in `backend/`. Verify all succeed.

## 4. Repository hygiene
- [x] 4.1 Review `.gitignore` covers `node_modules`, `dist`, `coverage`, `backend/bin`, `.env*`. Verify with `git status` that no generated files appear.
- [x] 4.2 Confirm no calculator logic, API routes, Docker, or CI files were added. Verify with `git status`.
