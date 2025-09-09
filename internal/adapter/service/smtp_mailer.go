package service

import (
	"fmt"
	"net/smtp"
)

// SMTPSender implements domain.MailerService for sending emails via SMTP.
type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// NewSMTPSender creates a new SMTPSender instance.
func NewSMTPSender(host, port, username, password, from string) *SMTPSender {
	return &SMTPSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendEmail sends an email using SMTP.
func (s *SMTPSender) SendEmail(to []string, subject, body string) error {
	// Authentication.
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// Email content.
	msg := []byte(
		"To: " + to[0] + "\r\n" +
			"From: " + s.From + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	// Sending email.
	err := smtp.SendMail(s.Host+":"+s.Port, auth, s.From, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
