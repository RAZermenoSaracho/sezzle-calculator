# Claude Code Prompts

These prompts define the intended OpenSpec workflow for the Sezzle calculator challenge.

Read `CLAUDE.md` before executing any prompt in this file.

# PROPOSAL PROMPTS

## Create All Proposals

Use this prompt once at the beginning of the project.

```text
Read CLAUDE.md and inspect the complete current repository.

We are beginning the proposal phase for the Sezzle calculator challenge.

Create the following OpenSpec changes using /opsx:propose:

1. bootstrap-fullstack-foundation
2. implement-go-expression-engine
3. implement-calculator-api
4. build-live-calculator-ui
5. add-quality-documentation-and-docker
6. add-github-ci-cd

This is PROPOSAL WORK ONLY.

Do not implement application functionality.
Do not commit.
Do not push.
Do not archive changes.

Use subagents in parallel where useful to inspect the repository and develop the proposals faster, but you are responsible for reconciling their findings into one coherent implementation plan. The proposals must not contradict each other or duplicate ownership of the same work.

The six changes should form an intentional dependency sequence.

CHANGE 1 — bootstrap-fullstack-foundation

Define the minimal project foundation required by later changes.

Scope should include:

- repository structure
- React + TypeScript + Vite frontend
- Tailwind CSS integration
- minimal frontend development configuration
- Go backend module and package structure
- minimal backend executable/server foundation where necessary for later changes
- local development commands
- sensible ignore files
- formatting/linting/testing foundations needed by later changes

Do not implement calculator behavior in this change.

Keep the structure simple. Do not introduce a monorepo framework or unnecessary tooling.

CHANGE 2 — implement-go-expression-engine

Implement the calculator domain independently from HTTP.

Define a small expression evaluator in Go capable of parsing a complete expression string.

Required semantics:

- decimal numbers
- whitespace
- parentheses
- unary + and -
- exponentiation with ^
- multiplication
- division
- addition
- subtraction
- sqrt(...)
- postfix percentage, where 50% means 0.5
- correct mathematical precedence
- controlled syntax errors
- controlled division-by-zero errors
- controlled invalid-square-root errors

The proposal must explicitly define precedence and associativity.

Prefer a tokenizer/lexer plus a small recursive-descent or Pratt parser rather than ad-hoc string replacement.

No eval-like mechanism or external expression evaluator.

The evaluator must be independent of HTTP and thoroughly unit tested, preferably using table-driven Go tests.

CHANGE 3 — implement-calculator-api

Expose the expression engine through a small REST API.

Preferred contract:

POST /api/calculate

Request:

{
  "expression": "1 + 2 * 3"
}

Success:

{
  "result": 7
}

Failure:

{
  "error": "invalid expression"
}

Define:

- JSON request/response types
- request validation
- status-code behavior
- malformed JSON handling
- expression error mapping
- HTTP handler tests
- separation between transport and calculator packages
- CORS or local-development routing strategy if needed
- optional lightweight health endpoint if useful for deployment verification

Do not duplicate calculator logic inside handlers.

CHANGE 4 — build-live-calculator-ui

Build the actual calculator experience.

Technology:

- React
- TypeScript
- Vite
- Tailwind CSS

Visual direction:

- dark mode only
- minimal
- responsive
- centered calculator
- expression input
- no conventional four-sided input box
- thin horizontal divider underneath the expression
- result/error below the divider
- no Calculate button
- no Enter requirement
- no unnecessary navigation, cards, gradients, or decorative UI

Behavior:

- calculation occurs automatically while typing
- use a short debounce if useful
- communicate only with the backend API
- display valid results
- display controlled syntax/API errors
- prevent stale/out-of-order responses from replacing newer results
- work on mobile and desktop

Tests should cover meaningful user behavior, including successful calculation and error display.

CHANGE 5 — add-quality-documentation-and-docker

Complete the reviewer-facing deliverables and local containerization.

Scope:

- final README
- setup instructions
- frontend/backend run instructions
- test instructions
- API documentation
- curl examples
- architecture/design rationale
- assumptions
- supported expression syntax
- AI tooling disclosure
- prompts disclosure/reference
- test coverage commands/reporting
- production-oriented Docker support
- simple full-stack container execution

Prefer the simplest reasonable deployment topology.

Avoid Kubernetes, orchestration platforms, databases, reverse-proxy complexity, or unrelated infrastructure.

The Dockerized application should be verifiable locally or through a compatible Docker environment.

Do not modify Ricardo's server in this change.

CHANGE 6 — add-github-ci-cd

Define GitHub Actions CI and the eventual deployment path.

CI should run automatically for pull requests and/or pushes as appropriate and verify at least:

Frontend:
- dependency installation
- lint
- tests
- build

Backend:
- formatting
- vet
- tests
- build

CD should target Ricardo's existing self-hosted Mac server architecture.

Important deployment constraints:

- Docker is not installed on the macOS host.
- Do not introduce Docker Desktop or Colima.
- Container workloads must use Vagrant.
- There is an existing shared Docker VM at ~/scripts/docker-test-vm.
- That VM uses Ubuntu and Docker and is intended as a reusable Docker host.
- Do not modify or interfere with the production n8n VM.
- The server already uses one GitHub Actions self-hosted runner per deployed repository.
- Any new runner registration or server mutation requires operator involvement.
- Cloudflare/DNS/server exposure is outside automatic implementation unless explicitly approved by the operator.

The proposal should clearly separate:
1. repository-owned CI/CD files that can be implemented locally;
2. server-side setup that requires Ricardo;
3. secrets/runner configuration that must never be committed.

Do not perform server-side changes during proposal or normal local implementation.

After creating all six proposals:

1. Validate all OpenSpec changes strictly.
2. Review the proposals together for duplicated scope, missing dependencies, contradictions, or unnecessary complexity.
3. Correct those issues in the proposals.
4. Show me the final dependency/order of implementation.
5. Summarize each change and its acceptance criteria.
6. Report OpenSpec validation results.
7. Stop.

Do not apply any change yet.
```

