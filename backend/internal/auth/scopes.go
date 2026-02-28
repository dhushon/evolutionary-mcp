package auth

const (
	ScopeOpenID      = "openid"
	ScopeProfile     = "profile"
	ScopeEmail       = "email"
	ScopeEvolveRead  = "evolve:read"
	ScopeEvolveWrite = "evolve:write"
)

// AllScopes defines the full set of scopes used by the Swagger UI / Frontend
var AllScopes = []string{
	ScopeOpenID,
	ScopeProfile,
	ScopeEmail,
	ScopeEvolveRead,
	ScopeEvolveWrite,
}
