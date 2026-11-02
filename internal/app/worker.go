package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/mail"
)

type dueReminder struct {
	ID            string
	Attempts      int
	UserID        string
	Email         string
	CaseID        string
	Title         string
	CaseNumber    string
	ClientName    string
	ScheduledAt   time.Time
	CourtName     string
	Purpose       string
	HearingStatus string
}

func (s *Store) Tick(ctx context.Context, sender mail.Sender) error {
	if !s.ticking.CompareAndSwap(false, true) {
		return nil
	}
	defer s.ticking.Store(false)

	if _, err := s.pool.Exec(ctx, `
		UPDATE reminders
		SET status = 'scheduled', claimed_at = NULL
		WHERE status = 'sending' AND claimed_at < now() - interval '10 minutes'`); err != nil {
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `
		SELECT id::text FROM reminders
		WHERE status = 'scheduled' AND remind_at <= now()
		ORDER BY remind_at
		LIMIT 50
		FOR UPDATE SKIP LOCKED`)
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := tx.Exec(ctx, `
			UPDATE reminders
			SET status = 'sending', attempts = attempts + 1, claimed_at = now()
			WHERE id = $1`, id); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	var sendErr error
	for _, id := range ids {
		if err := s.deliver(ctx, sender, id); err != nil && sendErr == nil {
			sendErr = err
		}
	}
	return sendErr
}

func (s *Store) deliver(ctx context.Context, sender mail.Sender, id string) error {
	var item dueReminder
	err := s.pool.QueryRow(ctx, `
		SELECT r.id::text, r.attempts, u.id::text, u.email, c.id::text, c.title, c.case_number, cl.full_name,
		       h.scheduled_at, h.court_name, h.purpose, h.status
		FROM reminders r
		JOIN hearings h ON h.id = r.hearing_id
		JOIN cases c ON c.id = h.case_id
		JOIN clients cl ON cl.id = c.client_id
		JOIN users u ON u.id = c.owner_id
		WHERE r.id = $1`, id).Scan(
		&item.ID, &item.Attempts, &item.UserID, &item.Email, &item.CaseID, &item.Title, &item.CaseNumber, &item.ClientName,
		&item.ScheduledAt, &item.CourtName, &item.Purpose, &item.HearingStatus)
	if err != nil {
		return err
	}
	if item.HearingStatus != "scheduled" {
		_, err = s.pool.Exec(ctx, `UPDATE reminders SET status = 'cancelled', claimed_at = NULL WHERE id = $1 AND status = 'sending'`, id)
		return err
	}
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*3600+30*60)
	}
	number := item.CaseNumber
	if number == "" {
		number = "Case number not allotted"
	}
	body := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		item.ClientName, item.Title, number, item.CourtName, item.Purpose,
		item.ScheduledAt.In(loc).Format("2 Jan 2006, 3:04 PM"))
	subject := "Court reminder: " + item.Title
	if err := sender.Send(ctx, item.Email, subject, body); err != nil {
		status := "scheduled"
		if item.Attempts >= 3 {
			status = "failed"
		}
		_, updateErr := s.pool.Exec(ctx, `
			UPDATE reminders
			SET status = $2, last_error = $3, claimed_at = NULL
			WHERE id = $1`, id, status, trimErr(err))
		if updateErr != nil {
			return updateErr
		}
		return err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		INSERT INTO notifications (id, user_id, reminder_id, title, body)
		VALUES ($1, $2, $3, $4, $5)`, newID(), item.UserID, id, "Hearing · "+item.Title, body); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE reminders SET status = 'sent', sent_at = now(), last_error = '', claimed_at = NULL
		WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func trimErr(err error) string {
	msg := err.Error()
	msg = strings.ReplaceAll(msg, "\n", " ")
	if len(msg) > 500 {
		return msg[:500]
	}
	return msg
}
