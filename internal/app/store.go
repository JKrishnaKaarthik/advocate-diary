package app

import (
	"context"
	_ "embed"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

type ObjectStore interface {
	Put(ctx context.Context, key, contentType string, body []byte) error
	Delete(ctx context.Context, key string) error
	PresignGet(ctx context.Context, key string) (string, error)
}

type Store struct {
	pool    *pgxpool.Pool
	objects ObjectStore
	ticking atomic.Bool
}

func NewStore(pool *pgxpool.Pool, objects ObjectStore) *Store {
	return &Store{pool: pool, objects: objects}
}

func (s *Store) Migrate(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT now()
		)`); err != nil {
		return err
	}
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM schema_migrations WHERE version = '001'`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	for _, stmt := range strings.Split(schemaSQL, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := s.pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ('001')`)
	return err
}

func dateString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func like(q string) string {
	q = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(strings.TrimSpace(q))
	return "%" + q + "%"
}
