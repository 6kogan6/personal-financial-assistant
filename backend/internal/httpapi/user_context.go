package httpapi

import (
	"context"
	"errors"
	"strings"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

var errUserIDMissingInContext = errors.New("user id missing in context")

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, strings.TrimSpace(userID))
}

// Старый публичный контракт для существующих handlers:
// userID, ok := UserIDFromContext(ctx)
func UserIDFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(userIDContextKey)

	userID, ok := value.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", false
	}

	return userID, true
}

// Новый внутренний контракт для новых handlers:
// userID, err := userIDFromContext(ctx)
func userIDFromContext(ctx context.Context) (string, error) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return "", errUserIDMissingInContext
	}
	return userID, nil
}

func withUserID(ctx context.Context, userID string) context.Context {
	return WithUserID(ctx, userID)
}
