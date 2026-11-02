package app

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash []byte

func init() {
	var err error
	dummyHash, err = bcrypt.GenerateFromPassword([]byte("missing-account"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
}

type Lawyer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Store) SignUp(ctx context.Context, name, email, password string) (Lawyer, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if name == "" || !strings.Contains(email, "@") || strings.Contains(email, " ") || len(password) < 8 {
		return Lawyer{}, Bad("Name, a valid email, and a password of at least 8 characters are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Lawyer{}, err
	}
	id := uuid.NewString()
	_, err = s.pool.Exec(ctx, `INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`, id, name, email, string(hash))
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return Lawyer{}, Conflict("An account with that email already exists")
	}
	if err != nil {
		return Lawyer{}, err
	}
	return Lawyer{ID: id, Name: name, Email: email}, nil
}

func (s *Store) SignIn(ctx context.Context, email, password string) (Lawyer, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var lawyer Lawyer
	var hash string
	err := s.pool.QueryRow(ctx, `SELECT id::text, name, email, password_hash FROM users WHERE email = $1`, email).
		Scan(&lawyer.ID, &lawyer.Name, &lawyer.Email, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		return Lawyer{}, Bad("Email or password is wrong")
	}
	if err != nil {
		return Lawyer{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return Lawyer{}, Bad("Email or password is wrong")
	}
	return lawyer, nil
}

func (s *Store) Lawyer(ctx context.Context, id string) (Lawyer, error) {
	var lawyer Lawyer
	err := s.pool.QueryRow(ctx, `SELECT id::text, name, email FROM users WHERE id = $1`, id).
		Scan(&lawyer.ID, &lawyer.Name, &lawyer.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return Lawyer{}, NotFound()
	}
	return lawyer, err
}
