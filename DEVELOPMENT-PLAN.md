# Development Plan & Roadmap

## Project Goals

1. **Evolutionary Memory**: A system that learns from feedback and evolves confidence over time.
2. **Hybrid Architecture**: High-performance Go backend + Flexible Python ML sidecar.
3. **Human-in-the-Loop**: A Vite SPA to visualize memory graphs, manage anchors, and provide manual feedback.
4. **Cloud Native**: Containerized deployment targeting GCP (Cloud Run).

---

## Immediate TODOs

### Phase 1: Restructuring & Core Backend (Current)

- [ ] **Refactor**: Move existing Go code to `backend/` directory.
- [ ] **Refactor**: Move existing Python code to `ml-sidecar/` directory.
- [ ] **Docker**: Update `docker-compose.yml` to point to new build contexts.
- [ ] **Go**: Implement `MemoryStore` interface backed by Postgres/pgvector.
- [ ] **Go**: Implement `MLClient` to talk to `ml-sidecar`.
- [ ] **Test**: Set up `testcontainers-go` for integration testing in `backend/`.

### Phase 2: ML Sidecar & Embeddings

- [ ] **Python**: Finalize `embeddings.py` with `sentence-transformers`.
- [ ] **Python**: Implement `trainer.py` for batch feedback processing.
- [ ] **API**: Verify contract between Go and Python (JSON schemas).

### Phase 3: Frontend (Vite SPA)

- [ ] **Init**: Initialize Vite project in `frontend/` (React + TS).
- [ ] **UI**: Create "Anchor Manager" view (CRUD for context anchors).
- [ ] **UI**: Create "Memory Inspector" view (Search memories, view confidence/versions).
- [ ] **Integration**: Connect Frontend to Go Backend via REST/SSE.

### Phase 4: MCP Integration

- [ ] **MCP**: Verify `remember`, `recall`, `give_feedback` tools with Claude Desktop.
- [ ] **Context**: Ensure anchors are correctly injected into recall context.

### Phase 5: Deployment (GCP)

- [ ] **Infra**: Create Terraform for:
  - Cloud Run (Backend, ML Sidecar, Frontend).
  - Cloud SQL (Postgres 16+ with pgvector and pg_trgm).
  - Memorystore (Redis).
- [ ] **CI/CD**: GitHub Actions to build and push images to Artifact Registry.

---

## Folder Structure & Responsibilities

### `backend/` (Go)

- **Responsibility**: The "Brain". Handles MCP requests, business logic, DB state, and coordinates ML tasks.
- **Key Tech**: Go 1.26, `mcp-go`, `pgx`, `viper`.

### `ml-sidecar/` (Python)

- **Responsibility**: The "Subconscious". Handles heavy math, vector generation, and model fine-tuning.
- **Key Tech**: FastAPI, PyTorch, Sentence Transformers.

### `frontend/` (TypeScript)

- **Responsibility**: The "Dashboard". Allows humans to audit the brain, adjust anchors, and correct memories manually.
- **Key Tech**: Vite, React, Tailwind, TanStack Query.

---

## Testing Strategy

We adopt a "Shift Left" testing approach with heavy reliance on containers.

1. **Unit Tests**:
   - Go: Standard `testing` package. Mock DB interfaces.
   - Python: `pytest` with mocked model loading.
   - Frontend: `vitest` for component logic.

2. **Integration Tests (Go)**:
   - Use `testcontainers-go` to spin up ephemeral Postgres and Redis containers.
   - Test the full `recall` pipeline (excluding ML inference, mock the ML client).

3. **End-to-End Tests**:
   - Run the full stack via `docker-compose`.
   - Use **Playwright** to click through the Frontend and verify backend state changes.

---

## Deployment Architecture (GCP)

The system is designed to run serverless on Google Cloud Run.

```mermaid
graph TD
    User[User / Claude] --> LB[Global Load Balancer]
    LB --> Frontend[Cloud Run: Frontend]
    LB --> Backend[Cloud Run: Backend (MCP)]
    Backend --> ML[Cloud Run: ML Sidecar]
    Backend --> DB[(Cloud SQL: Postgres)]
    Backend --> Cache[(Memorystore: Redis)]
```

### Docker Development

For local development, we use `docker-compose` to replicate the cloud environment:

- `backend`: Port 8080
- `ml-sidecar`: Port 8001
- `frontend`: Port 3000
- `postgres`: Port 5432
- `redis`: Port 6379

### Secrets Management

- **Local**: `.env` file (gitignored).
- **GCP**: Google Secret Manager injected as env vars into Cloud Run containers.
