package infrastructure

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type SMTPEmailService struct {
    Host     string
    Port     string
    Username string
    Password string
    From     string
}

func NewSMTPEmailService() *SMTPEmailService {
	return &SMTPEmailService{
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From: os.Getenv("SMTP_FROM"),
	}
}

func (s *SMTPEmailService) SendEmail(to []string, subject string, body string) error {
    auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

    msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
        "\r\n" + body)

    addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
    return smtp.SendMail(addr, auth, s.From, to, msg)
}