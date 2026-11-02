package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DiaryEntry struct {
	ID          string    `json:"id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CourtName   string    `json:"court_name"`
	Purpose     string    `json:"purpose"`
	CaseID      string    `json:"case_id"`
	Title       string    `json:"title"`
	CaseNumber  string    `json:"case_number"`
	CaseType    string    `json:"case_type"`
	ClientName  string    `json:"client_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
}

type Notification struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
	CaseID    string     `json:"case_id"`
}

func (s *Store) Diary(ctx context.Context, owner string, from, to time.Time) ([]DiaryEntry, error) {
	if !to.After(from) {
		return nil, Bad("The diary range is not valid")
	}
	rows, err := s.pool.Query(ctx, `
		SELECT h.id::text, h.scheduled_at, h.court_name, h.purpose,
		       c.id::text, c.title, c.case_number, c.case_type,
		       cl.full_name, cl.phone, cl.email, cl.address
		FROM hearings h
		JOIN cases c ON c.id = h.case_id
		JOIN clients cl ON cl.id = c.client_id
		WHERE c.owner_id = $1
		  AND h.status = 'scheduled'
		  AND h.scheduled_at >= $2
		  AND h.scheduled_at < $3
		ORDER BY h.scheduled_at`, owner, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DiaryEntry
	for rows.Next() {
		var e DiaryEntry
		if err := rows.Scan(&e.ID, &e.ScheduledAt, &e.CourtName, &e.Purpose, &e.CaseID, &e.Title, &e.CaseNumber, &e.CaseType, &e.ClientName, &e.Phone, &e.Email, &e.Address); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	if out == nil {
		out = []DiaryEntry{}
	}
	return out, rows.Err()
}

func (s *Store) Notifications(ctx context.Context, owner string) ([]Notification, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT n.id::text, n.title, n.body, n.read_at, n.created_at, c.id::text
		FROM notifications n
		JOIN reminders r ON r.id = n.reminder_id
		JOIN hearings h ON h.id = r.hearing_id
		JOIN cases c ON c.id = h.case_id
		WHERE n.user_id = $1
		ORDER BY (n.read_at IS NULL) DESC, n.created_at DESC`, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.ReadAt, &n.CreatedAt, &n.CaseID); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	if out == nil {
		out = []Notification{}
	}
	return out, rows.Err()
}

func (s *Store) MarkRead(ctx context.Context, owner, id string) error {
	var got string
	err := s.pool.QueryRow(ctx, `
		UPDATE notifications
		SET read_at = COALESCE(read_at, now())
		WHERE id = $1 AND user_id = $2
		RETURNING id::text`, id, owner).Scan(&got)
	if errors.Is(err, pgx.ErrNoRows) {
		return NotFound()
	}
	return err
}

func newID() string { return uuid.NewString() }
