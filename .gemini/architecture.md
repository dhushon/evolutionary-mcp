# System Architecture

## Component Interaction
The system operates as a hybrid local/containerized environment:

1.  **Go Backend (Host)**: Runs via `go run ./cmd/server`. Connects to Postgres, Redis, and the ML Sidecar via `localhost`.
2.  **Python ML Sidecar (Docker)**: Runs in a container, exposing port `8001`. Handles vector embeddings and semantic similarity.
3.  **Postgres/Redis (Docker)**: Persistent and short-term memory stores.
4.  **React Frontend (Host)**: Vite-based SPA. Proxies `/api` requests to the Go backend.

## Data Flow
- **Ingestion**: Frontend -> Go Backend -> ML Sidecar (Embeddings) -> Postgres (pgvector).
- **Retrieval**: AI Assistant -> MCP Server -> Go Backend -> ML Sidecar (Query Embedding) -> Postgres (Vector Search) -> Redis (Context Cache).

## Networking
- **Development**: Uses a "Hybrid" model where infrastructure is containerized but application code runs on the host for faster iteration.
- **Production**: Target is Google Cloud Run with a Global Load Balancer handling TLS.

## Security
- **Auth**: Okta OIDC with a dual-client setup (Confidential for Backend, Public/PKCE for SPA).
- **Authorization**: Tenant-based isolation enforced via `tenant_id` in the request context.