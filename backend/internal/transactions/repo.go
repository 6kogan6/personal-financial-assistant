package transactions

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

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
	var note *string
	var createdAt time.Time

	err := r.pool.QueryRow(
		ctx,
		q,
		t.UserID,
		t.OccurredAt,
		t.Type,
		t.AmountCents,
		t.CategoryID,
		t.Merchant,
		t.Note,
	).Scan(
		&out.ID,
		&out.UserID,
		&out.OccurredAt,
		&out.Type,
		&out.AmountCents,
		&out.CategoryID,
		&out.Merchant,
		&note,
		&out.Source,
		&createdAt,
	)
	if err != nil {
		return Transaction{}, err
	}

	out.Note = note
	out.CreatedAt = createdAt
	return out, nil
}

func (r *Repo) List(ctx context.Context, userID string, f ListFilter) ([]Transaction, int, error) {
	const q = `
SELECT
  id,
  user_id,
  occurred_at::text,
  type,
  amount_cents,
  category_id,
  merchant,
  note,
  source,
  created_at,
  COUNT(*) OVER() AS total
FROM transactions
WHERE user_id = $1
  AND ($2 = '' OR occurred_at >= $2::date)
  AND ($3 = '' OR occurred_at <= $3::date)
  AND ($4 = '' OR category_id = $4::uuid)
  AND (
    $5 = '' OR
    merchant ILIKE '%' || $5 || '%' OR
    COALESCE(note, '') ILIKE '%' || $5 || '%'
  )
ORDER BY occurred_at DESC, created_at DESC
LIMIT $6 OFFSET $7;
`

	rows, err := r.pool.Query(ctx, q, userID, f.From, f.To, f.CategoryID, f.Q, f.Limit, f.Offset)
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

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.OccurredAt,
			&t.Type,
			&t.AmountCents,
			&t.CategoryID,
			&t.Merchant,
			&note,
			&t.Source,
			&createdAt,
			&rowTotal,
		); err != nil {
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

func (r *Repo) Update(ctx context.Context, t Transaction) (Transaction, error) {
	const q = `
UPDATE transactions
SET
  occurred_at = $3,
  type = $4,
  amount_cents = $5,
  category_id = $6,
  merchant = $7,
  note = $8
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, occurred_at::text, type, amount_cents, category_id, merchant, note, source, created_at;
`

	var out Transaction
	var note *string
	var createdAt time.Time

	err := r.pool.QueryRow(
		ctx,
		q,
		t.ID,
		t.UserID,
		t.OccurredAt,
		t.Type,
		t.AmountCents,
		t.CategoryID,
		t.Merchant,
		t.Note,
	).Scan(
		&out.ID,
		&out.UserID,
		&out.OccurredAt,
		&out.Type,
		&out.AmountCents,
		&out.CategoryID,
		&out.Merchant,
		&note,
		&out.Source,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Transaction{}, ErrTransactionNotFound
		}
		return Transaction{}, err
	}

	out.Note = note
	out.CreatedAt = createdAt
	return out, nil
}

func (r *Repo) Delete(ctx context.Context, userID, transactionID string) error {
	const q = `
DELETE FROM transactions
WHERE id = $1 AND user_id = $2;
`

	tag, err := r.pool.Exec(ctx, q, transactionID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}
	return nil
}
