package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/reminders"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/validate"
)

type Hearing struct {
	ID                string    `json:"id"`
	ScheduledAt       time.Time `json:"scheduled_at"`
	CourtName         string    `json:"court_name"`
	Purpose           string    `json:"purpose"`
	Status            string    `json:"status"`
	Outcome           string    `json:"outcome"`
	PreviousHearingID *string   `json:"previous_hearing_id"`
	CreatedAt         time.Time `json:"created_at"`
}

type HearingResult struct {
	Hearing
	RemindersPlanned int `json:"reminders_planned"`
}

func (s *Store) AddHearing(ctx context.Context, owner, caseID, courtName, purpose string, at time.Time) (HearingResult, error) {
	if !validate.Purpose(purpose) {
		return HearingResult{}, Bad("Purpose is not valid")
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return HearingResult{}, err
	}
	defer tx.Rollback(ctx)

	lockedCourt, status, err := lockCase(ctx, tx, owner, caseID)
	if err != nil {
		return HearingResult{}, err
	}
	if status == "disposed" {
		return HearingResult{}, Conflict("A disposed case cannot take a new hearing")
	}
	courtName = strings.TrimSpace(courtName)
	if courtName == "" {
		courtName = lockedCourt
	}
	if courtName == "" {
		return HearingResult{}, Bad("Court is required")
	}
	h := Hearing{
		ID: uuid.NewString(), ScheduledAt: at.UTC(), CourtName: courtName, Purpose: purpose, Status: "scheduled",
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO hearings (id, case_id, scheduled_at, court_name, purpose, status)
		VALUES ($1, $2, $3, $4, $5, 'scheduled')
		RETURNING created_at`, h.ID, caseID, h.ScheduledAt, h.CourtName, h.Purpose).Scan(&h.CreatedAt)
	if err != nil {
		return HearingResult{}, err
	}
	n, err := planReminders(ctx, tx, h.ID, h.ScheduledAt, "scheduled")
	if err != nil {
		return HearingResult{}, err
	}
	if err := recomputeNext(ctx, tx, caseID); err != nil {
		return HearingResult{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return HearingResult{}, err
	}
	return HearingResult{Hearing: h, RemindersPlanned: n}, nil
}

func (s *Store) CompleteHearing(ctx context.Context, owner, id, outcome string) (Hearing, error) {
	outcome = strings.TrimSpace(outcome)
	if outcome == "" {
		return Hearing{}, Bad("Outcome is required")
	}
	return s.closeHearing(ctx, owner, id, "completed", outcome)
}

func (s *Store) CancelHearing(ctx context.Context, owner, id, outcome string) (Hearing, error) {
	return s.closeHearing(ctx, owner, id, "cancelled", strings.TrimSpace(outcome))
}

func (s *Store) AdjournHearing(ctx context.Context, owner, id, outcome, purpose string, nextAt time.Time) (HearingResult, error) {
	outcome = strings.TrimSpace(outcome)
	if outcome == "" {
		return HearingResult{}, Bad("Outcome is required")
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return HearingResult{}, err
	}
	defer tx.Rollback(ctx)

	old, caseID, err := lockHearing(ctx, tx, owner, id)
	if err != nil {
		return HearingResult{}, err
	}
	if old.Status != "scheduled" {
		return HearingResult{}, Conflict("Only a scheduled hearing can be adjourned")
	}
	if !nextAt.After(old.ScheduledAt) {
		return HearingResult{}, Bad("The next date must be later than this hearing")
	}
	if purpose == "" {
		purpose = old.Purpose
	}
	if !validate.Purpose(purpose) {
		return HearingResult{}, Bad("Purpose is not valid")
	}
	if _, err := tx.Exec(ctx, `
		UPDATE hearings SET status = 'adjourned', outcome = $2, updated_at = now() WHERE id = $1`, id, outcome); err != nil {
		return HearingResult{}, err
	}
	if err := cancelScheduledReminders(ctx, tx, id); err != nil {
		return HearingResult{}, err
	}
	next := Hearing{
		ID: uuid.NewString(), ScheduledAt: nextAt.UTC(), CourtName: old.CourtName, Purpose: purpose, Status: "scheduled",
		PreviousHearingID: &id,
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO hearings (id, case_id, scheduled_at, court_name, purpose, status, previous_hearing_id)
		VALUES ($1, $2, $3, $4, $5, 'scheduled', $6)
		RETURNING created_at`, next.ID, caseID, next.ScheduledAt, next.CourtName, next.Purpose, id).Scan(&next.CreatedAt)
	if err != nil {
		return HearingResult{}, err
	}
	n, err := planReminders(ctx, tx, next.ID, next.ScheduledAt, "scheduled")
	if err != nil {
		return HearingResult{}, err
	}
	if err := recomputeNext(ctx, tx, caseID); err != nil {
		return HearingResult{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return HearingResult{}, err
	}
	return HearingResult{Hearing: next, RemindersPlanned: n}, nil
}

func (s *Store) closeHearing(ctx context.Context, owner, id, status, outcome string) (Hearing, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Hearing{}, err
	}
	defer tx.Rollback(ctx)
	return finishInTx(ctx, tx, owner, id, status, outcome)
}

