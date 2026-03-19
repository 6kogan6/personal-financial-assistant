package httpapi

import (
	"context"
	"errors"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

var errUserIDMissingInContext = errors.New("user id missing in context")

func withUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func userIDFromContext(ctx context.Context) (int64, error) {
	value := ctx.Value(userIDContextKey)

	userID, ok := value.(int64)
	if !ok || userID <= 0 {
		return 0, errUserIDMissingInContext
	}

	return userID, nil
}
