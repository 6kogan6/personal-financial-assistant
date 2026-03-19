package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"personal-financial-assistant/backend/internal/apperrors"
	"personal-financial-assistant/backend/internal/transactions"
)

type TransactionsHandler struct {
	service *transactions.Service
	logger  *slog.Logger
}

func NewTransactionsHandler(service *transactions.Service, logger *slog.Logger) *TransactionsHandler {
	return &TransactionsHandler{
		service: service,
		logger:  logger,
	}
}

type createTransactionRequest struct {
	CategoryID  int64     `json:"category_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type updateTransactionRequest struct {
	CategoryID  int64     `json:"category_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func (h *TransactionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, h.logger, apperrors.New(apperrors.CodeUnauthorized, "unauthorized"))
		return
	}

	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid json body"))
		return
	}

	item, err := h.service.Create(r.Context(), userID, transactions.CreateInput{
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
		OccurredAt:  req.OccurredAt,
	})
	if err != nil {
		writeError(w, h.logger, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *TransactionsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, h.logger, apperrors.New(apperrors.CodeUnauthorized, "unauthorized"))
		return
	}

	query := r.URL.Query()

	var dateFrom *time.Time
	if value := query.Get("date_from"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid date_from"))
			return
		}
		dateFrom = &parsed
	}

	var dateTo *time.Time
	if value := query.Get("date_to"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid date_to"))
			return
		}
		dateTo = &parsed
	}

	var categoryID int64
	if value := query.Get("category_id"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed <= 0 {
			writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid category_id"))
			return
		}
		categoryID = parsed
	}

	limit := 50
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid limit"))
			return
		}
		limit = parsed
	}

	offset := 0
	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid offset"))
			return
		}
		offset = parsed
	}

	items, err := h.service.List(r.Context(), userID, transactions.ListInput{
		Type:       query.Get("type"),
		CategoryID: categoryID,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		writeError(w, h.logger, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *TransactionsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, h.logger, apperrors.New(apperrors.CodeUnauthorized, "unauthorized"))
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "transactionID"), 10, 64)
	if err != nil || transactionID <= 0 {
		writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid transaction id"))
		return
	}

	var req updateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid json body"))
		return
	}

	item, err := h.service.Update(r.Context(), userID, transactions.UpdateInput{
		ID:          transactionID,
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
		OccurredAt:  req.OccurredAt,
	})
	if err != nil {
		writeError(w, h.logger, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *TransactionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		writeError(w, h.logger, apperrors.New(apperrors.CodeUnauthorized, "unauthorized"))
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "transactionID"), 10, 64)
	if err != nil || transactionID <= 0 {
		writeError(w, h.logger, apperrors.New(apperrors.CodeInvalidArgument, "invalid transaction id"))
		return
	}

	if err := h.service.Delete(r.Context(), userID, transactionID); err != nil {
		writeError(w, h.logger, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
