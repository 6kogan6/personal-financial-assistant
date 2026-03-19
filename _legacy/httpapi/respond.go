package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"personal-financial-assistant/backend/internal/apperrors"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if value == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, logger *slog.Logger, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		message = appErr.Message

		switch appErr.Code {
		case apperrors.CodeInvalidArgument:
			status = http.StatusBadRequest
		case apperrors.CodeUnauthorized:
			status = http.StatusUnauthorized
		case apperrors.CodeForbidden:
			status = http.StatusForbidden
		case apperrors.CodeNotFound:
			status = http.StatusNotFound
		case apperrors.CodeConflict:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
			message = "internal server error"
		}
	}

	if status >= http.StatusInternalServerError && logger != nil {
		logger.Error("request failed", slog.Any("error", err))
	}

	writeJSON(w, status, errorResponse{Error: message})
}
