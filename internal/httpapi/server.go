package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/app"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/session"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/validate"
)

type Server struct {
	store    *app.Store
	sessions *session.Manager
}

func New(store *app.Store, sessions *session.Manager) http.Handler {
	mux := http.NewServeMux()
	Register(mux, store, sessions)
	return mux
}

func Register(mux *http.ServeMux, store *app.Store, sessions *session.Manager) {
	s := &Server{store: store, sessions: sessions}
	mux.HandleFunc("POST /api/sign-up", s.signUp)
	mux.HandleFunc("POST /api/session", s.signIn)
	mux.HandleFunc("DELETE /api/session", s.signOut)
	mux.HandleFunc("GET /api/session", s.withUser(s.current))
	mux.HandleFunc("GET /api/clients", s.withUser(s.listClients))
	mux.HandleFunc("POST /api/clients", s.withUser(s.createClient))
	mux.HandleFunc("GET /api/clients/{id}", s.withUser(s.getClient))
	mux.HandleFunc("PATCH /api/clients/{id}", s.withUser(s.patchClient))
	mux.HandleFunc("POST /api/cases", s.withUser(s.createCase))
	mux.HandleFunc("GET /api/cases", s.withUser(s.listCases))
	mux.HandleFunc("GET /api/cases/{id}", s.withUser(s.getCase))
	mux.HandleFunc("PATCH /api/cases/{id}", s.withUser(s.patchCase))
	mux.HandleFunc("POST /api/cases/{id}/documents", s.withUser(s.upload))
	mux.HandleFunc("GET /api/documents/{id}", s.withUser(s.download))
	mux.HandleFunc("DELETE /api/documents/{id}", s.withUser(s.deleteDocument))
	mux.HandleFunc("POST /api/cases/{id}/hearings", s.withUser(s.addHearing))
	mux.HandleFunc("POST /api/hearings/{id}/complete", s.withUser(s.completeHearing))
	mux.HandleFunc("POST /api/hearings/{id}/cancel", s.withUser(s.cancelHearing))
	mux.HandleFunc("POST /api/hearings/{id}/adjourn", s.withUser(s.adjournHearing))
	mux.HandleFunc("GET /api/diary", s.withUser(s.diary))
	mux.HandleFunc("GET /api/notifications", s.withUser(s.notifications))
	mux.HandleFunc("POST /api/notifications/{id}/read", s.withUser(s.readNotification))
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
	})
}

