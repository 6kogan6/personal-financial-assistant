package httpapi

import (
	"log/slog"
	"net/http"
	"strings"

	"personal-financial-assistant/backend/internal/apperrors"
	"personal-financial-assistant/backend/internal/auth"
)

func AuthMiddleware(jwtManager *auth.Manager, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				writeError(w, logger, apperrors.New(apperrors.CodeUnauthorized, "missing authorization header"))
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				writeError(w, logger, apperrors.New(apperrors.CodeUnauthorized, "invalid authorization header"))
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if tokenString == "" {
				writeError(w, logger, apperrors.New(apperrors.CodeUnauthorized, "empty bearer token"))
				return
			}

			claims, err := jwtManager.Parse(tokenString)
			if err != nil {
				writeError(w, logger, apperrors.New(apperrors.CodeUnauthorized, "invalid token"))
				return
			}

			ctx := withUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
