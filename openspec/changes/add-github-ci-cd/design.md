# Design: GitHub CI/CD

## Boundaries

| Category | Owner | Contents |
|---|---|---|
| Repository-owned (this change) | Claude, locally | `ci.yml`, `deploy.yml`, `deploy/deploy.sh`, `deploy/README.md`, README CI/CD section |
| Requires Ricardo | Operator | runner registration for this repo and its labels, GitHub Environment `production`, repository variables/secrets, verifying the Vagrant VM sync/port-forward details, any Cloudflare/DNS/exposure |
| Never committed | Operator | tokens, SSH keys, Cloudflare credentials, `.env` files |

## CI (`ci.yml`)

Triggers: `pull_request`, `push` to `main`. `permissions: contents: read`. Concurrency group per ref with cancel-in-progress.

- **frontend** (ubuntu-latest, `working-directory: frontend`): `actions/setup-node` with Node 24 and npm cache, `npm ci`, `npm run lint`, `npm test -- --run`, `npm run build`.
- **backend** (ubuntu-latest, `working-directory: backend`): `actions/setup-go` with `go-version-file: backend/go.mod`, formatting check (`test -z "$(gofmt -l .)"`), `go vet ./...`, `go test ./...`, `go build ./...`.
- **docker**: `docker build .` with no push, validating the Dockerfile on a runner that has Docker.

Actions pinned to major versions. No matrix, no coverage upload, no third-party actions beyond `actions/*`.

## CD (`deploy.yml`)

- Trigger: `workflow_run` of `CI` completed on `main` with conclusion `success` (plus `workflow_dispatch` for manual runs).
- Job condition: `vars.DEPLOY_ENABLED == 'true'`, so merging this change deploys nothing until Ricardo opts in.
- `runs-on: [self-hosted, macOS]` (label finalized with the operator's runner), `environment: production`, `concurrency: deploy`.
- Steps: checkout the tested commit (`github.event.workflow_run.head_sha`), run `deploy/deploy.sh`, then poll `GET /api/health` on the published URL, failing the job if it does not become healthy.

## `deploy/deploy.sh`

Runs on the Mac host and never installs Docker on it. It drives the existing shared VM through Vagrant:

1. `cd "$DOCKER_VM_DIR"` (default `~/scripts/docker-test-vm`) and use `vagrant ssh -c ...`.
2. Transfer the commit into the VM (mechanism: `git archive HEAD | vagrant ssh -c 'tar -x -C <dir>'`, which needs no shared-folder assumptions).
3. In the VM: `docker build -t sezzle-calculator:<sha> .`, then replace only the container named `sezzle-calculator` (`docker rm -f` of that exact name, then `docker run -d --restart unless-stopped -p <HOST_PORT>:8080 --name sezzle-calculator ...`).
4. Never references other containers or VMs (n8n is untouched); refuses to run if `DOCKER_VM_DIR` is unset/absent; uses `set -euo pipefail`.

All environment-specific values (`DOCKER_VM_DIR`, `HOST_PORT`, `CONTAINER_NAME`) come from GitHub repository variables or defaults documented in `deploy/README.md`; no secrets are needed by the script. The script is validated locally only with `bash -n` and `shellcheck` if available; it is never executed against the server during local implementation.

## Open questions for the operator (recorded in `deploy/README.md`)

1. Is the shared VM acceptable for a long-running container, given it is described as ephemeral/test?
2. Which host port is free and how does the VM forward it to the Mac (Vagrant `forwarded_port`), and how is it exposed (Cloudflare tunnel/DNS)?
3. Runner labels and the user the runner runs as (needs `vagrant` on PATH).
4. Whether the `production` Environment should require manual approval.

Answers may adjust `deploy.sh`; they do not change CI.
