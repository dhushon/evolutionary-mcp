package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"evolutionary-mcp/backend/internal/config"
	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Logger defines the logging interface compatible with the application logger.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// Auth contains configuration and helpers for performing OpenID Connect
// authentication with an Okta tenant.
type Auth struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	apiVerifier  *oidc.IDTokenVerifier
	repo         repository.Repository
	logger       Logger
	devMode      bool
	authBypass   bool
}

// New creates a new Auth object using values from the application
// configuration. It establishes a connection to the provider and prepares an
// ID token verifier.
func New(ctx context.Context, cfg *config.Config, repo repository.Repository, logger Logger) (*Auth, error) {
	isDev := strings.ToUpper(cfg.Environment) == "DEV"
	shouldBypass := isDev && cfg.DevModeBypass

	var oauth2Config *oauth2.Config
	var verifier *oidc.IDTokenVerifier
	var apiVerifier *oidc.IDTokenVerifier

	if !shouldBypass {
		if cfg.Auth.OktaDomain == "" || cfg.Auth.ClientID == "" ||
			cfg.Auth.ClientSecret == "" || cfg.Auth.RedirectURL == "" {
			return nil, errors.New("auth configuration is incomplete")
		}

		provider, err := oidc.NewProvider(ctx, cfg.Auth.OktaDomain)
		if err != nil {
			return nil, err
		}

		oauth2Config = &oauth2.Config{
			ClientID:     cfg.Auth.ClientID,
			ClientSecret: cfg.Auth.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  cfg.Auth.RedirectURL,
			Scopes:       []string{ScopeOpenID},
		}

		verifier = provider.Verifier(&oidc.Config{ClientID: cfg.Auth.ClientID})

		// Create a separate verifier for Access Tokens (Bearer).
		// We skip ClientID check because Access Tokens often have a different audience (e.g. "api://default")
		apiVerifier = provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	}

	return &Auth{
		oauth2Config: oauth2Config,
		verifier:     verifier,
		apiVerifier:  apiVerifier,
		repo:         repo,
		logger:       logger,
		devMode:      isDev,
		authBypass:   shouldBypass,
	}, nil
}

// LoginHandler initiates the OAuth2 authorization code flow by redirecting the
// user to the Okta authorization endpoint. A random state value is stored in a
// cookie to mitigate CSRF attacks.
func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if a.authBypass {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "failed to generate state", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		// For production you should set Secure: true and SameSite=strict
	})

	http.Redirect(w, r, a.oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// CallbackHandler handles the redirect back from Okta. It verifies the state
// parameter, exchanges the code for tokens, validates the ID token, and sets a
// session cookie containing the raw ID token.
func (a *Auth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if a.authBypass {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// verify state
	cookie, err := r.Cookie("oauthstate")
	if err != nil || r.URL.Query().Get("state") != cookie.Value {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	// exchange code for token
	token, err := a.oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token in token response", http.StatusInternalServerError)
		return
	}

	idToken, err := a.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, "failed to verify id token", http.StatusUnauthorized)
		return
	}

	// optionally parse claims (not used here, but could be stored in session)
	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	_ = idToken.Claims(&claims) // ignore error; claims not required for simple flow

	// set session cookie with raw id token
	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    rawIDToken,
		HttpOnly: true,
		Path:     "/",
		// Secure: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RequireAuth is middleware that ensures a valid ID token cookie is present.
// If the token is missing or invalid the user is redirected to the login page.
func (a *Auth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var email string

		if a.authBypass {
			email = "dev@localhost"
		} else {
			var token *oidc.IDToken
			var err error

			// Check for Authorization header first (for Swagger/API clients)
			if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
				rawToken := strings.TrimPrefix(authHeader, "Bearer ")
				token, err = a.apiVerifier.Verify(r.Context(), rawToken)
				if err != nil {
					http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
					return
				}
			} else {
				cookie, err := r.Cookie("id_token")
				if err != nil {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
				token, err = a.verifier.Verify(r.Context(), cookie.Value)
				if err != nil {
					http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
					return
				}
			}

			// Extract claims to identify the user and tenant
			var claims struct {
				Email string `json:"email"`
			}
			if err := token.Claims(&claims); err != nil {
				http.Error(w, "failed to parse token claims", http.StatusUnauthorized)
				return
			}
			email = claims.Email
		}

		// Resolve Tenant ID from Email Domain
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			http.Error(w, "invalid email format in token", http.StatusUnauthorized)
			return
		}
		domain := parts[1]

		// Lookup or Auto-Provision Tenant
		tenant, err := a.repo.GetTenantByDomain(r.Context(), domain)
		if err != nil {
			// Auto-provisioning for Day 1 experience
			tenant = &models.Tenant{Name: domain, Domain: domain}
			if createErr := a.repo.CreateTenant(r.Context(), tenant); createErr != nil {
				if a.logger != nil {
					a.logger.Error("failed to provision tenant", "domain", domain, "error", createErr)
				}
				http.Error(w, "failed to provision tenant: "+createErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Inject tenant_id into context
		ctx := context.WithValue(r.Context(), "tenant_id", tenant.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LogoutHandler clears the session cookie and redirects to the home page.
func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "id_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
