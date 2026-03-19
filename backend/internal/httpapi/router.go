package httpapi

import (
	"net/http"
	"time"

	"personal-financial-assistant/backend/internal/auth"
	"personal-financial-assistant/backend/internal/categories"
	"personal-financial-assistant/backend/internal/reports"
	"personal-financial-assistant/backend/internal/transactions"
	"personal-financial-assistant/backend/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	UsersRepo        *users.Repo
	CategoriesRepo   *categories.Repo
	JWT              *auth.Manager
	TransactionsRepo *transactions.Repo
	ReportsRepo      *reports.Repo
}

func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler)

		authHandlers := NewAuthHandlers(deps.UsersRepo, deps.JWT)
		r.Post("/auth/register", authHandlers.Register)
		r.Post("/auth/login", authHandlers.Login)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(deps.JWT))

			profileHandlers := NewProfileHandlers(deps.UsersRepo)
			r.Get("/me", profileHandlers.Me)

			catHandlers := NewCategoriesHandlers(deps.CategoriesRepo)
			r.Get("/categories", catHandlers.List)
			r.Post("/categories", catHandlers.Create)

			txService := transactions.NewService(deps.TransactionsRepo)
			txHandlers := NewTransactionsHandlers(txService)

			r.Get("/transactions", txHandlers.List)
			r.Post("/transactions", txHandlers.Create)
			r.Put("/transactions/{transactionID}", txHandlers.Update)
			r.Delete("/transactions/{transactionID}", txHandlers.Delete)

			repHandlers := NewReportsHandlers(deps.ReportsRepo)
			r.Get("/reports/summary", repHandlers.Summary)
		})
	})

	return r
}
