package users

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, email, passwordHash string) (User, error) {
	const q = `
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, created_at;
`
	var u User
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, q, email, passwordHash).Scan(&u.ID, &u.Email, &createdAt)
	if err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrEmailExists
		}
		return User{}, err
	}
	u.CreatedAt = createdAt
	return u, nil
}

func (r *Repo) GetAuthDataByEmail(ctx context.Context, email string) (userID string, passwordHash string, err error) {
	const q = `SELECT id, password_hash FROM users WHERE email = $1;`
	err = r.pool.QueryRow(ctx, q, email).Scan(&userID, &passwordHash)
	return userID, passwordHash, err
}

func (r *Repo) GetByID(ctx context.Context, userID string) (User, error) {
	const q = `
SELECT id, email, created_at
FROM users
WHERE id = $1;
`

	var u User
	var createdAt time.Time

	err := r.pool.QueryRow(ctx, q, userID).Scan(&u.ID, &u.Email, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	u.CreatedAt = createdAt
	return u, nil
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
