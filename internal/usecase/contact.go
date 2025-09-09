package usecase

import (
	"programming_blog_go/internal/domain"
)

// SendContactMessageUseCase handles sending contact form emails.
type SendContactMessageUseCase struct {
	MailerService domain.MailerService
}

type SendContactMessageRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Content string `json:"content" binding:"required"`
}

func (uc *SendContactMessageUseCase) Execute(req SendContactMessageRequest) error {
	subject := "Contact Form Message from " + req.Name
	body := "From: " + req.Email + "\n\n" + req.Content

	// In a real application, you might want to send to a specific admin email configured in settings.
	// For now, let's assume a dummy admin email.
	adminEmail := []string{"admin@example.com"}

	return uc.MailerService.SendEmail(adminEmail, subject, body)
}
