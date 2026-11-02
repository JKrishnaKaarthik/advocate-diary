package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const cookieName = "diary_session"

type Manager struct {
	secret []byte
	secure bool
	ttl    time.Duration
}

func New(secret string, secure bool) *Manager {
	return &Manager{secret: []byte(secret), secure: secure, ttl: 14 * 24 * time.Hour}
}

func (m *Manager) Set(w http.ResponseWriter, userID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    m.sign(userID, time.Now().Add(m.ttl)),
		Path:     "/",
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(m.ttl.Seconds()),
	})
}

func (m *Manager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (m *Manager) UserID(r *http.Request) (string, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	id, exp, err := m.open(c.Value)
	if err != nil {
		return "", err
	}
	if time.Now().After(exp) {
		return "", errors.New("session expired")
	}
	if _, err := uuid.Parse(id); err != nil {
		return "", err
	}
	return id, nil
}

func (m *Manager) sign(userID string, exp time.Time) string {
	payload := userID + "|" + strconv.FormatInt(exp.Unix(), 10)
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (m *Manager) open(value string) (string, time.Time, error) {
	payloadB64, macB64, ok := strings.Cut(value, ".")
	if !ok {
		return "", time.Time{}, errors.New("bad session")
	}
	payload, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return "", time.Time{}, err
	}
	got, err := base64.RawURLEncoding.DecodeString(macB64)
	if err != nil {
		return "", time.Time{}, err
	}
	mac := hmac.New(sha256.New, m.secret)
	mac.Write(payload)
	if !hmac.Equal(got, mac.Sum(nil)) {
		return "", time.Time{}, errors.New("bad session")
	}
	id, expText, ok := strings.Cut(string(payload), "|")
	if !ok {
		return "", time.Time{}, errors.New("bad session")
	}
	unix, err := strconv.ParseInt(expText, 10, 64)
	if err != nil {
		return "", time.Time{}, err
	}
	return id, time.Unix(unix, 0), nil
}
