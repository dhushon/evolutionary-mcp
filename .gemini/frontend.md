# Gemini Style Guide — Frontend

This document outlines the architectural principles and conventions for the frontend development of the Evolutionary Memory MCP project.

---

## 1. Data Fetching & State Management

- **Use TanStack Query**: All server-side state (data fetching, caching, and synchronization) must be managed by `@tanstack/react-query`. This replaces manual `useEffect` and `useState` for data fetching.
- **Centralized API Logic**: API calls should be encapsulated in dedicated functions (e.g., `getWorkflows()`) and organized in the `frontend/src/api` directory.

## 2. API Client

- **Use Axios**: All HTTP requests to the backend must be made through a pre-configured `axios` instance. This provides a robust foundation for handling features like request cancellation, interceptors, and future needs like SSE.
- **Base URL**: The Axios instance should be configured with the base API path (`/api/v1`) to keep individual API calls clean.

## 3. Styling & Theming

- **Tailwind CSS**: The UI is built exclusively with Tailwind CSS.
- **Theming Strategy**:
  - **Dark/Light Mode**: The application must support both modes, controlled by a `dark` class on the `<html>` or `<body>` tag.
  - **CSS Variables**: Core theme colors (primary, secondary, etc.) should be defined as CSS variables in `src/index.css`. This allows for dynamic theme switching (e.g., to match a client's color palette) by updating the variables at runtime.

## 4. Authentication

- **SPA OAuth Flow**: The frontend is a public client and must use the OAuth 2.0 Authorization Code Flow with PKCE.
- **Authentication State**: The user's authentication status is determined by the presence and validity of the `id_token` cookie, which is managed by the backend.
- **API Requests**: As a same-origin application (via Vite proxy), the browser will automatically include the `id_token` cookie in all API requests made by Axios. No manual `Authorization` header is needed for standard user interactions.
