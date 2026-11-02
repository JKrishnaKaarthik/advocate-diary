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

type CaseSummary struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	CaseNumber    string     `json:"case_number"`
	CourtName     string     `json:"court_name"`
	CaseType      string     `json:"case_type"`
	Status        string     `json:"status"`
	OppositeParty string     `json:"opposite_party"`
	FilingDate    *string    `json:"filing_date"`
	NextHearingAt *time.Time `json:"next_hearing_at"`
	Notes         string     `json:"notes"`
	ClientID      string     `json:"client_id"`
	ClientName    string     `json:"client_name"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CaseDetail struct {
	CaseSummary
	Documents []Document `json:"documents"`
	Hearings  []Hearing  `json:"hearings"`
}

type CaseInput struct {
	ClientID      string
	Title         string
	CourtName     string
	CaseType      string
	CaseNumber    string
	OppositeParty string
	FilingDate    string
	Notes         string
}

type CasePatch struct {
	ClientID      *string
	Title         *string
	CourtName     *string
	CaseType      *string
	Status        *string
	CaseNumber    *string
	OppositeParty *string
	FilingDate    *string
	Notes         *string
}

func (s *Store) CreateCase(ctx context.Context, owner string, in CaseInput) (CaseSummary, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.CourtName = strings.TrimSpace(in.CourtName)
	if in.Title == "" || in.CourtName == "" || !validate.CaseType(in.CaseType) {
		return CaseSummary{}, Bad("Title, court, and a valid case type are required")
	}
	filing, err := parseDate(in.FilingDate)
	if err != nil {
		return CaseSummary{}, Bad("Filing date must be YYYY-MM-DD")
	}
	var owns bool
	if err := s.pool.QueryRow(ctx, `SELECT true FROM clients WHERE id = $1 AND owner_id = $2`, in.ClientID, owner).Scan(&owns); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CaseSummary{}, NotFound()
		}
		return CaseSummary{}, err
	}
	id := uuid.NewString()
	row := s.pool.QueryRow(ctx, `
		INSERT INTO cases (
			id, owner_id, client_id, title, case_number, court_name, case_type, status, opposite_party, filing_date, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', $8, $9, $10)
		RETURNING id`, id, owner, in.ClientID, in.Title, strings.TrimSpace(in.CaseNumber), in.CourtName, in.CaseType, strings.TrimSpace(in.OppositeParty), filing, strings.TrimSpace(in.Notes))
	if err := row.Scan(&id); err != nil {
		return CaseSummary{}, err
	}
	return s.getCaseSummary(ctx, owner, id)
}

func (s *Store) ListCases(ctx context.Context, owner, status, clientID, q string) ([]CaseSummary, error) {
	if status != "" && !validate.CaseStatus(status) {
		return nil, Bad("Status is not valid")
	}
	query := `
		SELECT c.id::text, c.title, c.case_number, c.court_name, c.case_type, c.status,
		       c.opposite_party, c.filing_date, c.next_hearing_at, c.notes, c.client_id::text,
		       cl.full_name, c.created_at, c.updated_at
		FROM cases c
		JOIN clients cl ON cl.id = c.client_id
		WHERE c.owner_id = $1`
	args := []any{owner}
	if status != "" {
		args = append(args, status)
		query += " AND c.status = $2"
	}
	n := len(args)
	if clientID != "" {
		n++
		args = append(args, clientID)
		query += " AND c.client_id = $" + itoa(n)
	}
	if strings.TrimSpace(q) != "" {
		n++
		args = append(args, like(q))
		query += " AND (c.title ILIKE $" + itoa(n) + " ESCAPE '\\' OR c.case_number ILIKE $" + itoa(n) + " ESCAPE '\\' OR cl.full_name ILIKE $" + itoa(n) + " ESCAPE '\\')"
	}
	query += " ORDER BY c.next_hearing_at NULLS LAST, c.title"
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CaseSummary
	for rows.Next() {
		item, err := scanCase(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if out == nil {
		out = []CaseSummary{}
	}
	return out, rows.Err()
}

func (s *Store) GetCase(ctx context.Context, owner, id string) (CaseDetail, error) {
	summary, err := s.getCaseSummary(ctx, owner, id)
	if err != nil {
		return CaseDetail{}, err
	}
	docs, err := s.listDocuments(ctx, owner, id)
	if err != nil {
		return CaseDetail{}, err
	}
	hearings, err := s.listHearings(ctx, owner, id)
	if err != nil {
		return CaseDetail{}, err
	}
	return CaseDetail{CaseSummary: summary, Documents: docs, Hearings: hearings}, nil
}

func (s *Store) UpdateCase(ctx context.Context, owner, id string, patch CasePatch) (CaseSummary, error) {
	current, err := s.getCaseSummary(ctx, owner, id)
	if err != nil {
		return CaseSummary{}, err
	}
	if patch.Title != nil {
		current.Title = strings.TrimSpace(*patch.Title)
	}
	if patch.CourtName != nil {
		current.CourtName = strings.TrimSpace(*patch.CourtName)
	}
	if patch.CaseType != nil {
		current.CaseType = *patch.CaseType
	}
	if patch.Status != nil {
		current.Status = *patch.Status
	}
	if patch.CaseNumber != nil {
		current.CaseNumber = strings.TrimSpace(*patch.CaseNumber)
	}
	if patch.OppositeParty != nil {
		current.OppositeParty = strings.TrimSpace(*patch.OppositeParty)
	}
	if patch.Notes != nil {
		current.Notes = strings.TrimSpace(*patch.Notes)
	}
	if patch.ClientID != nil {
		current.ClientID = *patch.ClientID
	}
	filingText := ""
	if current.FilingDate != nil {
		filingText = *current.FilingDate
	}
	if patch.FilingDate != nil {
		filingText = strings.TrimSpace(*patch.FilingDate)
	}
	if current.Title == "" || current.CourtName == "" || !validate.CaseType(current.CaseType) || !validate.CaseStatus(current.Status) {
		return CaseSummary{}, Bad("Title, court, case type, and status are required")
	}
	filing, err := parseDate(filingText)
	if err != nil {
		return CaseSummary{}, Bad("Filing date must be YYYY-MM-DD")
	}
	var owns bool
	if err := s.pool.QueryRow(ctx, `SELECT true FROM clients WHERE id = $1 AND owner_id = $2`, current.ClientID, owner).Scan(&owns); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CaseSummary{}, NotFound()
		}
		return CaseSummary{}, err
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE cases
		SET client_id = $3, title = $4, case_number = $5, court_name = $6, case_type = $7,
		    status = $8, opposite_party = $9, filing_date = $10, notes = $11, updated_at = now()
		WHERE id = $1 AND owner_id = $2`,
		id, owner, current.ClientID, current.Title, current.CaseNumber, current.CourtName, current.CaseType,
		current.Status, current.OppositeParty, filing, current.Notes)
	if err != nil {
		return CaseSummary{}, err
	}
	if tag.RowsAffected() == 0 {
		return CaseSummary{}, NotFound()
	}
	return s.getCaseSummary(ctx, owner, id)
}

