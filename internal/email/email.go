// Package email sends transactional emails (password reset, verification).
//
// In local development the default Sender just logs the message (including any
// link) to the server log, so the whole flow can be tested without a real mail
// server. Set SMTP_HOST to send through a real SMTP server (e.g. Mailpit at
// localhost:1025, or a production provider).
package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
)

// Message is one email to send.
type Message struct {
	To      string
	Subject string
	HTML    string
}

// Sender sends an email message.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// LogSender "sends" emails by logging them. Perfect for local development:
// you read the reset/verification link straight from the server log.
type LogSender struct{}

func (LogSender) Send(_ context.Context, msg Message) error {
	slog.Info("email (dev: not actually sent)",
		"to", msg.To, "subject", msg.Subject, "body", msg.HTML)
	return nil
}

// SMTPSender sends real emails over SMTP (plain, for local Mailpit or a
// provider). From is the sender address; Addr is host:port.
type SMTPSender struct {
	Addr string
	From string
	Auth smtp.Auth // nil for servers without auth (e.g. Mailpit)
}

func (s SMTPSender) Send(_ context.Context, msg Message) error {
	body := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.From, msg.To, msg.Subject, msg.HTML,
	)
	if err := smtp.SendMail(s.Addr, s.Auth, s.From, []string{msg.To}, []byte(body)); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
