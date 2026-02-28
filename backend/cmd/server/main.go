package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

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

	// Parse command line flags
	envFile := flag.String("env", "", "Path to .env file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*envFile)
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		log.Fatalf("Configuration loading failed: %v", err)
	}
	logger.Info("Configuration loaded",
		"okta_client_id", cfg.Auth.ClientID,
		"okta_domain", cfg.Auth.OktaDomain,
		"secret_len", len(cfg.Auth.ClientSecret),
		"swagger_client_id", cfg.Auth.SwaggerClientID,
		"config_file", viper.ConfigFileUsed(),
	)

	if cfg.Auth.SwaggerClientID == cfg.Auth.ClientID {
		logger.Warn("⚠️  Swagger Client ID matches Backend Client ID. This will fail if Backend is a Web App (requires secret) and Swagger uses PKCE (no secret). Check your config.yaml.")
	}

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

	// Create Echo server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize authentication
	authz, err := auth.New(ctx, cfg)
	if err != nil {
		logger.Error("failed to initialize auth", "error", err)
		log.Fatalf("auth initialization failed: %v", err)
	}

	// Register auth handlers
	e.GET("/login", echo.WrapHandler(http.HandlerFunc(authz.LoginHandler)))
	e.GET("/auth/callback", echo.WrapHandler(http.HandlerFunc(authz.CallbackHandler)))
	e.GET("/logout", echo.WrapHandler(http.HandlerFunc(authz.LogoutHandler)))

	// Mount REST API handlers
	// Create a group for /api/v1 to match OpenAPI spec and apply auth middleware
	apiGroup := e.Group("/api/v1")
	apiGroup.Use(echo.WrapMiddleware(authz.RequireAuth))
	apiHandler := api.NewHandler()
	api.RegisterHandlers(apiGroup, apiHandler)

	logger.Info("REST API handlers mounted")

	// Mount MCP protocol handlers
	mcpServer := mcp.NewServer(memoryService)
	mcpHandlers := http.NewServeMux()
	mcp.MountHTTPHandlers(mcpHandlers, mcpServer.GetMCPServer())
	e.Any("/mcp/*", echo.WrapHandler(mcpHandlers))

	logger.Info("MCP protocol handlers mounted")

	// expose OpenAPI spec (with runtime substitution) and Swagger UI
	e.GET("/openapi.yaml", echo.WrapHandler(http.HandlerFunc(api.SpecHandler(cfg.Auth.OktaDomain))))
	e.GET("/docs", echo.WrapHandler(http.HandlerFunc(api.SwaggerHandler(cfg.Auth.OktaDomain, cfg.Auth.SwaggerClientID))))
	e.GET("/docs/oauth2-redirect.html", echo.WrapHandler(api.OAuth2RedirectHandler()))

	// Create HTTP server
	addr := ":8080"
	if cfg.TLS.Enable {
		// use TLS port 8443
		addr = ":8443"
	}
	server := &http.Server{
		Addr:         addr,
		Handler:      e,
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
		if err != http.ErrServerClosed {
			logger.Error("Server error: %v", err)
			log.Fatalf("Server error: %v", err)
		}
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
