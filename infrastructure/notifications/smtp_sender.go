package notifications

import (
	"context"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"net/smtp"
	"strings"
)

type SmtpSender struct {
	host     string
	port     string
	email    string
	password string
}

func NewSmtpSender(host, port, email, password string) ports.NotificationSender {
	return &SmtpSender{
		host:     host,
		port:     port,
		email:    email,
		password: password,
	}
}

func (s *SmtpSender) Send(ctx context.Context, data domain.EmailData) error {
	auth := smtp.PlainAuth("", s.email, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.email))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", data.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))
	msg.WriteString("MIME-version: 1.0;\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	msg.WriteString(data.BodyHTML)

	toList := []string{data.To}

	err := smtp.SendMail(addr, auth, s.email, toList, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
