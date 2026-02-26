package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"personal-financial-assistant/backend/internal/transactions"
)

type TransactionsHandlers struct {
	repo *transactions.Repo
}

func NewTransactionsHandlers(repo *transactions.Repo) *TransactionsHandlers {
	return &TransactionsHandlers{repo: repo}
}

type transactionDTO struct {
	ID          string  `json:"id"`
	OccurredAt  string  `json:"occurred_at"`
	Type        string  `json:"type"`
	AmountCents int64   `json:"amount_cents"`
	CategoryID  string  `json:"category_id"`
	Merchant    string  `json:"merchant"`
	Note        *string `json:"note,omitempty"`
	Source      string  `json:"source"`
	CreatedAt   string  `json:"created_at,omitempty"`
}

type listTransactionsResp struct {
	Items []transactionDTO `json:"items"`
	Total int              `json:"total"`
}

func (h *TransactionsHandlers) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "нужен токен")
		return
	}

	f := transactions.ListFilter{
		From:       r.URL.Query().Get("from"),
		To:         r.URL.Query().Get("to"),
		CategoryID: r.URL.Query().Get("category_id"),
		Q:          r.URL.Query().Get("q"),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	items, total, err := h.repo.List(ctx, userID, f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	resp := listTransactionsResp{Items: make([]transactionDTO, 0, len(items)), Total: total}
	for _, t := range items {
		dto := transactionDTO{
			ID:          t.ID,
			OccurredAt:  t.OccurredAt,
			Type:        t.Type,
			AmountCents: t.AmountCents,
			CategoryID:  t.CategoryID,
			Merchant:    t.Merchant,
			Note:        t.Note,
			Source:      t.Source,
			CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
		}
		resp.Items = append(resp.Items, dto)
	}

	writeJSON(w, http.StatusOK, resp)
}

type createTransactionReq struct {
	OccurredAt  string  `json:"occurred_at"`
	Type        string  `json:"type"`
	AmountCents any     `json:"amount_cents"` // чтобы ловить и числа и строки
	CategoryID  string  `json:"category_id"`
	Merchant    string  `json:"merchant"`
	Note        *string `json:"note"`
}

func parseAmountCents(v any) (int64, bool) {
	switch x := v.(type) {
	case float64:
		// JSON numbers decode as float64
		if x < 0 {
			return 0, false
		}
		return int64(x), true
	case string:
		n, err := strconv.ParseInt(x, 10, 64)
		if err != nil || n < 0 {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

func (h *TransactionsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "нужен токен")
		return
	}

	var req createTransactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный json")
		return
	}

	if req.OccurredAt == "" || req.Type == "" || req.CategoryID == "" || req.Merchant == "" {
		writeError(w, http.StatusBadRequest, "occurred_at, type, category_id, merchant обязательны")
		return
	}
	if req.Type != "expense" && req.Type != "income" {
		writeError(w, http.StatusBadRequest, "type должен быть expense или income")
		return
	}
	amount, okAmt := parseAmountCents(req.AmountCents)
	if !okAmt {
		writeError(w, http.StatusBadRequest, "amount_cents должен быть целым числом >= 0")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	allowed, err := h.repo.CategoryAllowed(ctx, userID, req.CategoryID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}
	if !allowed {
		writeError(w, http.StatusNotFound, "категория не найдена")
		return
	}

	created, err := h.repo.Create(ctx, transactions.Transaction{
		UserID:      userID,
		OccurredAt:  req.OccurredAt,
		Type:        req.Type,
		AmountCents: amount,
		CategoryID:  req.CategoryID,
		Merchant:    req.Merchant,
		Note:        req.Note,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	writeJSON(w, http.StatusCreated, transactionDTO{
		ID:          created.ID,
		OccurredAt:  created.OccurredAt,
		Type:        created.Type,
		AmountCents: created.AmountCents,
		CategoryID:  created.CategoryID,
		Merchant:    created.Merchant,
		Note:        created.Note,
		Source:      created.Source,
		CreatedAt:   created.CreatedAt.UTC().Format(time.RFC3339),
	})
}
