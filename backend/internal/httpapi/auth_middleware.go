package httpapi

import (
	"net/http"
	"strings"

	"personal-financial-assistant/backend/internal/auth"
)

func AuthMiddleware(jwtManager *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				writeError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if tokenString == "" {
				writeError(w, http.StatusUnauthorized, "empty bearer token")
				return
			}

			claims, err := jwtManager.Parse(tokenString)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := withUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(jwtManager *auth.Manager) func(http.Handler) http.Handler {
	return AuthMiddleware(jwtManager)
}
