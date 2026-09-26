## ADDED Requirements

### Requirement: Gated deployment workflow
The repository SHALL provide a deployment workflow that runs on the self-hosted runner only after CI succeeds on `main` (or by manual dispatch), only when repository variable `DEPLOY_ENABLED` equals `true`, and within a GitHub Environment named `production`.

#### Scenario: Not enabled
- **WHEN** `DEPLOY_ENABLED` is unset
- **THEN** the deploy job is skipped and nothing is deployed

#### Scenario: CI failed
- **WHEN** CI concludes with failure on `main`
- **THEN** no deployment runs

### Requirement: Vagrant-based container deployment
The deployment script SHALL build and run the application image inside the existing Docker VM via Vagrant, SHALL NOT require or install Docker on the macOS host, and SHALL replace only the container named `sezzle-calculator`.

#### Scenario: Other workloads untouched
- **WHEN** the script runs
- **THEN** it does not stop, modify, or reference any other container or VM, including n8n

### Requirement: Post-deploy verification
The deployment workflow SHALL fail unless `GET /api/health` returns healthy after deployment.

#### Scenario: Unhealthy deploy
- **WHEN** the health check does not succeed within its timeout
- **THEN** the workflow fails

### Requirement: Operator boundary
This change SHALL NOT register runners, modify Vagrant VMs, install software on the server, edit Cloudflare or DNS, expose ports, or commit secrets; `deploy/README.md` SHALL list every action requiring the operator and every unresolved question.

#### Scenario: Local implementation only
- **WHEN** the change is applied
- **THEN** only repository files are created and `deploy.sh` is validated statically (`bash -n`, `shellcheck` if available) without being executed against the server

### Requirement: No committed secrets
Credentials, tokens, and keys SHALL be supplied only through GitHub secrets, variables, or the runner environment.

#### Scenario: Repository scan
- **WHEN** searching the repository for tokens or private keys
- **THEN** none are found