func (s *Store) casesForClient(ctx context.Context, owner, clientID string) ([]CaseSummary, error) {
	return s.ListCases(ctx, owner, "", clientID, "")
}

func (s *Store) getCaseSummary(ctx context.Context, owner, id string) (CaseSummary, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT c.id::text, c.title, c.case_number, c.court_name, c.case_type, c.status,
		       c.opposite_party, c.filing_date, c.next_hearing_at, c.notes, c.client_id::text,
		       cl.full_name, c.created_at, c.updated_at
		FROM cases c
		JOIN clients cl ON cl.id = c.client_id
		WHERE c.id = $1 AND c.owner_id = $2`, id, owner)
	if err != nil {
		return CaseSummary{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return CaseSummary{}, NotFound()
	}
	return scanCase(rows)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCase(row scanner) (CaseSummary, error) {
	var c CaseSummary
	var filing *time.Time
	err := row.Scan(&c.ID, &c.Title, &c.CaseNumber, &c.CourtName, &c.CaseType, &c.Status,
		&c.OppositeParty, &filing, &c.NextHearingAt, &c.Notes, &c.ClientID, &c.ClientName, &c.CreatedAt, &c.UpdatedAt)
	c.FilingDate = dateString(filing)
	return c, err
}

func parseDate(v string) (any, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
