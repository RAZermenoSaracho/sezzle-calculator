# Proposal: GitHub CI/CD

## Why

The completed build should be verified automatically on every change and deployable to Ricardo's self-hosted Mac server through the existing per-repository runner model, without the repository owning any secret or server state.

## What Changes

- Add `.github/workflows/ci.yml`: on pull requests and pushes to `main`, run frontend install/lint/test/build, backend gofmt check/vet/test/build, and a Docker image build (no push) on GitHub-hosted runners.
- Add `.github/workflows/deploy.yml`: runs on the repository's self-hosted runner after CI succeeds on `main`, guarded by a repository variable and a GitHub Environment, invoking a repo-owned `deploy/deploy.sh`.
- Add `deploy/deploy.sh` and `deploy/README.md` describing the parameterized deployment into the shared Vagrant Docker VM and the exact operator prerequisites.
- Update the README's CI/CD section (only that section).
- **No server-side, runner, Vagrant, Cloudflare, DNS, or secret changes are made by this change.**

## Capabilities

### New Capabilities
- `continuous-integration`: automated verification workflow.
- `continuous-deployment`: repository-owned deployment workflow and its operator boundary.

### Modified Capabilities
None.

## Non-goals

Registries, release tagging, multi-environment promotion, Docker Desktop/Colima, touching the production n8n VM, registering runners, exposing ports or DNS.

## Impact

- New: `.github/workflows/{ci,deploy}.yml`, `deploy/`.
- Depends on: `add-quality-documentation-and-docker` (Dockerfile, health endpoint, README).
