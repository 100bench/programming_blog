package domain

// MailerService defines the interface for sending emails.
type MailerService interface {
	SendEmail(to []string, subject, body string) error
}
