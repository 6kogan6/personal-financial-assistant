package httpapi

import (
	"context"
	"errors"
	"net/http"
	"time"

	"personal-financial-assistant/backend/internal/users"
)

type ProfileHandlers struct {
	usersRepo *users.Repo
}

func NewProfileHandlers(usersRepo *users.Repo) *ProfileHandlers {
	return &ProfileHandlers{usersRepo: usersRepo}
}

type profileResp struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *ProfileHandlers) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "нужен токен")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	u, err := h.usersRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "пользователь не найден")
			return
		}
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	writeJSON(w, http.StatusOK, profileResp{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	})
}
