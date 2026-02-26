package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is empty")
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 5
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, err
	}

	db := &DB{Pool: pool}

	if err := db.ensureSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return db, nil
}

func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

func (d *DB) ensureSchema(ctx context.Context) error {
	const q = `
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NULL REFERENCES users(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	is_system BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT categories_type_check CHECK (type IN ('expense','income'))
);

CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  occurred_at DATE NOT NULL,
  type TEXT NOT NULL,
  amount_cents BIGINT NOT NULL CHECK (amount_cents >= 0),
  category_id UUID NOT NULL REFERENCES categories(id),
  merchant TEXT NOT NULL,
  note TEXT NULL,
  source TEXT NOT NULL DEFAULT 'manual',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT transactions_type_check CHECK (type IN ('expense','income')),
  CONSTRAINT transactions_source_check CHECK (source IN ('manual','import'))
);

CREATE INDEX IF NOT EXISTS ix_transactions_user_date
ON transactions (user_id, occurred_at DESC);

CREATE INDEX IF NOT EXISTS ix_transactions_user_category
ON transactions (user_id, category_id);

-- Уникальность для системных категорий (по name+type)
CREATE UNIQUE INDEX IF NOT EXISTS ux_categories_system_name_type
ON categories (name, type)
WHERE is_system = true;

-- Уникальность для пользовательских категорий (в рамках user_id)
CREATE UNIQUE INDEX IF NOT EXISTS ux_categories_user_name_type
ON categories (user_id, name, type)
WHERE is_system = false;

-- Сиды системных категорий (можно расширять)
INSERT INTO categories (user_id, name, type, is_system)
VALUES
  (NULL, 'Еда', 'expense', true),
  (NULL, 'Транспорт', 'expense', true),
  (NULL, 'Дом', 'expense', true),
  (NULL, 'Развлечения', 'expense', true),
  (NULL, 'Здоровье', 'expense', true),
  (NULL, 'Зарплата', 'income', true),
  (NULL, 'Подарки', 'income', true)
ON CONFLICT DO NOTHING;
`
	_, err := d.Pool.Exec(ctx, q)
	return err
}
