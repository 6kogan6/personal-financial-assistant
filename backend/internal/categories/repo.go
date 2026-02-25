package categories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAlreadyExists = errors.New("category already exists")
	ErrInvalidType   = errors.New("invalid category type")
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// ListForUser возвращает системные + пользовательские категории.
func (r *Repo) ListForUser(ctx context.Context, userID string) ([]Category, error) {
	const q = `
SELECT id, user_id, name, type, is_system, created_at
FROM categories
WHERE is_system = true OR user_id = $1
ORDER BY is_system DESC, type ASC, name ASC;
`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Category
	for rows.Next() {
		var c Category
		var uid *string
		var createdAt time.Time
		if err := rows.Scan(&c.ID, &uid, &c.Name, &c.Type, &c.IsSystem, &createdAt); err != nil {
			return nil, err
		}
		c.UserID = uid
		c.CreatedAt = createdAt
		res = append(res, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUserCategory создаёт пользовательскую категорию (is_system=false).
func (r *Repo) CreateUserCategory(ctx context.Context, userID, name, typ string) (Category, error) {
	const q = `
INSERT INTO categories (user_id, name, type, is_system)
VALUES ($1, $2, $3, false)
RETURNING id, user_id, name, type, is_system, created_at;
`
	var c Category
	var uid *string
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, q, userID, name, typ).Scan(&c.ID, &uid, &c.Name, &c.Type, &c.IsSystem, &createdAt)
	if err != nil {
		if isUniqueViolation(err) {
			return Category{}, ErrAlreadyExists
		}
		return Category{}, err
	}
	c.UserID = uid
	c.CreatedAt = createdAt
	return c, nil
}

func isUniqueViolation(err error) bool {
	s := err.Error()
	return contains(s, "duplicate key") || contains(s, "unique constraint") || contains(s, "23505")
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