func (s *Server) withUser(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.sessions.UserID(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Sign in required"})
			return
		}
		next(w, r, id)
	}
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	lawyer, err := s.store.SignUp(r.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	s.sessions.Set(w, lawyer.ID)
	writeJSON(w, http.StatusCreated, lawyer)
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	lawyer, err := s.store.SignIn(r.Context(), body.Email, body.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	s.sessions.Set(w, lawyer.ID)
	writeJSON(w, http.StatusOK, lawyer)
}

func (s *Server) signOut(w http.ResponseWriter, r *http.Request) {
	s.sessions.Clear(w)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) current(w http.ResponseWriter, r *http.Request, userID string) {
	lawyer, err := s.store.Lawyer(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, lawyer)
}

func (s *Server) listClients(w http.ResponseWriter, r *http.Request, userID string) {
	items, err := s.store.ListClients(r.Context(), userID, r.URL.Query().Get("q"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) createClient(w http.ResponseWriter, r *http.Request, userID string) {
	var body struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Address  string `json:"address"`
		Notes    string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	item, err := s.store.CreateClient(r.Context(), userID, app.ClientInput{
		FullName: body.FullName, Phone: body.Phone, Email: body.Email, Address: body.Address, Notes: body.Notes,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) getClient(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, err := s.store.GetClient(r.Context(), userID, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) patchClient(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var body struct {
		FullName *string `json:"full_name"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		Address  *string `json:"address"`
		Notes    *string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	item, err := s.store.UpdateClient(r.Context(), userID, id, app.ClientPatch{
		FullName: body.FullName, Phone: body.Phone, Email: body.Email, Address: body.Address, Notes: body.Notes,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) createCase(w http.ResponseWriter, r *http.Request, userID string) {
	var body struct {
		ClientID      string `json:"client_id"`
		Title         string `json:"title"`
		CourtName     string `json:"court_name"`
		CaseType      string `json:"case_type"`
		CaseNumber    string `json:"case_number"`
		OppositeParty string `json:"opposite_party"`
		FilingDate    string `json:"filing_date"`
		Notes         string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	if _, err := uuid.Parse(body.ClientID); err != nil {
		writeErr(w, app.NotFound())
		return
	}
	item, err := s.store.CreateCase(r.Context(), userID, app.CaseInput{
		ClientID: body.ClientID, Title: body.Title, CourtName: body.CourtName, CaseType: body.CaseType,
		CaseNumber: body.CaseNumber, OppositeParty: body.OppositeParty, FilingDate: body.FilingDate, Notes: body.Notes,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) listCases(w http.ResponseWriter, r *http.Request, userID string) {
	clientID := r.URL.Query().Get("clientId")
	if clientID != "" {
		if _, err := uuid.Parse(clientID); err != nil {
			writeErr(w, app.Bad("Client is not valid"))
			return
		}
	}
	items, err := s.store.ListCases(r.Context(), userID, r.URL.Query().Get("status"), clientID, r.URL.Query().Get("q"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) getCase(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, err := s.store.GetCase(r.Context(), userID, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) patchCase(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var body struct {
		ClientID      *string `json:"client_id"`
		Title         *string `json:"title"`
		CourtName     *string `json:"court_name"`
		CaseType      *string `json:"case_type"`
		Status        *string `json:"status"`
		CaseNumber    *string `json:"case_number"`
		OppositeParty *string `json:"opposite_party"`
		FilingDate    *string `json:"filing_date"`
		Notes         *string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	item, err := s.store.UpdateCase(r.Context(), userID, id, app.CasePatch{
		ClientID: body.ClientID, Title: body.Title, CourtName: body.CourtName, CaseType: body.CaseType,
		Status: body.Status, CaseNumber: body.CaseNumber, OppositeParty: body.OppositeParty,
		FilingDate: body.FilingDate, Notes: body.Notes,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, validate.MaxBytes+1<<20)
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		writeErr(w, app.Bad("Upload could not be read. Files must be 20 MB or smaller."))
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeErr(w, app.Bad("Choose a file to upload"))
		return
	}
	defer file.Close()
	body, err := io.ReadAll(file)
	if err != nil {
		writeErr(w, app.Bad("Upload could not be read. Files must be 20 MB or smaller."))
		return
	}
	doc, err := s.store.Upload(r.Context(), userID, id, header.Filename, header.Header.Get("Content-Type"), r.FormValue("label"), body)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, doc)
}

func (s *Server) download(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, err := s.store.Download(r.Context(), userID, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) deleteDocument(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := s.store.DeleteDocument(r.Context(), userID, id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) addHearing(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var body struct {
		ScheduledAt string `json:"scheduled_at"`
		CourtName   string `json:"court_name"`
		Purpose     string `json:"purpose"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	at, err := time.Parse(time.RFC3339, body.ScheduledAt)
	if err != nil {
		writeErr(w, app.Bad("Hearing time is not a valid timestamp"))
		return
	}
	item, err := s.store.AddHearing(r.Context(), userID, id, body.CourtName, body.Purpose, at)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) completeHearing(w http.ResponseWriter, r *http.Request, userID string) {
	s.withOutcome(w, r, func(outcome string) (any, error) {
		return s.store.CompleteHearing(r.Context(), userID, r.PathValue("id"), outcome)
	})
}

func (s *Server) cancelHearing(w http.ResponseWriter, r *http.Request, userID string) {
	s.withOutcome(w, r, func(outcome string) (any, error) {
		return s.store.CancelHearing(r.Context(), userID, r.PathValue("id"), outcome)
	})
}

func (s *Server) withOutcome(w http.ResponseWriter, r *http.Request, fn func(string) (any, error)) {
	if _, err := pathID(r); err != nil {
		writeErr(w, err)
		return
	}
	var body struct {
		Outcome string `json:"outcome"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	item, err := fn(body.Outcome)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) adjournHearing(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var body struct {
		NextAt  string `json:"next_at"`
		Outcome string `json:"outcome"`
		Purpose string `json:"purpose"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	nextAt, err := time.Parse(time.RFC3339, body.NextAt)
	if err != nil {
		writeErr(w, app.Bad("The next date is not a valid timestamp"))
		return
	}
	item, err := s.store.AdjournHearing(r.Context(), userID, id, body.Outcome, body.Purpose, nextAt)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) diary(w http.ResponseWriter, r *http.Request, userID string) {
	from, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	if err != nil {
		writeErr(w, app.Bad("Diary dates are not valid"))
		return
	}
	to, err := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if err != nil {
		writeErr(w, app.Bad("Diary dates are not valid"))
		return
	}
	items, err := s.store.Diary(r.Context(), userID, from, to)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) notifications(w http.ResponseWriter, r *http.Request, userID string) {
	items, err := s.store.Notifications(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) readNotification(w http.ResponseWriter, r *http.Request, userID string) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := s.store.MarkRead(r.Context(), userID, id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		return "", app.NotFound()
	}
	return id, nil
}

func readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return app.Bad("Request body is not valid")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeErr(w http.ResponseWriter, err error) {
	var se *app.StatusError
	if errors.As(err, &se) {
		writeJSON(w, se.Status, map[string]string{"error": se.Msg})
		return
	}
	log.Printf("request failed: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
}
