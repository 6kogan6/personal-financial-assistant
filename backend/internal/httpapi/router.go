package httpapi

import (
	"net/http"

	"personal-financial-assistant/backend/internal/auth"
	"personal-financial-assistant/backend/internal/categories"
	"personal-financial-assistant/backend/internal/users"

	"github.com/go-chi/chi/v5"
)

type Deps struct {
	UsersRepo      *users.Repo
	CategoriesRepo *categories.Repo
	JWT            *auth.Manager
}

func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		// Public
		r.Get("/health", healthHandler)

		authHandlers := NewAuthHandlers(deps.UsersRepo, deps.JWT)
		r.Post("/auth/register", authHandlers.Register)
		r.Post("/auth/login", authHandlers.Login)

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(deps.JWT))

			catHandlers := NewCategoriesHandlers(deps.CategoriesRepo)
			r.Get("/categories", catHandlers.List)
			r.Post("/categories", catHandlers.Create)
		})
	})

	return r
}
