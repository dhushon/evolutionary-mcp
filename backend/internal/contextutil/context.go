package contextutil

import (
	"context"
)

type contextKey string

const (
	tenantIDKey contextKey = "tenant_id"
	userIDKey   contextKey = "user_id"
)

// WithTenant returns a new context with the tenant ID attached.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// GetTenant extracts the tenant ID from the context.
func GetTenant(ctx context.Context) string {
	val, _ := ctx.Value(tenantIDKey).(string)
	return val
}

// WithUser returns a new context with the user ID attached.
func WithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUser extracts the user ID from the context.
func GetUser(ctx context.Context) string {
	val, _ := ctx.Value(userIDKey).(string)
	return val
}
