package transactions

import (
	"context"
	"errors"
	"testing"
)

type stubRepo struct {
	categoryAllowed bool
	categoryErr     error

	createFn func(ctx context.Context, t Transaction) (Transaction, error)
	listFn   func(ctx context.Context, userID string, f ListFilter) ([]Transaction, int, error)
	updateFn func(ctx context.Context, t Transaction) (Transaction, error)
	deleteFn func(ctx context.Context, userID, transactionID string) error

	lastCreate Transaction
	lastList   struct {
		userID string
		filter ListFilter
	}
	lastUpdate Transaction
	lastDelete struct {
		userID        string
		transactionID string
	}
}

func (r *stubRepo) CategoryAllowed(ctx context.Context, userID, categoryID string) (bool, error) {
	if r.categoryErr != nil {
		return false, r.categoryErr
	}
	return r.categoryAllowed, nil
}

func (r *stubRepo) Create(ctx context.Context, t Transaction) (Transaction, error) {
	r.lastCreate = t
	if r.createFn != nil {
		return r.createFn(ctx, t)
	}
	t.ID = "tx-1"
	return t, nil
}

func (r *stubRepo) List(ctx context.Context, userID string, f ListFilter) ([]Transaction, int, error) {
	r.lastList.userID = userID
	r.lastList.filter = f
	if r.listFn != nil {
		return r.listFn(ctx, userID, f)
	}
	return []Transaction{}, 0, nil
}

func (r *stubRepo) Update(ctx context.Context, t Transaction) (Transaction, error) {
	r.lastUpdate = t
	if r.updateFn != nil {
		return r.updateFn(ctx, t)
	}
	return t, nil
}

func (r *stubRepo) Delete(ctx context.Context, userID, transactionID string) error {
	r.lastDelete.userID = userID
	r.lastDelete.transactionID = transactionID
	if r.deleteFn != nil {
		return r.deleteFn(ctx, userID, transactionID)
	}
	return nil
}

func strPtr(v string) *string {
	return &v
}

func TestServiceCreateSuccess(t *testing.T) {
	repo := &stubRepo{
		categoryAllowed: true,
	}
	svc := NewService(repo)

	got, err := svc.Create(context.Background(), CreateInput{
		UserID:      "user-1",
		OccurredAt:  "2026-03-18",
		Type:        "expense",
		AmountCents: 150000,
		CategoryID:  "cat-1",
		Merchant:    "Лента",
		Note:        strPtr(" продукты "),
	})
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if got.ID != "tx-1" {
		t.Fatalf("Create() returned unexpected id: %q", got.ID)
	}

	if repo.lastCreate.UserID != "user-1" {
		t.Fatalf("expected user_id user-1, got %q", repo.lastCreate.UserID)
	}
	if repo.lastCreate.Type != "expense" {
		t.Fatalf("expected type expense, got %q", repo.lastCreate.Type)
	}
	if repo.lastCreate.AmountCents != 150000 {
		t.Fatalf("expected amount 150000, got %d", repo.lastCreate.AmountCents)
	}
	if repo.lastCreate.Merchant != "Лента" {
		t.Fatalf("expected merchant Лента, got %q", repo.lastCreate.Merchant)
	}
	if repo.lastCreate.Note == nil || *repo.lastCreate.Note != "продукты" {
		t.Fatalf("expected trimmed note продукты, got %#v", repo.lastCreate.Note)
	}
}

func TestServiceCreateRejectsInvalidType(t *testing.T) {
	repo := &stubRepo{categoryAllowed: true}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateInput{
		UserID:      "user-1",
		OccurredAt:  "2026-03-18",
		Type:        "transfer",
		AmountCents: 1000,
		CategoryID:  "cat-1",
		Merchant:    "Bank",
	})
	if !errors.Is(err, ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}

func TestServiceCreateRejectsInvalidAmount(t *testing.T) {
	repo := &stubRepo{categoryAllowed: true}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateInput{
		UserID:      "user-1",
		OccurredAt:  "2026-03-18",
		Type:        "expense",
		AmountCents: 0,
		CategoryID:  "cat-1",
		Merchant:    "Лента",
	})
	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestServiceCreateRejectsCategoryNotAllowed(t *testing.T) {
	repo := &stubRepo{
		categoryAllowed: false,
	}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateInput{
		UserID:      "user-1",
		OccurredAt:  "2026-03-18",
		Type:        "expense",
		AmountCents: 1000,
		CategoryID:  "cat-1",
		Merchant:    "Лента",
	})
	if !errors.Is(err, ErrCategoryNotAllowed) {
		t.Fatalf("expected ErrCategoryNotAllowed, got %v", err)
	}
}

