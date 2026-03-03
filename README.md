# Evolutionary Memory MCP Server

## 1. Project Overview

This project is an **Evolutionary Memory MCP (Model Context Protocol) server** designed to provide AI assistants like Claude with robust short-term and long-term memory. It features feedback-driven learning and context anchoring, allowing the system's knowledge base to grow and adapt alongside the business it supports.

The core purpose is to create a persistent, evolving memory layer that enhances AI interactions by:

- **Remembering** key information from conversations.
- **Recalling** relevant context on demand using semantic search.
- **Learning** from user feedback to improve the confidence and accuracy of its memories.
- **Anchoring** critical business concepts to ensure they are always prioritized in recall.

## 2. Design & Architecture

The system uses a powerful and flexible hybrid architecture that separates high-performance request handling from computationally intensive machine learning tasks.

### Core Components

- **Go Backend (`/backend`)**: The "Brain" of the operation. It's a high-performance Go application that handles:
- The MCP server endpoint for the AI assistant.
- API routing and business logic for all memory operations.
- CRUD operations for short-term (Redis) and long-term (PostgreSQL) memory stores.
- Coordination with the Python ML sidecar for embedding and similarity tasks.

- **Python ML Sidecar (`/ml-sidecar`)**: The "Subconscious." This FastAPI service handles all heavy machine learning tasks, keeping the Go backend lean and fast. Its responsibilities include:
- Generating vector embeddings from text for semantic search.
- Calculating semantic similarity between queries and stored memories.
- Handling feedback-driven model adaptation and training jobs.

- **Vite Frontend (`/frontend`)**: The "Dashboard." A React-based single-page application that provides a human-in-the-loop interface for:
- Visualizing memory graphs and confidence scores.
- Managing and auditing context anchors.
- Manually providing feedback or corrections to memories.

### Technology Stack

| Layer | Technology | Purpose |
| --- | --- | --- |
| **MCP Server** | Go 1.26+ | Primary runtime for business logic and API |
| **Short-Term Memory** | Redis 7+ | Session-scoped context with a TTL |
| **Long-Term Memory** | PostgreSQL 16+ (with `pgvector` and `pg_trgm` extensions) | Persistent storage for memories, anchors, and feedback, with vector and fuzzy search capabilities. |
| **Embeddings / ML** | Python 3.11, FastAPI, `sentence-transformers` | Vector generation and semantic similarity |
| **Frontend** | React 18, Vite, TailwindCSS | Dashboard and management UI |
| **Orchestration** | Docker Compose | Local development environment |
| **Deployment** | Docker, Google Cloud Run | Target production environment |

### System Architecture (GCP)

The application is designed to be deployed as a set of containerized, serverless components on Google Cloud.

```mermaid
graph TD
    subgraph "User Interaction"
        User[User / Claude Assistant]
    end

    subgraph "Google Cloud Platform"
        LB[Global Load Balancer]

        subgraph "Cloud Run Services"
            Frontend[Frontend UI]
            Backend[Backend MCP Server]
            ML[ML Sidecar]
        end

        subgraph "Data Stores"
            DB[(Cloud SQL: Postgres + pgvector)]
            Cache[(Memorystore: Redis)]
        end
    end

    User --> LB
    LB --> Frontend
    LB --> Backend
    Backend --> ML
    Backend --> DB
    Backend --> Cache
```

## 3. Getting Started

Follow these steps to set up your local development environment and run the application.

### Prerequisites

- Docker and Docker Compose
- Go (version 1.26 or later)
- Python (version 3.11 or later)
- Node.js and npm

### 1. Clone the Repository

```sh
git clone <repository-url>
cd evolutionary-mcp
```

### 2. Configure Environment Variables

Copy the example environment file and customize it as needed.

```sh
cp .env.example .env
```

At a minimum, ensure the `DATABASE_URL` and `REDIS_URL` point to the services that will be started by Docker Compose. The default values are typically sufficient for local development.

### 3. Build and Run the Application

The easiest way to start the development environment is via the provided `Makefile`.

```sh
# bring up Postgres & Redis, then launch the Go backend locally
make run
```

