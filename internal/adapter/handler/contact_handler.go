package handler

import (
	"net/http"

	"programming_blog_go/internal/domain"
	"programming_blog_go/internal/usecase"

	"github.com/gin-gonic/gin"
)

// ContactHandler handles HTTP requests related to the contact form.
type ContactHandler struct {
	SendContactMessageUseCase *usecase.SendContactMessageUseCase
}

// NewContactHandler creates a new ContactHandler.
func NewContactHandler(sendContactMessageUC *usecase.SendContactMessageUseCase) *ContactHandler {
	return &ContactHandler{SendContactMessageUseCase: sendContactMessageUC}
}

// SendContactMessage handles the request to send a contact message.
func (h *ContactHandler) SendContactMessage(c *gin.Context) {
	var req usecase.SendContactMessageRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	err := h.SendContactMessageUseCase.Execute(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// ShowContactPage renders the contact form page.
// This is already present as a dummy in blog_handler.go, but will be moved here.
func (h *ContactHandler) ShowContactPage(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", gin.H{"title": "Обратная связь"})
}
