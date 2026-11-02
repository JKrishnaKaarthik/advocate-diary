package app

import (
	"context"
	"errors"
	"mime"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/validate"
)

type Document struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	ByteSize     int64     `json:"byte_size"`
	Label        string    `json:"label"`
	CreatedAt    time.Time `json:"created_at"`
}

type Download struct {
	URL       string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}

func (s *Store) Upload(ctx context.Context, owner, caseID, filename, contentType, label string, body []byte) (Document, error) {
	if !validate.Label(label) {
		return Document{}, Bad("Label is not valid")
	}
	if err := validate.CheckFile(filename, contentType, int64(len(body))); err != nil {
		return Document{}, Bad(err.Error())
	}
	var owns bool
	err := s.pool.QueryRow(ctx, `SELECT true FROM cases WHERE id = $1 AND owner_id = $2`, caseID, owner).Scan(&owns)
	if errors.Is(err, pgx.ErrNoRows) {
		return Document{}, NotFound()
	}
	if err != nil {
		return Document{}, err
	}
	if s.objects == nil {
		return Document{}, errors.New("document storage is not configured")
	}
	id := uuid.NewString()
	name := safeName(filename)
	key := path.Join("lawyers", owner, "cases", caseID, id, name)
	media := storedType(filename, contentType)
	if err := s.objects.Put(ctx, key, media, body); err != nil {
		return Document{}, err
	}
	var doc Document
	err = s.pool.QueryRow(ctx, `
		INSERT INTO documents (id, case_id, original_name, content_type, byte_size, storage_key, label)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, original_name, content_type, byte_size, label, created_at`,
		id, caseID, filepath.Base(strings.ReplaceAll(filename, "\\", "/")), media, len(body), key, label).
		Scan(&doc.ID, &doc.OriginalName, &doc.ContentType, &doc.ByteSize, &doc.Label, &doc.CreatedAt)
	if err != nil {
		_ = s.objects.Delete(ctx, key)
		return Document{}, err
	}
	return doc, nil
}

func (s *Store) Download(ctx context.Context, owner, id string) (Download, error) {
	key, err := s.objectKey(ctx, owner, id)
	if err != nil {
		return Download{}, err
	}
	url, err := s.objects.PresignGet(ctx, key)
	if err != nil {
		return Download{}, err
	}
	return Download{URL: url, ExpiresIn: 300}, nil
}

func (s *Store) DeleteDocument(ctx context.Context, owner, id string) error {
	key, err := s.objectKey(ctx, owner, id)
	if err != nil {
		return err
	}
	if err := s.objects.Delete(ctx, key); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM documents d USING cases c
		WHERE d.id = $1 AND d.case_id = c.id AND c.owner_id = $2`, id, owner)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return NotFound()
	}
	return nil
}

func (s *Store) objectKey(ctx context.Context, owner, id string) (string, error) {
	var key string
	err := s.pool.QueryRow(ctx, `
		SELECT d.storage_key
		FROM documents d
		JOIN cases c ON c.id = d.case_id
		WHERE d.id = $1 AND c.owner_id = $2`, id, owner).Scan(&key)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", NotFound()
	}
	return key, err
}

func (s *Store) listDocuments(ctx context.Context, owner, caseID string) ([]Document, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT d.id::text, d.original_name, d.content_type, d.byte_size, d.label, d.created_at
		FROM documents d
		JOIN cases c ON c.id = d.case_id
		WHERE d.case_id = $1 AND c.owner_id = $2
		ORDER BY d.created_at DESC`, caseID, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Document
	for rows.Next() {
		var d Document
		if err := rows.Scan(&d.ID, &d.OriginalName, &d.ContentType, &d.ByteSize, &d.Label, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	if out == nil {
		out = []Document{}
	}
	return out, rows.Err()
}

func safeName(name string) string {
	base := filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	var b strings.Builder
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '.', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	s := b.String()
	if s == "" || s == "." || s == ".." {
		return "file"
	}
	return s
}

func storedType(filename, contentType string) string {
	media, _, err := mime.ParseMediaType(contentType)
	if err != nil || media == "" || strings.EqualFold(media, "application/octet-stream") {
		return "application/pdf"
	}
	return media
}
