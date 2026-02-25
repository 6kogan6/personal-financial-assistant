package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"personal-financial-assistant/backend/internal/categories"
)

type CategoriesHandlers struct {
	repo *categories.Repo
}

func NewCategoriesHandlers(repo *categories.Repo) *CategoriesHandlers {
	return &CategoriesHandlers{repo: repo}
}

type categoryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type listCategoriesResp struct {
	Items []categoryDTO `json:"items"`
}

func (h *CategoriesHandlers) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "нужен токен")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	items, err := h.repo.ListForUser(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	resp := listCategoriesResp{Items: make([]categoryDTO, 0, len(items))}
	for _, c := range items {
		resp.Items = append(resp.Items, categoryDTO{
			ID:   c.ID,
			Name: c.Name,
			Type: c.Type,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

type createCategoryReq struct {
	Name string `json:"name"`
	Type string `json:"type"` // expense|income
}

func (h *CategoriesHandlers) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "нужен токен")
		return
	}

	var req createCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный json")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name обязателен")
		return
	}
	if req.Type != "expense" && req.Type != "income" {
		writeError(w, http.StatusBadRequest, "type должен быть expense или income")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	c, err := h.repo.CreateUserCategory(ctx, userID, req.Name, req.Type)
	if err != nil {
		if errors.Is(err, categories.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "категория уже существует")
			return
		}
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	writeJSON(w, http.StatusCreated, categoryDTO{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
	})
}
