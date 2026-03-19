package transactions

import (
	"context"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

type Repository interface {
	CategoryAllowed(ctx context.Context, userID, categoryID string) (bool, error)
	Create(ctx context.Context, t Transaction) (Transaction, error)
	List(ctx context.Context, userID string, f ListFilter) ([]Transaction, int, error)
	Update(ctx context.Context, t Transaction) (Transaction, error)
	Delete(ctx context.Context, userID, transactionID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	UserID      string
	OccurredAt  string
	Type        string
	AmountCents int64
	CategoryID  string
	Merchant    string
	Note        *string
}

type UpdateInput struct {
	ID          string
	UserID      string
	OccurredAt  string
	Type        string
	AmountCents int64
	CategoryID  string
	Merchant    string
	Note        *string
}

type ListInput struct {
	UserID     string
	From       string
	To         string
	CategoryID string
	Q          string
	Limit      int
	Offset     int
}

func (s *Service) Create(ctx context.Context, in CreateInput) (Transaction, error) {
	in.UserID = strings.TrimSpace(in.UserID)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.Type = normalizeType(in.Type)
	in.Merchant = strings.TrimSpace(in.Merchant)
	in.Note = normalizeNote(in.Note)

	if in.UserID == "" {
		return Transaction{}, ErrInvalidUserID
	}
	if in.CategoryID == "" {
		return Transaction{}, ErrInvalidCategoryID
	}
	if in.Type != "income" && in.Type != "expense" {
		return Transaction{}, ErrInvalidType
	}
	if in.AmountCents <= 0 {
		return Transaction{}, ErrInvalidAmount
	}
	if in.Merchant == "" {
		return Transaction{}, ErrInvalidMerchant
	}

	occurredAt, err := normalizeDateOrToday(in.OccurredAt)
	if err != nil {
		return Transaction{}, err
	}

	ok, err := s.repo.CategoryAllowed(ctx, in.UserID, in.CategoryID)
	if err != nil {
		return Transaction{}, err
	}
	if !ok {
		return Transaction{}, ErrCategoryNotAllowed
	}

	return s.repo.Create(ctx, Transaction{
		UserID:      in.UserID,
		OccurredAt:  occurredAt,
		Type:        in.Type,
		AmountCents: in.AmountCents,
		CategoryID:  in.CategoryID,
		Merchant:    in.Merchant,
		Note:        in.Note,
	})
}

func (s *Service) List(ctx context.Context, in ListInput) ([]Transaction, int, error) {
	in.UserID = strings.TrimSpace(in.UserID)
	in.From = strings.TrimSpace(in.From)
	in.To = strings.TrimSpace(in.To)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.Q = strings.TrimSpace(in.Q)

	if in.UserID == "" {
		return nil, 0, ErrInvalidUserID
	}
	if in.From != "" && !isValidDate(in.From) {
		return nil, 0, ErrInvalidFromDate
	}
	if in.To != "" && !isValidDate(in.To) {
		return nil, 0, ErrInvalidToDate
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	offset := in.Offset
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, in.UserID, ListFilter{
		From:       in.From,
		To:         in.To,
		CategoryID: in.CategoryID,
		Q:          in.Q,
		Limit:      limit,
		Offset:     offset,
	})
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (Transaction, error) {
	in.ID = strings.TrimSpace(in.ID)
	in.UserID = strings.TrimSpace(in.UserID)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.Type = normalizeType(in.Type)
	in.Merchant = strings.TrimSpace(in.Merchant)
	in.Note = normalizeNote(in.Note)

	if in.ID == "" {
		return Transaction{}, ErrInvalidTransactionID
	}
	if in.UserID == "" {
		return Transaction{}, ErrInvalidUserID
	}
	if in.CategoryID == "" {
		return Transaction{}, ErrInvalidCategoryID
	}
	if in.Type != "income" && in.Type != "expense" {
		return Transaction{}, ErrInvalidType
	}
	if in.AmountCents <= 0 {
		return Transaction{}, ErrInvalidAmount
	}
	if in.Merchant == "" {
		return Transaction{}, ErrInvalidMerchant
	}
	if !isValidDate(in.OccurredAt) {
		return Transaction{}, ErrInvalidOccurredAt
	}

	ok, err := s.repo.CategoryAllowed(ctx, in.UserID, in.CategoryID)
	if err != nil {
		return Transaction{}, err
	}
	if !ok {
		return Transaction{}, ErrCategoryNotAllowed
	}

	return s.repo.Update(ctx, Transaction{
		ID:          in.ID,
		UserID:      in.UserID,
		OccurredAt:  in.OccurredAt,
		Type:        in.Type,
		AmountCents: in.AmountCents,
		CategoryID:  in.CategoryID,
		Merchant:    in.Merchant,
		Note:        in.Note,
	})
}

func (s *Service) Delete(ctx context.Context, userID, transactionID string) error {
	userID = strings.TrimSpace(userID)
	transactionID = strings.TrimSpace(transactionID)

	if userID == "" {
		return ErrInvalidUserID
	}
	if transactionID == "" {
		return ErrInvalidTransactionID
	}

	return s.repo.Delete(ctx, userID, transactionID)
}

func normalizeType(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func normalizeNote(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}

func normalizeDateOrToday(v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Now().UTC().Format(dateLayout), nil
	}
	if !isValidDate(v) {
		return "", ErrInvalidOccurredAt
	}
	return v, nil
}

func isValidDate(v string) bool {
	_, err := time.Parse(dateLayout, v)
	return err == nil
}
