# Gemini Style Guide — Evolutionary Memory MCP

## Project Overview

This project is an **Evolutionary Memory MCP (Model Context Protocol) server**. Your primary goal is to assist in its development by providing and maintaining capabilities for short-term and long-term memory, feedback-driven learning, and context anchoring for AI assistants.

The architecture is a **Go/Python hybrid**. Adhere to this separation of concerns:
- **Go**: Handle the MCP server, API routing, short-term memory (Redis), and long-term storage CRUD (Postgres).
- **Python**: Handle all ML inference, including embeddings, semantic similarity, and feedback-driven model adaptation.

---

## Repository Structure

Familiarize yourself with and strictly follow this repository structure. When adding new files, place them in the appropriate directory.

```
mcp-memory/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point; wires dependencies, starts MCP server
├── internal/
│   ├── mcp/
│   │   ├── server.go                # MCP protocol handler (stdio or SSE transport)
│   │   └── tools.go                 # Tool registration, input schemas, dispatch logic
│   ├── memory/
│   │   ├── shortterm.go             # In-process or Redis short-term store with TTL
│   │   ├── longterm.go              # Postgres CRUD for memories, anchors, feedback
│   │   └── retrieval.go             # Hybrid retrieval: vector search + anchor boosting + session merge
│   ├── feedback/
│   │   ├── ingestion.go             # Capture feedback events, validate, persist
│   │   └── reinforcement.go         # Apply confidence delta to memories based on feedback
│   ├── anchor/
│   │   └── context.go               # Manage context anchors; inject into every retrieval
│   ├── mlclient/
│   │   └── client.go                # HTTP client to Python ML sidecar (embed, similarity, train)
│   └── config/
│       └── config.go                # Env-driven config (DB DSN, Redis URL, ML sidecar URL, etc.)
├── migrations/
│   ├── 001_init.up.sql              # memories, feedback_log, context_anchors tables
│   ├── 002_pgvector.up.sql          # Enable pgvector, add embedding column + index
│   └── 003_versioning.up.sql        # memory_versions table for evolution lineage
├── python/
│   ├── main.py                      # FastAPI app entry point
│   ├── embeddings.py                # POST /embed — generate vector from text
│   ├── similarity.py                # POST /similarity — rank candidates by cosine similarity
│   ├── trainer.py                   # POST /feedback_train — batch feedback adaptation
│   ├── requirements.txt
│   └── Dockerfile
├── docker-compose.yml               # Postgres+pgvector, Redis, Python sidecar, Go server
├── Makefile                         # build, test, migrate, run, lint targets
├── .env.example
└── README.md
```

---

## Technology Stack

Utilize the established technology stack. Do not introduce new technologies without explicit instruction.

| Layer | Technology | Notes |
| --- | --- | --- |
| MCP Server | Go 1.26+, `github.com/mark3labs/mcp-go` | Primary runtime |
| Short-term memory | Redis 7+ (or `sync.Map` for local dev) | TTL-scoped per session |
| Long-term memory | PostgreSQL 16+ (with `pgvector` and `pg_trgm` extensions) | Vector + relational in one place |
| Embeddings / ML | Python 3.11, FastAPI, `sentence-transformers` | Swap for OpenAI API if preferred |
| Frontend | React 18, Vite, TailwindCSS, TanStack Query | Dashboard & Management UI |
| Go ↔ Python | HTTP REST (JSON) | gRPC acceptable for high-throughput |
| DB migrations | `golang-migrate/migrate` | SQL files in `/migrations` |
| Config | `github.com/spf13/viper` or plain `os.Getenv` | 12-factor, env-driven |
| Logging | `go.uber.org/zap` (structured JSON) | |
| Testing | `testing` stdlib + `testcontainers-go` | Integration tests spin up real Postgres/Redis |

---

## MCP Tools — Canonical Definitions

When implementing or modifying tools in `internal/mcp/tools.go`, treat these definitions as the source of truth.

