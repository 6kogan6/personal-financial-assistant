package transactions

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCategoryNotAllowed = errors.New("category not allowed")
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// CategoryAllowed: категория допустима, если is_system=true OR (user_id = текущий пользователь)
func (r *Repo) CategoryAllowed(ctx context.Context, userID, categoryID string) (bool, error) {
	const q = `
SELECT EXISTS(
  SELECT 1
  FROM categories
  WHERE id = $1 AND (is_system = true OR user_id = $2)
);
`
	var ok bool
	if err := r.pool.QueryRow(ctx, q, categoryID, userID).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *Repo) Create(ctx context.Context, t Transaction) (Transaction, error) {
	const q = `
INSERT INTO transactions (user_id, occurred_at, type, amount_cents, category_id, merchant, note, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'manual')
RETURNING id, user_id, occurred_at::text, type, amount_cents, category_id, merchant, note, source, created_at;
`
	var out Transaction
	var createdAt time.Time
	var note *string

	err := r.pool.QueryRow(ctx, q,
		t.UserID, t.OccurredAt, t.Type, t.AmountCents, t.CategoryID, t.Merchant, t.Note,
	).Scan(&out.ID, &out.UserID, &out.OccurredAt, &out.Type, &out.AmountCents, &out.CategoryID,
		&out.Merchant, &note, &out.Source, &createdAt)
	if err != nil {
		return Transaction{}, err
	}

	out.Note = note
	out.CreatedAt = createdAt
	return out, nil
}

func (r *Repo) List(ctx context.Context, userID string, f ListFilter) ([]Transaction, int, error) {
	// Фильтры делаем через условия "($X = '' OR ...)" чтобы не городить билдера.
	// Для MVP норм.
	const q = `
SELECT
  id, user_id, occurred_at::text, type, amount_cents, category_id, merchant, note, source, created_at,
  COUNT(*) OVER() AS total
FROM transactions
WHERE user_id = $1
  AND ($2 = '' OR occurred_at >= $2::date)
  AND ($3 = '' OR occurred_at <= $3::date)
  AND ($4 = '' OR category_id = $4::uuid)
  AND (
    $5 = '' OR
    merchant ILIKE '%' || $5 || '%' OR
    COALESCE(note,'') ILIKE '%' || $5 || '%'
  )
ORDER BY occurred_at DESC, created_at DESC
LIMIT 200;
`
	rows, err := r.pool.Query(ctx, q, userID, f.From, f.To, f.CategoryID, f.Q)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var res []Transaction
	total := 0

	for rows.Next() {
		var t Transaction
		var note *string
		var createdAt time.Time
		var rowTotal int
		if err := rows.Scan(&t.ID, &t.UserID, &t.OccurredAt, &t.Type, &t.AmountCents, &t.CategoryID,
			&t.Merchant, &note, &t.Source, &createdAt, &rowTotal); err != nil {
			return nil, 0, err
		}
		t.Note = note
		t.CreatedAt = createdAt
		total = rowTotal
		res = append(res, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return res, total, nil
}
