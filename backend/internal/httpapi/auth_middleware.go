package httpapi

import (
	"context"
	"net/http"
	"strings"

	"personal-financial-assistant/backend/internal/auth"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	s, ok := v.(string)
	return s, ok && s != ""
}

func RequireAuth(jwtMgr *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "нужен токен")
				return
			}

			token := strings.TrimPrefix(h, "Bearer ")
			userID, err := jwtMgr.ParseToken(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "неверный токен")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