func finishInTx(ctx context.Context, tx pgx.Tx, owner, id, status, outcome string) (Hearing, error) {
	h, caseID, err := lockHearing(ctx, tx, owner, id)
	if err != nil {
		return Hearing{}, err
	}
	if h.Status != "scheduled" {
		return Hearing{}, Conflict("Only a scheduled hearing can be updated")
	}
	h.Status = status
	h.Outcome = outcome
	if _, err := tx.Exec(ctx, `
		UPDATE hearings SET status = $2, outcome = $3, updated_at = now() WHERE id = $1`, id, status, outcome); err != nil {
		return Hearing{}, err
	}
	if err := cancelScheduledReminders(ctx, tx, id); err != nil {
		return Hearing{}, err
	}
	if err := recomputeNext(ctx, tx, caseID); err != nil {
		return Hearing{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Hearing{}, err
	}
	return h, nil
}

func (s *Store) listHearings(ctx context.Context, owner, caseID string) ([]Hearing, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT h.id::text, h.scheduled_at, h.court_name, h.purpose, h.status, h.outcome,
		       h.previous_hearing_id::text, h.created_at
		FROM hearings h
		JOIN cases c ON c.id = h.case_id
		WHERE h.case_id = $1 AND c.owner_id = $2
		ORDER BY h.scheduled_at DESC`, caseID, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Hearing
	for rows.Next() {
		h, err := scanHearing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	if out == nil {
		out = []Hearing{}
	}
	return out, rows.Err()
}

func scanHearing(row scanner) (Hearing, error) {
	var h Hearing
	var prev *string
	err := row.Scan(&h.ID, &h.ScheduledAt, &h.CourtName, &h.Purpose, &h.Status, &h.Outcome, &prev, &h.CreatedAt)
	if prev != nil && *prev != "" {
		h.PreviousHearingID = prev
	}
	return h, err
}

func lockCase(ctx context.Context, tx pgx.Tx, owner, caseID string) (court, status string, err error) {
	err = tx.QueryRow(ctx, `
		SELECT court_name, status FROM cases WHERE id = $1 AND owner_id = $2 FOR UPDATE`, caseID, owner).
		Scan(&court, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", NotFound()
	}
	return court, status, err
}

func lockHearing(ctx context.Context, tx pgx.Tx, owner, id string) (Hearing, string, error) {
	var h Hearing
	var caseID string
	var prev *string
	err := tx.QueryRow(ctx, `
		SELECT h.id::text, h.case_id::text, h.scheduled_at, h.court_name, h.purpose, h.status, h.outcome,
		       h.previous_hearing_id::text, h.created_at
		FROM hearings h
		JOIN cases c ON c.id = h.case_id
		WHERE h.id = $1 AND c.owner_id = $2
		FOR UPDATE`, id, owner).Scan(&h.ID, &caseID, &h.ScheduledAt, &h.CourtName, &h.Purpose, &h.Status, &h.Outcome, &prev, &h.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Hearing{}, "", NotFound()
	}
	if err != nil {
		return Hearing{}, "", err
	}
	if prev != nil && *prev != "" {
		h.PreviousHearingID = prev
	}
	return h, caseID, nil
}

func planReminders(ctx context.Context, tx pgx.Tx, hearingID string, scheduled time.Time, status string) (int, error) {
	if _, err := tx.Exec(ctx, `
		UPDATE reminders SET status = 'cancelled'
		WHERE hearing_id = $1 AND status = 'scheduled'`, hearingID); err != nil {
		return 0, err
	}
	if status != "scheduled" {
		return 0, nil
	}
	n := 0
	for _, slot := range reminders.Slots(scheduled, time.Now()) {
		tag, err := tx.Exec(ctx, `
			INSERT INTO reminders (id, hearing_id, kind, remind_at, status, attempts)
			VALUES ($1, $2, $3, $4, 'scheduled', 0)
			ON CONFLICT (hearing_id, kind) DO UPDATE
			SET remind_at = EXCLUDED.remind_at,
			    status = 'scheduled',
			    attempts = 0,
			    last_error = '',
			    sent_at = NULL,
			    claimed_at = NULL
			WHERE reminders.status = 'cancelled'`,
			uuid.NewString(), hearingID, slot.Kind, slot.At)
		if err != nil {
			return 0, err
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

func cancelScheduledReminders(ctx context.Context, tx pgx.Tx, hearingID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE reminders SET status = 'cancelled'
		WHERE hearing_id = $1 AND status = 'scheduled'`, hearingID)
	return err
}

func recomputeNext(ctx context.Context, tx pgx.Tx, caseID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE cases
		SET next_hearing_at = (
			SELECT MIN(scheduled_at) FROM hearings
			WHERE case_id = $1 AND status = 'scheduled' AND scheduled_at > now()
		),
		updated_at = now()
		WHERE id = $1`, caseID)
	return err
}
