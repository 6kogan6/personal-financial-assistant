package transactions

import (
	"context"
	"strings"
	"time"

	"personal-financial-assistant/backend/internal/apperrors"
)

type Repository interface {
	Create(ctx context.Context, params CreateParams) (Transaction, error)
	ListByUser(ctx context.Context, filter ListFilter) ([]Transaction, error)
	Update(ctx context.Context, params UpdateParams) (Transaction, error)
	DeleteByID(ctx context.Context, userID int64, transactionID int64) error
}

type CreateParams struct {
	UserID      int64
	CategoryID  int64
	Type        string
	Amount      float64
	Description string
	OccurredAt  time.Time
}

type UpdateParams struct {
	ID          int64
	UserID      int64
	CategoryID  int64
	Type        string
	Amount      float64
	Description string
	OccurredAt  time.Time
}

type ListFilter struct {
	UserID     int64
	Type       string
	CategoryID int64
	DateFrom   *time.Time
	DateTo     *time.Time
	Limit      int
	Offset     int
}

type CreateInput struct {
	CategoryID  int64
	Type        string
	Amount      float64
	Description string
	OccurredAt  time.Time
}

type UpdateInput struct {
	ID          int64
	CategoryID  int64
	Type        string
	Amount      float64
	Description string
	OccurredAt  time.Time
}

type ListInput struct {
	Type       string
	CategoryID int64
	DateFrom   *time.Time
	DateTo     *time.Time
	Limit      int
	Offset     int
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID int64, input CreateInput) (Transaction, error) {
	if userID <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeUnauthorized, "invalid user id")
	}

	if err := validateType(input.Type); err != nil {
		return Transaction{}, err
	}

	if input.CategoryID <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "category_id must be positive")
	}

	if input.Amount <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "amount must be positive")
	}

	description := normalizeDescription(input.Description)
	if len(description) > 255 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "description is too long")
	}

	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	transaction, err := s.repo.Create(ctx, CreateParams{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Type:        normalizeType(input.Type),
		Amount:      input.Amount,
		Description: description,
		OccurredAt:  occurredAt,
	})
	if err != nil {
		return Transaction{}, apperrors.Wrap(apperrors.CodeInternal, "failed to create transaction", err)
	}

	return transaction, nil
}

func (s *Service) List(ctx context.Context, userID int64, input ListInput) ([]Transaction, error) {
	if userID <= 0 {
		return nil, apperrors.New(apperrors.CodeUnauthorized, "invalid user id")
	}

	if input.Type != "" {
		if err := validateType(input.Type); err != nil {
			return nil, err
		}
	}

	limit := input.Limit
	offset := input.Offset

	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	items, err := s.repo.ListByUser(ctx, ListFilter{
		UserID:     userID,
		Type:       normalizeType(input.Type),
		CategoryID: input.CategoryID,
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "failed to load transactions", err)
	}

	return items, nil
}

func (s *Service) Update(ctx context.Context, userID int64, input UpdateInput) (Transaction, error) {
	if userID <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeUnauthorized, "invalid user id")
	}

	if input.ID <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "transaction id must be positive")
	}

	if err := validateType(input.Type); err != nil {
		return Transaction{}, err
	}

	if input.CategoryID <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "category_id must be positive")
	}

	if input.Amount <= 0 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "amount must be positive")
	}

	description := normalizeDescription(input.Description)
	if len(description) > 255 {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "description is too long")
	}

	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		return Transaction{}, apperrors.New(apperrors.CodeInvalidArgument, "occurred_at is required")
	}

	transaction, err := s.repo.Update(ctx, UpdateParams{
		ID:          input.ID,
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Type:        normalizeType(input.Type),
		Amount:      input.Amount,
		Description: description,
		OccurredAt:  occurredAt,
	})
	if err != nil {
		return Transaction{}, apperrors.Wrap(apperrors.CodeInternal, "failed to update transaction", err)
	}

	return transaction, nil
}

func (s *Service) Delete(ctx context.Context, userID int64, transactionID int64) error {
	if userID <= 0 {
		return apperrors.New(apperrors.CodeUnauthorized, "invalid user id")
	}

	if transactionID <= 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "transaction id must be positive")
	}

	if err := s.repo.DeleteByID(ctx, userID, transactionID); err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "failed to delete transaction", err)
	}

	return nil
}

func validateType(value string) error {
	switch normalizeType(value) {
	case "income", "expense":
		return nil
	default:
		return apperrors.New(apperrors.CodeInvalidArgument, "type must be either income or expense")
	}
}

func normalizeType(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeDescription(value string) string {
	return strings.TrimSpace(value)
}
