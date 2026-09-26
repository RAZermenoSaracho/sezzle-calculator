# Deployment

The `Deploy` workflow (`.github/workflows/deploy.yml`) builds the committed tree inside the shared Docker VM through Vagrant and replaces the `sezzle-calculator` container. Docker is never installed or used on the macOS host.

It is **disabled by default**: nothing runs until the repository variable `DEPLOY_ENABLED` is `true`.

## Flow

1. `CI` succeeds for a push to `main` (or someone runs `Deploy` manually).
2. On the self-hosted runner, `deploy/deploy.sh` archives `HEAD`, streams it into the VM with `vagrant ssh`, runs `docker build`, replaces only the container named `sezzle-calculator`, and starts it with `--restart unless-stopped -p $HOST_PORT:8080`.
3. The workflow polls `GET /api/health` and fails if it does not become healthy within about a minute.

`deploy.sh` touches no other container, image, or VM, and needs no secrets.

## What lives where

| Category | Items |
|---|---|
| Repository-owned (this repo) | `ci.yml`, `deploy.yml`, `deploy/deploy.sh`, this file |
| GitHub configuration (Ricardo) | runner registered for this repository; Environment `production`; repository variables below |
| Server-side (Ricardo) | runner service and its `PATH` (needs `vagrant`); the Vagrant VM and its port forward; Cloudflare Tunnel and DNS |
| Never committed | tokens, SSH keys, Cloudflare credentials, `.env` files |

## Repository variables

| Variable | Required | Default | Meaning |
|---|---|---|---|
| `DEPLOY_ENABLED` | yes | (unset, so disabled) | set to `true` to allow deployments |
| `DOCKER_VM_DIR` | no | `~/scripts/docker-test-vm` | Vagrant directory of the Docker VM |
| `HOST_PORT` | no | `8080` | port published inside the VM |
| `DEPLOY_HEALTH_URL` | no | `http://127.0.0.1:8080/api/health` | URL polled after deploy, e.g. the public URL |

## Enabling deployment (operator checklist)

1. Register the repository's self-hosted runner. The workflow targets labels `self-hosted` and `macOS`.
2. Make sure `vagrant` is on the runner service's `PATH` and the VM is running.
3. Create the GitHub Environment `production`. Consider requiring manual approval.
4. Set `DEPLOY_ENABLED=true` and any optional variables above.
5. Run `Deploy` once via **Run workflow** and check the health step.

## Verified facts and remaining questions

Ricardo verified manually that the image builds and runs in the Ubuntu Vagrant VM, that port 8080 is forwarded to the Mac, and that the app is reachable through the Cloudflare Tunnel at https://sezzle-calculator.razs.dev.

Still to confirm before enabling:

- Is the shared VM, described as an ephemeral test host, acceptable for a long-running container?
- The runner labels and the user the runner service runs as.
- Whether `production` deployments should require approval.
- Whether the runner can run `vagrant ssh` non-interactively (no TTY); `deploy.sh` passes `-T` for this.

`deploy.sh` has only been validated statically (`bash -n`); it has not been run against the server.
