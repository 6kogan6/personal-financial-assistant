package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"personal-financial-assistant/backend/internal/auth"
	"personal-financial-assistant/backend/internal/users"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	users *users.Repo
	jwt   *auth.Manager
}

func NewAuthHandlers(usersRepo *users.Repo, jwtMgr *auth.Manager) *AuthHandlers {
	return &AuthHandlers{users: usersRepo, jwt: jwtMgr}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResp struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный json")
		return
	}
	if req.Email == "" || len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "email обязателен, пароль минимум 8 символов")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	u, err := h.users.Create(ctx, req.Email, string(hash))
	if err != nil {
		if errors.Is(err, users.ErrEmailExists) {
			writeError(w, http.StatusConflict, "email уже существует")
			return
		}
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	writeJSON(w, http.StatusCreated, registerResp{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный json")
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email и пароль обязательны")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	userID, passwordHash, err := h.users.GetAuthDataByEmail(ctx, req.Email)
	if err != nil {
		// не палим, существует ли email
		writeError(w, http.StatusUnauthorized, "неверный email/пароль")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "неверный email/пароль")
		return
	}

	token, err := h.jwt.NewToken(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	writeJSON(w, http.StatusOK, loginResp{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}
