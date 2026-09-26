# Tasks

## 1. CI
- [ ] 1.1 Write `.github/workflows/ci.yml` per `design.md`. Verify with `actionlint` if available, otherwise YAML parse and manual review of triggers, working directories, and versions.
- [ ] 1.2 Run every CI command locally (`npm ci`, lint, test, build; gofmt check, vet, test, build). Verify all succeed.

## 2. Deployment (repository files only)
- [ ] 2.1 Write `deploy/deploy.sh` per `design.md`. Verify `bash -n` and `shellcheck` (if available); do not execute it.
- [ ] 2.2 Write `.github/workflows/deploy.yml` with the `DEPLOY_ENABLED` guard, `production` environment, `workflow_run` trigger, and health check. Verify with `actionlint` or manual review.
- [ ] 2.3 Write `deploy/README.md` listing operator actions and the open questions from `design.md`. Verify it separates repository work, GitHub configuration, and server actions.

## 3. Docs and hygiene
- [ ] 3.1 Update the README CI/CD section. Verify it matches the workflows.
- [ ] 3.2 Search the repo for secrets, tokens, keys, and `.env` files. Verify none exist.

## 4. Report
- [ ] 4.1 Final report separates: repository work completed; automated verification; GitHub configuration still required; server-side actions still required; manual verification. Verify no server, runner, VM, Cloudflare, or DNS change was made.
