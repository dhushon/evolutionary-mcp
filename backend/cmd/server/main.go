package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"evolutionary-mcp/backend/internal/api"
	"evolutionary-mcp/backend/internal/auth"
	"evolutionary-mcp/backend/internal/config"
	"evolutionary-mcp/backend/internal/logging"
	"evolutionary-mcp/backend/internal/mcp"
	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/internal/services"
	"evolutionary-mcp/backend/internal/tls"
)

func main() {
	ctx := context.Background()

	// Initialize logging
	logger := logging.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		log.Fatalf("Configuration loading failed: %v", err)
	}
	logger.Info("Configuration loaded", "okta_client_id", cfg.Auth.ClientID, "okta_domain", cfg.Auth.OktaDomain)

	logger.Info("Starting Evolutionary Memory Service")

	// Initialize database connection
	dbPool, err := initDatabase(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize database: %v", err)
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer dbPool.Close()

	logger.Info("Database connected")

	// Initialize repository layer
	memoryStore := repository.NewPostgresMemoryStore(dbPool)

	// Initialize service layer
	mlClient := services.NewHTTPMLClient(cfg.MLSidecar.URL)
	memoryService := services.NewMemoryService(memoryStore, mlClient)

	logger.Info("Service layer initialized")

	// Create HTTP mux
	mux := http.NewServeMux()

	// Initialize authentication
	authz, err := auth.New(ctx, cfg)
	if err != nil {
		logger.Error("failed to initialize auth", "error", err)
		log.Fatalf("auth initialization failed: %v", err)
	}

	// Register auth handlers
	mux.HandleFunc("/login", authz.LoginHandler)
	mux.HandleFunc("/auth/callback", authz.CallbackHandler)
	mux.HandleFunc("/logout", authz.LogoutHandler)

	// Mount REST API handlers
	apiHandler := api.NewHandler()
	mux.HandleFunc("/api/v1/health", apiHandler.HandleHealth)

	logger.Info("REST API handlers mounted")

	// Mount MCP protocol handlers behind authentication
	mcpServer := mcp.NewServer(memoryService)
	mcpHandlers := http.NewServeMux()
	mcp.MountHTTPHandlers(mcpHandlers, mcpServer.GetMCPServer())
	mux.Handle("/mcp/", authz.RequireAuth(mcpHandlers))

	logger.Info("MCP protocol handlers mounted")

	// expose OpenAPI spec (with runtime substitution) and Swagger UI
	mux.HandleFunc("/openapi.yaml", api.SpecHandler(cfg.Auth.OktaDomain))

	mux.HandleFunc("/docs", api.SwaggerHandler(cfg.Auth.OktaDomain, cfg.Auth.ClientID))
	mux.HandleFunc("/docs/oauth2-redirect.html", api.OAuthRedirectHandler)

	// Wrap with logging middleware
	handler := loggingMiddleware(mux, logger)

	// Create HTTP server
	addr := ":8080"
	if cfg.TLS.Enable {
		// use TLS port 8443
		addr = ":8443"
	}
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown handling
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("Server starting", "address", addr, "tls", cfg.TLS.Enable)
		if cfg.TLS.Enable {
			// ensure certificate exists if requested
			if cfg.TLS.CertFile == "" || cfg.TLS.KeyFile == "" {
				logger.Error("TLS enabled but cert/key file not provided")
				httpErr := server.ListenAndServe()
				serverErrors <- httpErr
				return
			}
			// generate if missing and hostnames provided
			if _, err := os.Stat(cfg.TLS.CertFile); os.IsNotExist(err) {
				if len(cfg.TLS.Hostnames) > 0 {
					if err := tls.GenerateSelfSignedCert(cfg.TLS.CertFile, cfg.TLS.KeyFile, cfg.TLS.Hostnames); err != nil {
						logger.Error("failed to generate self-signed cert", "error", err)
					}
				}
			}
			serverErrors <- server.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			serverErrors <- server.ListenAndServe()
		}
	}()

	// Wait for shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("Server error: %v", err)
		log.Fatalf("Server error: %v", err)
	case sig := <-shutdown:
		logger.Info("Shutdown signal received: %v", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown error: %v", err)
			if err := server.Close(); err != nil {
				logger.Error("Server close error: %v", err)
			}
		}

		logger.Info("Server stopped gracefully")
	}
}

func initDatabase(ctx context.Context, cfg *config.Config, logger *logging.Logger) (*pgxpool.Pool, error) {
	logger.Debug("Initializing database connection")

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func loggingMiddleware(next http.Handler, logger *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call next handler
		next.ServeHTTP(w, r)

		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"duration", time.Since(start).String(),
		)
	})
}