Under the hood `make run` spins up the `postgres` and `redis` containers and
waits for the database to be ready before running `go run ./backend/cmd/server`.
If Docker or the image pull fails you'll see a helpful error message and the
Makefile will exit; in that case either fix your Docker installation or start
Postgres/Redis manually and ensure the connection parameters in
`config.yaml` or environment variables point to a running instance.

Alternatively you can bring up the full compose stack (including frontend and
ml-sidecar) with:

```sh
docker-compose up --build -d
```

- The Go backend will be available on port `8080`.
- The Python ML sidecar will be available on port `8001`.
- The frontend runs on port `5173`.

### 4. Run the Frontend

In a separate terminal, navigate to the `frontend` directory to install dependencies and start the Vite development server.

```sh
cd frontend
npm install
npm run dev
```

The frontend dashboard will be available at `http://localhost:5173`.

### 5. Verify Installation

Once all services are running, you can connect your AI assistant (e.g., Claude Desktop) to the local MCP server. Configure the assistant to point to the memory server, typically by updating its local MCP settings to use the command for the running backend process.

You can verify the connection by asking the assistant to use one of the memory tools, such as `list_anchors`.

## 6. Authentication (Okta OAuth)

The backend supports user authentication via Okta using a dual-client architecture to support both server-side and browser-based flows securely.

1. **Create Okta Applications**
   - **Backend App**: Create a "Web Application".
     - *Redirect URI*: `http://localhost:8080/auth/callback`
     - Note the *Client ID* and *Client Secret*.
   - **Frontend App**: Create a "Single Page App (SPA)".
     - *Redirect URI*: `http://localhost:8080/docs/oauth2-redirect.html`
     - Enable **PKCE**.
     - Note the *Client ID* (no secret required).

2. **Update Configuration**
   - Add the Okta issuer URL and credentials to `config.yaml`.

     ```yaml
     auth:
       # full issuer URL, e.g. "https://dev-123456.okta.com/oauth2/default"
       okta_domain: "https://dev-123456.okta.com/oauth2/default"
       client_id: "BACKEND_WEB_APP_CLIENT_ID"
       client_secret: "BACKEND_WEB_APP_CLIENT_SECRET"
       # New field for Swagger/SPA:
       swagger_client_id: "FRONTEND_SPA_CLIENT_ID"
       redirect_url: "http://localhost:8080/auth/callback"
     ```

   - The backend will read and normalize these values on startup.

3. **Using the Frontend**
   - Navigate to the frontend UI (e.g. `http://localhost:5173`).
   - Click the **Login with Okta** button. You will be redirected to Okta to sign in.
   - After successful authentication you’ll be returned to the dashboard with a session cookie set.
   - Use the **Logout** button to clear the session.

4. **Swagger / OpenAPI Docs**
   - A Swagger UI is available at `http://localhost:8080/docs`.
   - The UI loads the OpenAPI spec from `/openapi.yaml` and includes an **Authorize** button.
   - The **Swagger Client ID** is injected automatically. It uses PKCE, so no client secret is required or sent.
   - Click **Authorize** to redirect to Okta for login.
   - After authorizing, all requests made from the UI will include the access token.
   - A read‑only text box above the Swagger UI will also show the raw Bearer token; you can copy/paste it to share with other users or for manual API calls.
   - Ensure your Okta application’s allowed redirect URLs include
     `http://localhost:8080/docs/oauth2-redirect.html`.

5. **API Access**
   - Any requests to `/mcp/*` paths require an authenticated user. Unauthenticated requests will be redirected to `/login`.
   - You can test by hitting `http://localhost:8080/api/v1/health` with credentials included; a 200 response indicates a valid session.

## 7. Active Development Tasks (Context for Next Session)

**Current Status:**

- **Auth**: Dual-client architecture (Backend Confidential + SPA Public) is fully configured and documented. Scopes are centralized in `backend/internal/auth/scopes.go`.
- **Config**: Setup utility (`make setup-env`) is working and generates `.env` files.
- **Frontend**: Vite proxy is configured in `frontend/vite.config.ts` to forward `/api`, `/login`, and `/logout` to the backend.

**Next Steps:**

1. **Workflows API**: Create DB migration for `workflows` table and implement `listWorkflows` handler logic in Go.
2. **Frontend Data**: Wire up React components to fetch from `/api/v1/workflows`.
