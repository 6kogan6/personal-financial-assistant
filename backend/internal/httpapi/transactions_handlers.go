package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"personal-financial-assistant/backend/internal/transactions"
)

type TransactionsService interface {
	Create(ctx context.Context, in transactions.CreateInput) (transactions.Transaction, error)
	List(ctx context.Context, in transactions.ListInput) ([]transactions.Transaction, int, error)
	Update(ctx context.Context, in transactions.UpdateInput) (transactions.Transaction, error)
	Delete(ctx context.Context, userID, transactionID string) error
}

type TransactionsHandlers struct {
	service TransactionsService
}

func NewTransactionsHandlers(service TransactionsService) *TransactionsHandlers {
	return &TransactionsHandlers{service: service}
}

type createTransactionRequest struct {
	OccurredAt  string  `json:"occurred_at"`
	Type        string  `json:"type"`
	AmountCents int64   `json:"amount_cents"`
	CategoryID  string  `json:"category_id"`
	Merchant    string  `json:"merchant"`
	Note        *string `json:"note"`
}

type updateTransactionRequest struct {
	OccurredAt  string  `json:"occurred_at"`
	Type        string  `json:"type"`
	AmountCents int64   `json:"amount_cents"`
	CategoryID  string  `json:"category_id"`
	Merchant    string  `json:"merchant"`
	Note        *string `json:"note"`
}

type listTransactionsResponse struct {
	Items  []transactions.Transaction `json:"items"`
	Total  int                        `json:"total"`
	Limit  int                        `json:"limit"`
	Offset int                        `json:"offset"`
}

func (h *TransactionsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createTransactionRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный json")
		return
	}

	item, err := h.service.Create(r.Context(), transactions.CreateInput{
		UserID:      userID,
		OccurredAt:  req.OccurredAt,
		Type:        req.Type,
		AmountCents: req.AmountCents,
		CategoryID:  req.CategoryID,
		Merchant:    req.Merchant,
		Note:        req.Note,
	})
	if err != nil {
		writeTransactionsError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *TransactionsHandlers) List(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	query := r.URL.Query()

	limit := 50
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = parsed
	}

	offset := 0
	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = parsed
	}

	items, total, err := h.service.List(r.Context(), transactions.ListInput{
		UserID:     userID,
		From:       query.Get("from"),
		To:         query.Get("to"),
		CategoryID: query.Get("category_id"),
		Q:          query.Get("q"),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		writeTransactionsError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, listTransactionsResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *TransactionsHandlers) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	transactionID := chi.URLParam(r, "transactionID")
	if transactionID == "" {
		writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req updateTransactionRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный json")
		return
	}

	item, err := h.service.Update(r.Context(), transactions.UpdateInput{
		ID:          transactionID,
		UserID:      userID,
		OccurredAt:  req.OccurredAt,
		Type:        req.Type,
		AmountCents: req.AmountCents,
		CategoryID:  req.CategoryID,
		Merchant:    req.Merchant,
		Note:        req.Note,
	})
	if err != nil {
		writeTransactionsError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *TransactionsHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	transactionID := chi.URLParam(r, "transactionID")
	if transactionID == "" {
		writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	if err := h.service.Delete(r.Context(), userID, transactionID); err != nil {
		writeTransactionsError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeTransactionsError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, transactions.ErrInvalidUserID),
		errors.Is(err, transactions.ErrInvalidTransactionID),
		errors.Is(err, transactions.ErrInvalidCategoryID),
		errors.Is(err, transactions.ErrInvalidType),
		errors.Is(err, transactions.ErrInvalidAmount),
		errors.Is(err, transactions.ErrInvalidMerchant),
		errors.Is(err, transactions.ErrInvalidOccurredAt),
		errors.Is(err, transactions.ErrInvalidFromDate),
		errors.Is(err, transactions.ErrInvalidToDate),
		errors.Is(err, transactions.ErrCategoryNotAllowed):
		writeError(w, http.StatusBadRequest, err.Error())
		return

	case errors.Is(err, transactions.ErrTransactionNotFound):
		writeError(w, http.StatusNotFound, err.Error())
		return

	default:
		writeError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}
}