### `remember`
- **Goal**: Store information in memory.
- **Inputs**: `content` (string), `memory_type` (enum: `context|preference|decision|fact`), `scope` (enum: `short|long`), `anchor_key` (optional string), `session_id` (string), `metadata` (optional object)
- **Process**: Generate embedding via the ML sidecar, persist to the appropriate store (Redis for short-term, Postgres for long-term), and link to an anchor if provided.

### `recall`
- **Goal**: Retrieve relevant memories.
- **Inputs**: `query` (string), `session_id` (string), `top_k` (int, default 5), `include_anchors` (bool, default true)
- **Process**: Embed the query, perform a similarity search in `pgvector`, boost results associated with anchors, merge with short-term session context, and return ranked results.

### `give_feedback`
- **Goal**: Adjust memory confidence based on user feedback.
- **Inputs**: `memory_id` (UUID), `feedback_type` (enum: `positive|negative|correction`), `correction_text` (optional string), `session_id` (string), `weight` (float, default 1.0)
- **Process**: Log the feedback event, apply the confidence delta, and trigger `evolve_memory` if a correction is provided.

*(... and so on for all other tool definitions ...)*

---

## Database Schema — Canonical

All schema changes **must** be performed by creating new, numbered migration files in `/migrations/`. Do not alter tables directly.

```sql
-- memories
CREATE TABLE memories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    memory_type     TEXT NOT NULL CHECK (memory_type IN ('context','preference','decision','fact','correction')),
    content         TEXT NOT NULL,
    embedding       vector(384),
    metadata        JSONB DEFAULT '{}',
    confidence      FLOAT NOT NULL DEFAULT 0.7,
    -- ... and other columns
);

-- feedback_log
CREATE TABLE feedback_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    memory_id       UUID NOT NULL REFERENCES memories(id),
    -- ... and other columns
);

-- context_anchors
CREATE TABLE context_anchors (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anchor_key      TEXT UNIQUE NOT NULL,
    -- ... and other columns
);

-- memory_versions
CREATE TABLE memory_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    memory_id       UUID NOT NULL REFERENCES memories(id),
    -- ... and other columns
);
```

---

## Key Design Invariants & Rules

You must adhere to these core principles at all times:

1.  **Never Delete Memories**: Mark memories as `deprecated = true` or create new versions. The audit trail is a critical feature.
2.  **Anchors are Always Injected**: Anchors must be included in `recall` results by default.
3.  **Short-Term Memory is Session-Scoped**: `session_id` is mandatory for `remember` and `recall`.
4.  **Embeddings are Centralized**: The Go backend **must not** generate embeddings. It must always call the Python ML sidecar.
5.  **Feedback is Asynchronous**: The `give_feedback` tool must return immediately. Apply confidence updates in a background process.
6.  **Versioning on Correction**: An `evolve_memory` call must always create a snapshot of the prior state.

---

## Coding Conventions

### Go
- **Interfaces for Data Access**: All database access must be through an interface (e.g., `MemoryStore`). Do not use raw `*sql.DB` in business logic. This is for testability.
- **Context Propagation**: Thread `context.Context` through all function calls. Respect cancellation.
- **Error Handling**: Use wrapped errors, e.g., `fmt.Errorf("recall: %w", err)`.

### Python
- **Model Loading**: Load the embedding model once at startup, not on each request.
- **Error Responses**: Return errors as a JSON object: `{"error": "description"}`.
- **Testing**: Use `pytest` and mock the transformer model.

### SQL
- **No SQL Injection**: Always use parameterized queries.
- **Migrations**: Write irreversible `up` migrations.

---

## Testing Strategy

- **Integration tests** (`backend/internal/*/integration_test.go`): Use `testcontainers-go` to spin up Postgres & Redis.
- **Python tests** (`ml-sidecar/tests/`): Mock the `SentenceTransformer`.
- **E2E Tests**: Use Playwright against the full `docker-compose` stack.
- **Run tests using the Makefile**: `make test`
