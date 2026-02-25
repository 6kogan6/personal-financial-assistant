package httpapi

import (
	"net/http"

	"personal-financial-assistant/backend/internal/auth"
	"personal-financial-assistant/backend/internal/users"

	"github.com/go-chi/chi/v5"
)

type Deps struct {
	UsersRepo *users.Repo
	JWT       *auth.Manager
}

func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()

	// Public
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler)

		authHandlers := NewAuthHandlers(deps.UsersRepo, deps.JWT)
		r.Post("/auth/register", authHandlers.Register)
		r.Post("/auth/login", authHandlers.Login)

		// Protected group — пригодится для следующих задач
		// r.Group(func(r chi.Router) {
		// 	r.Use(RequireAuth(deps.JWT))
		// 	// тут будут categories/transactions/reports
		// })
	})

	return r
}
