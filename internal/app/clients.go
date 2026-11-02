package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/validate"
)

type Client struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClientDetail struct {
	Client
	Cases []CaseSummary `json:"cases"`
}

type ClientInput struct {
	FullName string
	Phone    string
	Email    string
	Address  string
	Notes    string
}

type ClientPatch struct {
	FullName *string
	Phone    *string
	Email    *string
	Address  *string
	Notes    *string
}

func (s *Store) CreateClient(ctx context.Context, owner string, in ClientInput) (Client, error) {
	in.FullName = strings.TrimSpace(in.FullName)
	in.Email = strings.TrimSpace(in.Email)
	if in.FullName == "" {
		return Client{}, Bad("Name is required")
	}
	if !validate.Email(in.Email) {
		return Client{}, Bad("Enter a valid email")
	}
	c := Client{
		ID: uuid.NewString(), FullName: in.FullName,
		Phone: strings.TrimSpace(in.Phone), Email: in.Email,
		Address: strings.TrimSpace(in.Address), Notes: strings.TrimSpace(in.Notes),
	}
	err := s.pool.QueryRow(ctx, `
		INSERT INTO clients (id, owner_id, full_name, phone, email, address, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`,
		c.ID, owner, c.FullName, c.Phone, c.Email, c.Address, c.Notes).Scan(&c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (s *Store) ListClients(ctx context.Context, owner, q string) ([]Client, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id::text, full_name, phone, email, address, notes, created_at, updated_at
		FROM clients
		WHERE owner_id = $1 AND ($2 = '' OR full_name ILIKE $3 ESCAPE '\')
		ORDER BY full_name`, owner, strings.TrimSpace(q), like(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.FullName, &c.Phone, &c.Email, &c.Address, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []Client{}
	}
	return out, rows.Err()
}

func (s *Store) GetClient(ctx context.Context, owner, id string) (ClientDetail, error) {
	var d ClientDetail
	err := s.pool.QueryRow(ctx, `
		SELECT id::text, full_name, phone, email, address, notes, created_at, updated_at
		FROM clients WHERE id = $1 AND owner_id = $2`, id, owner).
		Scan(&d.ID, &d.FullName, &d.Phone, &d.Email, &d.Address, &d.Notes, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ClientDetail{}, NotFound()
	}
	if err != nil {
		return ClientDetail{}, err
	}
	d.Cases, err = s.casesForClient(ctx, owner, id)
	return d, err
}

func (s *Store) UpdateClient(ctx context.Context, owner, id string, patch ClientPatch) (Client, error) {
	current, err := s.GetClient(ctx, owner, id)
	if err != nil {
		return Client{}, err
	}
	c := current.Client
	if patch.FullName != nil {
		c.FullName = strings.TrimSpace(*patch.FullName)
		if c.FullName == "" {
			return Client{}, Bad("Client name is required")
		}
	}
	if patch.Phone != nil {
		c.Phone = strings.TrimSpace(*patch.Phone)
	}
	if patch.Email != nil {
		c.Email = strings.TrimSpace(*patch.Email)
	}
	if patch.Address != nil {
		c.Address = strings.TrimSpace(*patch.Address)
	}
	if patch.Notes != nil {
		c.Notes = strings.TrimSpace(*patch.Notes)
	}
	if c.FullName == "" {
		return Client{}, Bad("Name is required")
	}
	if !validate.Email(c.Email) {
		return Client{}, Bad("Enter a valid email")
	}
	err = s.pool.QueryRow(ctx, `
		UPDATE clients
		SET full_name = $3, phone = $4, email = $5, address = $6, notes = $7, updated_at = now()
		WHERE id = $1 AND owner_id = $2
		RETURNING created_at, updated_at`,
		id, owner, c.FullName, c.Phone, c.Email, c.Address, c.Notes).Scan(&c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Client{}, NotFound()
	}
	return c, err
}