# IMPLEMENTATION PROMPTS

Run these one at a time and in this order.

Do not start the next change until the current change has passed the manual verification and closure workflow.

## 1. Foundation

```text
/opsx:apply bootstrap-fullstack-foundation

Follow CLAUDE.md exactly.

Implement only this OpenSpec change, run all relevant automated checks, update its task state accurately, and stop before archive, commit, or push.

Report what I need to verify manually.
```

## 2. Expression Engine

```text
/opsx:apply implement-go-expression-engine

Follow CLAUDE.md exactly.

Implement only this OpenSpec change, run all relevant automated checks, update its task state accurately, and stop before archive, commit, or push.

Report what I need to verify manually.
```

## 3. Calculator API

```text
/opsx:apply implement-calculator-api

Follow CLAUDE.md exactly.

Implement only this OpenSpec change, run all relevant automated checks, update its task state accurately, and stop before archive, commit, or push.

Report what I need to verify manually.
```

## 4. Live Calculator UI

```text
/opsx:apply build-live-calculator-ui

Follow CLAUDE.md exactly.

Implement only this OpenSpec change, run all relevant automated checks, update its task state accurately, and stop before archive, commit, or push.

Start the local application if useful for operator verification, but do not make server-side changes.

Report the exact local URLs and the manual behaviors I should verify.
```

## 5. Documentation and Docker

```text
/opsx:apply add-quality-documentation-and-docker

Follow CLAUDE.md exactly.

Implement only this OpenSpec change, run all relevant automated checks, update its task state accurately, and stop before archive, commit, or push.

Do not modify Ricardo's server.

If Docker cannot be executed directly on the development machine, distinguish between artifacts that were statically validated and behavior that still requires execution in a Docker-capable environment.

Report what I need to verify manually.
```

## 6. CI/CD

```text
/opsx:apply add-github-ci-cd

Follow CLAUDE.md exactly.

Implement repository-owned CI/CD configuration only.

Do not register runners, modify Vagrant VMs, install software on the server, edit Cloudflare configuration, create DNS records, expose ports, or mutate production infrastructure without explicit operator approval.

Run every local/static validation that can reasonably be performed.

At the end, separate the report into:

1. repository work completed;
2. automated verification completed;
3. GitHub configuration still required;
4. server-side actions still required;
5. manual verification still required.

Stop before archive, commit, or push.
```

# CLOSURE PROMPT

Use this only after Ricardo has manually reviewed an applied change and explicitly approved it.

Replace `<change-name>` with the current OpenSpec change.

```text
The change <change-name> has been manually reviewed and approved.

Follow the closure procedure defined in CLAUDE.md.

Before closing:

1. Inspect git status and git diff.
2. Re-read the approved OpenSpec change and its tasks.
3. Confirm that the working tree changes belonging to this change match the approved scope.
4. Do not include unrelated modifications.
5. Run the appropriate final automated verification for this change.
6. If verification fails, stop and report the failure instead of closing.

If verification succeeds:

7. Mark the approved OpenSpec tasks complete accurately.
8. Validate OpenSpec.
9. Archive the completed change according to OpenSpec conventions.
10. Inspect the resulting diff again.
11. Commit only this approved change using a concise conventional commit message.
12. Push the current branch to its configured remote.
13. Confirm the final git status.

Report:

- archived OpenSpec change
- verification commands and results
- commit hash
- commit message
- pushed branch
- final git status
- next change in the implementation sequence, if any

Do not begin implementing the next change.
```

# Expected Change Order

```text
bootstrap-fullstack-foundation
        ↓
implement-go-expression-engine
        ↓
implement-calculator-api
        ↓
build-live-calculator-ui
        ↓
add-quality-documentation-and-docker
        ↓
add-github-ci-cd
```

The order is intentional.

The expression engine must exist independently before HTTP is added. The API must exist before the frontend depends on it. The working product should exist before packaging/documentation are finalized. CI/CD comes last so it validates and deploys the actual completed build rather than an evolving scaffold.