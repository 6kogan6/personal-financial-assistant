package reports

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Summary(ctx context.Context, userID string, start time.Time, end time.Time, month string) (Summary, error) {
	// 1) totals
	const totalsQ = `
SELECT
  COALESCE(SUM(CASE WHEN type='income' THEN amount_cents ELSE 0 END), 0) AS income_cents,
  COALESCE(SUM(CASE WHEN type='expense' THEN amount_cents ELSE 0 END), 0) AS expense_cents
FROM transactions
WHERE user_id = $1
  AND occurred_at >= $2::date
  AND occurred_at <  $3::date;
`
	var income, expense int64
	if err := r.pool.QueryRow(ctx, totalsQ, userID, start, end).Scan(&income, &expense); err != nil {
		return Summary{}, err
	}

	// 2) by_category
	const byCatQ = `
SELECT
  t.category_id,
  c.name,
  c.type,
  COALESCE(SUM(t.amount_cents), 0) AS amount_cents
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
  AND t.occurred_at >= $2::date
  AND t.occurred_at <  $3::date
GROUP BY t.category_id, c.name, c.type
ORDER BY amount_cents DESC, c.name ASC;
`
	rows, err := r.pool.Query(ctx, byCatQ, userID, start, end)
	if err != nil {
		return Summary{}, err
	}
	defer rows.Close()

	byCategory := make([]ByCategoryItem, 0)
	for rows.Next() {
		var it ByCategoryItem
		if err := rows.Scan(&it.CategoryID, &it.CategoryName, &it.Type, &it.AmountCents); err != nil {
			return Summary{}, err
		}
		byCategory = append(byCategory, it)
	}
	if err := rows.Err(); err != nil {
		return Summary{}, err
	}

	// 3) daily
	const dailyQ = `
SELECT
  occurred_at::text AS date,
  COALESCE(SUM(CASE WHEN type='income' THEN amount_cents ELSE 0 END), 0) AS income_cents,
  COALESCE(SUM(CASE WHEN type='expense' THEN amount_cents ELSE 0 END), 0) AS expense_cents
FROM transactions
WHERE user_id = $1
  AND occurred_at >= $2::date
  AND occurred_at <  $3::date
GROUP BY occurred_at
ORDER BY occurred_at ASC;
`
	rows2, err := r.pool.Query(ctx, dailyQ, userID, start, end)
	if err != nil {
		return Summary{}, err
	}
	defer rows2.Close()

	daily := make([]DailyItem, 0)
	for rows2.Next() {
		var d DailyItem
		if err := rows2.Scan(&d.Date, &d.IncomeCents, &d.ExpenseCents); err != nil {
			return Summary{}, err
		}
		daily = append(daily, d)
	}
	if err := rows2.Err(); err != nil {
		return Summary{}, err
	}

	return Summary{
		Month: month,
		Totals: Totals{
			IncomeCents:  income,
			ExpenseCents: expense,
			BalanceCents: income - expense,
		},
		ByCategory: byCategory,
		Daily:      daily,
	}, nil
}
