# Design Choices

## 1. Frontend Layout Strategy: Shell-and-Component
- **Application Shell**: A centralized `Shell.tsx` component manages the responsive layout. It handles the sidebar (hidden on mobile, visible on desktop) and the main content flow.
- **Data Consistency**: Standardized `DataCard.tsx` is used for all primary data entities (Memories, Workflows). This ensures confidence scores, status badges, and metadata are rendered uniformly.
- **Theming**: Core colors and semantic backgrounds are defined as CSS variables in `index.css` (e.g., `--bg-base`, `--primary`). This supports seamless Dark/Light mode transitions using Tailwind's `dark:` variants.

## 2. Backend Architecture
- **TLS Termination**: Removed internal TLS logic from the Go backend. TLS will be handled by an Nginx router or Cloud Load Balancer in production. The backend now listens on plain HTTP (port 8080).
- **Context Bridging**: Implemented a custom middleware in `main.go` to bridge the standard `http.Request` context to the Echo context. This ensures that `tenant_id` (injected during auth bypass) is accessible to all API handlers via `c.Get("tenant_id")`.
- **ML Sidecar Integration**: The Python ML service is treated as a core infrastructure dependency. It is started automatically via `make run` using Docker Compose to ensure embedding capabilities are always available during development.

## 3. Project Structure & Hygiene
- **Case Sensitivity**: Standardized on PascalCase for React components (e.g., `Sidebar.tsx`). Duplicate lowercase files (e.g., `sidebar.tsx`) were removed to prevent build conflicts on case-insensitive filesystems.
- **Configuration**: Centralized environment variables in a root `.env` file, managed by a Go-based setup utility.

## 4. Authentication
- **Dev Mode Bypass**: Supports a `DEV_MODE_BYPASS` flag. When enabled, the `RequireAuth` middleware injects a mock `tenant_id` into the context, allowing development without an active Okta session.