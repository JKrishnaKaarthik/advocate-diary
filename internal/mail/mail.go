package mail

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
)

type Sender interface {
	Send(ctx context.Context, to, subject, body string) error
}

type LogSender struct{}

func (LogSender) Send(_ context.Context, to, subject, body string) error {
	log.Printf("reminder email (no SMTP configured)\nto: %s\nsubject: %s\n%s", to, subject, body)
	return nil
}

type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (s SMTPSender) Send(_ context.Context, to, subject, body string) error {
	subject = strings.NewReplacer("\r", " ", "\n", " ").Replace(subject)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", s.From, to, subject, body)
	var auth smtp.Auth
	if s.Username != "" {
		auth = smtp.PlainAuth("", s.Username, s.Password, s.Host)
	}
	return smtp.SendMail(net.JoinHostPort(s.Host, s.Port), auth, s.From, []string{to}, []byte(msg))
}