func TestServiceCreateUsesTodayWhenOccurredAtEmpty(t *testing.T) {
	repo := &stubRepo{
		categoryAllowed: true,
	}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateInput{
		UserID:      "user-1",
		OccurredAt:  "",
		Type:        "income",
		AmountCents: 500000,
		CategoryID:  "cat-1",
		Merchant:    "Работа",
	})
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if repo.lastCreate.OccurredAt == "" {
		t.Fatal("expected occurred_at to be auto-filled")
	}
}

func TestServiceListNormalizesLimitAndOffset(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	_, _, err := svc.List(context.Background(), ListInput{
		UserID: "user-1",
		From:   "2026-03-01",
		To:     "2026-03-31",
		Limit:  999,
		Offset: -10,
	})
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	if repo.lastList.userID != "user-1" {
		t.Fatalf("expected userID user-1, got %q", repo.lastList.userID)
	}
	if repo.lastList.filter.Limit != 200 {
		t.Fatalf("expected capped limit 200, got %d", repo.lastList.filter.Limit)
	}
	if repo.lastList.filter.Offset != 0 {
		t.Fatalf("expected normalized offset 0, got %d", repo.lastList.filter.Offset)
	}
}

func TestServiceListRejectsInvalidFromDate(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	_, _, err := svc.List(context.Background(), ListInput{
		UserID: "user-1",
		From:   "18-03-2026",
	})
	if !errors.Is(err, ErrInvalidFromDate) {
		t.Fatalf("expected ErrInvalidFromDate, got %v", err)
	}
}

func TestServiceUpdateSuccess(t *testing.T) {
	repo := &stubRepo{
		categoryAllowed: true,
	}
	svc := NewService(repo)

	_, err := svc.Update(context.Background(), UpdateInput{
		ID:          "tx-1",
		UserID:      "user-1",
		OccurredAt:  "2026-03-18",
		Type:        "expense",
		AmountCents: 1000,
		CategoryID:  "cat-1",
		Merchant:    "Лента",
		Note:        strPtr(" ужин "),
	})
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	if repo.lastUpdate.ID != "tx-1" {
		t.Fatalf("expected id tx-1, got %q", repo.lastUpdate.ID)
	}
	if repo.lastUpdate.Note == nil || *repo.lastUpdate.Note != "ужин" {
		t.Fatalf("expected trimmed note ужин, got %#v", repo.lastUpdate.Note)
	}
}

func TestServiceUpdateRejectsInvalidOccurredAt(t *testing.T) {
	repo := &stubRepo{
		categoryAllowed: true,
	}
	svc := NewService(repo)

	_, err := svc.Update(context.Background(), UpdateInput{
		ID:          "tx-1",
		UserID:      "user-1",
		OccurredAt:  "03/18/2026",
		Type:        "expense",
		AmountCents: 1000,
		CategoryID:  "cat-1",
		Merchant:    "Лента",
	})
	if !errors.Is(err, ErrInvalidOccurredAt) {
		t.Fatalf("expected ErrInvalidOccurredAt, got %v", err)
	}
}

func TestServiceDeleteRejectsEmptyTransactionID(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	err := svc.Delete(context.Background(), "user-1", "")
	if !errors.Is(err, ErrInvalidTransactionID) {
		t.Fatalf("expected ErrInvalidTransactionID, got %v", err)
	}
}

func TestServiceDeletePassesValuesToRepo(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	err := svc.Delete(context.Background(), "user-1", "tx-1")
	if err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}

	if repo.lastDelete.userID != "user-1" {
		t.Fatalf("expected userID user-1, got %q", repo.lastDelete.userID)
	}
	if repo.lastDelete.transactionID != "tx-1" {
		t.Fatalf("expected transactionID tx-1, got %q", repo.lastDelete.transactionID)
	}
}
