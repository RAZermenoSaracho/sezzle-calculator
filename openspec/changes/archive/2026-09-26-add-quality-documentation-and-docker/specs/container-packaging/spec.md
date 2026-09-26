## ADDED Requirements

### Requirement: Single-image deployment
The repository SHALL provide one multi-stage `Dockerfile` producing a single image that serves the built frontend and the `/api` endpoints from one port (default 8080), running as a non-root user, with no additional services.

#### Scenario: Build and run
- **WHEN** the image is built and started with `-p 8080:8080` in a Docker-capable environment
- **THEN** `GET /` returns the calculator page and `GET /api/health` returns `{"status":"ok"}`

#### Scenario: End-to-end calculation
- **WHEN** `POST /api/calculate` with `{"expression":"1 + 2 * 3"}` is sent to the container
- **THEN** the response is `{"result":7}`

### Requirement: Optional static file serving
The Go server SHALL serve static files from the directory named by `STATIC_DIR` at `/` when the variable is set, SHALL leave `/api` routes unaffected, and SHALL serve no static files when it is unset.

#### Scenario: Static enabled
- **WHEN** `STATIC_DIR` points to a directory containing `index.html`
- **THEN** `GET /` returns that file and `GET /api/health` still returns JSON

#### Scenario: Static disabled
- **WHEN** `STATIC_DIR` is unset
- **THEN** `GET /` returns 404

### Requirement: Minimal infrastructure
The packaging SHALL NOT introduce docker-compose, reverse proxies, orchestration, databases, or registry publishing.

#### Scenario: Only a Dockerfile and ignore file
- **WHEN** reviewing files added for containerization
- **THEN** only `Dockerfile` and `.dockerignore` (plus documentation) are added

### Requirement: Honest verification reporting
The implementation report SHALL distinguish artifacts validated statically from behavior that requires a Docker-capable environment, and SHALL NOT modify any server or VM.

#### Scenario: No local Docker
- **WHEN** Docker is unavailable on the development host
- **THEN** the report lists what was validated statically and what remains to be verified by the operator
