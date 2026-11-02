package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/app"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/session"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/storage"
)

type captureMail struct{ n int }

func (c *captureMail) Send(context.Context, string, string, string) error {
	c.n++
	return nil
}

func TestPracticeFlow(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	objects := storage.New(os.Getenv("S3_ENDPOINT"), env("S3_REGION", "us-east-1"), os.Getenv("S3_BUCKET"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), os.Getenv("S3_PATH_STYLE") == "true")
	store := app.NewStore(pool, objects)
	if err := store.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	handler := New(store, session.New("test-session-secret", false))
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	a := newClient(t, srv.URL)
	b := newClient(t, srv.URL)
	stamp := time.Now().UnixNano()
	suffix := strconv.FormatInt(stamp, 10)
	lawyerA := a.do(t, http.MethodPost, "/api/sign-up", map[string]string{
		"name": "Asha Rao", "email": "asha+" + suffix + "@example.com", "password": "password1",
	}, http.StatusCreated)
	_ = b.do(t, http.MethodPost, "/api/sign-up", map[string]string{
		"name": "Bharat Shah", "email": "bharat+" + suffix + "@example.com", "password": "password2",
	}, http.StatusCreated)
	if lawyerA["name"] != "Asha Rao" {
		t.Fatalf("signup body %#v", lawyerA)
	}

	client := a.do(t, http.MethodPost, "/api/clients", map[string]string{"full_name": "Meera Iyer", "phone": "99999"}, http.StatusCreated)
	caseItem := a.do(t, http.MethodPost, "/api/cases", map[string]string{
		"client_id": client["id"].(string), "title": "Iyer vs State", "court_name": "City Civil Court", "case_type": "civil",
	}, http.StatusCreated)
	caseID := caseItem["id"].(string)
	b.do(t, http.MethodGet, "/api/clients/"+client["id"].(string), nil, http.StatusNotFound)
	b.do(t, http.MethodGet, "/api/cases/"+caseID, nil, http.StatusNotFound)

	when := time.Now().Add(72 * time.Hour).UTC().Format(time.RFC3339)
	hearing := a.do(t, http.MethodPost, "/api/cases/"+caseID+"/hearings", map[string]string{
		"scheduled_at": when, "court_name": "City Civil Court", "purpose": "arguments",
	}, http.StatusCreated)
	if int(hearing["reminders_planned"].(float64)) != 2 {
		t.Fatalf("reminders_planned = %v", hearing["reminders_planned"])
	}
	past := a.do(t, http.MethodPost, "/api/cases/"+caseID+"/hearings", map[string]string{
		"scheduled_at": time.Now().Add(-2 * time.Hour).UTC().Format(time.RFC3339),
		"purpose":      "mention",
	}, http.StatusCreated)
	if int(past["reminders_planned"].(float64)) != 0 {
		t.Fatalf("past hearing planned %v", past["reminders_planned"])
	}

	nextAt := time.Now().Add(96 * time.Hour).UTC().Format(time.RFC3339)
	adjourned := a.do(t, http.MethodPost, "/api/hearings/"+hearing["id"].(string)+"/adjourn", map[string]string{
		"next_at": nextAt, "outcome": "Adjourned for documents",
	}, http.StatusCreated)
	if adjourned["status"] != "scheduled" {
		t.Fatalf("new hearing %#v", adjourned)
	}

	from := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	to := time.Now().Add(200 * time.Hour).UTC().Format(time.RFC3339)
	diary := a.list(t, "/api/diary?from="+from+"&to="+to)
	if len(diary) == 0 {
		t.Fatal("diary is empty")
	}
	otherDiary := b.list(t, "/api/diary?from="+from+"&to="+to)
	if len(otherDiary) != 0 {
		t.Fatal("other lawyer saw a hearing")
	}

	var sent captureMail
	if _, err := pool.Exec(ctx, `UPDATE reminders SET remind_at = now() - interval '1 minute' WHERE hearing_id = $1`, adjourned["id"]); err != nil {
		t.Fatal(err)
	}
	if err := store.Tick(ctx, &sent); err != nil {
		t.Fatal(err)
	}
	if sent.n == 0 {
		t.Fatal("no reminder email")
	}
	notes := a.list(t, "/api/notifications")
	if len(notes) == 0 || !strings.Contains(notes[0]["title"].(string), "Iyer vs State") {
		t.Fatalf("notifications %#v", notes)
	}
	again := captureMail{}
	if err := store.Tick(ctx, &again); err != nil {
		t.Fatal(err)
	}
	if again.n != 0 {
		t.Fatal("reminder sent twice")
	}
	a.do(t, http.MethodPatch, "/api/cases/"+caseID, map[string]string{"status": "disposed"}, http.StatusOK)
	a.doStatus(t, http.MethodPost, "/api/cases/"+caseID+"/hearings", map[string]string{
		"scheduled_at": time.Now().Add(120 * time.Hour).UTC().Format(time.RFC3339),
		"purpose":      "orders",
	}, http.StatusConflict)

	if os.Getenv("S3_BUCKET") == "" {
		return
	}
	doc := a.upload(t, "/api/cases/"+caseID+"/documents", "note.txt", "text/plain", "petition", []byte("grounds of petition"))
	got := a.do(t, http.MethodGet, "/api/documents/"+doc["id"].(string), nil, http.StatusOK)
	if !strings.Contains(got["url"].(string), "note.txt") && got["url"] == "" {
		t.Fatalf("download %#v", got)
	}
	b.do(t, http.MethodGet, "/api/documents/"+doc["id"].(string), nil, http.StatusNotFound)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type apiClient struct {
	base string
	http *http.Client
}

func newClient(t *testing.T, base string) *apiClient {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	return &apiClient{base: base, http: &http.Client{Jar: jar}}
}

func (c *apiClient) do(t *testing.T, method, path string, body map[string]string, want int) map[string]any {
	t.Helper()
	res := c.raw(t, method, path, body)
	defer res.Body.Close()
	if res.StatusCode != want {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("%s %s = %d %s", method, path, res.StatusCode, b)
	}
	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

func (c *apiClient) doStatus(t *testing.T, method, path string, body map[string]string, want int) {
	t.Helper()
	res := c.raw(t, method, path, body)
	defer res.Body.Close()
	if res.StatusCode != want {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("%s %s = %d %s, want %d", method, path, res.StatusCode, b, want)
	}
}

func (c *apiClient) list(t *testing.T, path string) []map[string]any {
	t.Helper()
	res := c.raw(t, http.MethodGet, path, nil)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("GET %s = %d %s", path, res.StatusCode, b)
	}
	var out []map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

func (c *apiClient) raw(t *testing.T, method, path string, body map[string]string) *http.Response {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.base+path, r)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.http.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func (c *apiClient) upload(t *testing.T, path, name, contentType, label string, content []byte) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("label", label); err != nil {
		t.Fatal(err)
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+name+`"`)
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	w.Close()
	req, err := http.NewRequest(http.MethodPost, c.base+path, &buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := c.http.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("upload = %d %s", res.StatusCode, b)
	}
	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}
