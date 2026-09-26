# Tasks

## 1. CI
- [x] 1.1 Write `.github/workflows/ci.yml` per `design.md`. Verify with `actionlint` if available, otherwise YAML parse and manual review of triggers, working directories, and versions.
- [x] 1.2 Run every CI command locally (`npm ci`, lint, test, build; gofmt check, vet, test, build). Verify all succeed.

## 2. Deployment (repository files only)
- [x] 2.1 Write `deploy/deploy.sh` per `design.md`. Verify `bash -n` and `shellcheck` (if available); do not execute it.
- [x] 2.2 Write `.github/workflows/deploy.yml` with the `DEPLOY_ENABLED` guard, `production` environment, `workflow_run` trigger, and health check. Verify with `actionlint` or manual review.
- [x] 2.3 Write `deploy/README.md` listing operator actions and the open questions from `design.md`. Verify it separates repository work, GitHub configuration, and server actions.

## 3. Docs and hygiene
- [x] 3.1 Update the README CI/CD section. Verify it matches the workflows.
- [x] 3.2 Search the repo for secrets, tokens, keys, and `.env` files. Verify none exist.

## 4. Report
- [x] 4.1 Final report separates: repository work completed; automated verification; GitHub configuration still required; server-side actions still required; manual verification. Verify no server, runner, VM, Cloudflare, or DNS change was made.

## Verification record

Static validation during apply: `actionlint` v1.7.12 clean on both workflows, `bash -n deploy/deploy.sh`, all CI commands run locally, `deploy.sh` refusal paths (missing Vagrantfile, non-numeric port) exit 1 without side effects, and a secret scan found nothing. `openspec validate --strict` passed.

Manually verified by Ricardo in the real deployment environment:

- The repository self-hosted runner is registered with labels `self-hosted`, `macOS`, and `X64`.
- It runs as the macOS launchd service `actions.runner.RAZermenoSaracho-sezzle-calculator.sezzle-calculator`.
- The runner can reach Vagrant at `/usr/local/bin/vagrant`.
- The Docker VM is `/Users/razs/production/n8n-vm`; it is running and Docker is available inside it.
- `vagrant ssh -c ... -- -T` works non-interactively.
- The GitHub Environment `production` is configured with `DOCKER_VM_DIR=/Users/razs/production/n8n-vm`, `HOST_PORT=8080`, and `DEPLOY_HEALTH_URL=http://127.0.0.1:8080/api/health`.
- The repository variable `DEPLOY_ENABLED=true` is set, so the job-level gate is evaluated before the environment is entered.
- A manual run of the Deploy workflow against commit `87ca08b` succeeded. It rebuilt and replaced the production container, which was verified as `Image=sezzle-calculator:87ca08b`, `Restart=unless-stopped`, newly created, exposing port 8080.
- Local health `http://127.0.0.1:8080/api/health` returned `{"status":"ok"}`.
- Public health `https://sezzle-calculator.razs.dev/api/health` through Cloudflare returned `{"status":"ok"}`.
- A public calculation of `sqrt(16) + 2^3` returned `{"result":12}`.

The manually triggered path is therefore verified end to end: GitHub Actions, self-hosted macOS runner, `deploy.sh`, Vagrant VM, Docker build and container replacement, local health, and the public application.

Not yet verified at the time of closure: the automatic `push to main`, CI success, `workflow_run` Deploy trigger. The closeout push to `main` is intended to exercise it; the result is observed in GitHub Actions after the push and is not part of this change's automated checks.
