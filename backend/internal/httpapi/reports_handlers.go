package httpapi

import (
	"context"
	"net/http"
	"time"

	"personal-financial-assistant/backend/internal/reports"
)

type ReportsHandlers struct {
	repo *reports.Repo
}

func NewReportsHandlers(repo *reports.Repo) *ReportsHandlers {
	return &ReportsHandlers{repo: repo}
}

func (h *ReportsHandlers) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "нужен токен")
		return
	}

	month := r.URL.Query().Get("month")
	if month == "" {
		writeError(w, http.StatusBadRequest, "month обязателен (YYYY-MM)")
		return
	}

	// Парсим YYYY-MM
	t, err := time.Parse("2006-01", month)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат month, нужен YYYY-MM")
		return
	}

	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	s, err := h.repo.Summary(ctx, userID, start, end, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	writeJSON(w, http.StatusOK, s)
}
