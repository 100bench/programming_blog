package handler

import (
	"net/http"

	"programming_blog_go/internal/domain"
	"programming_blog_go/internal/usecase"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TODO: Move this to config
var jwtSecret = []byte("your_secret_key")

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	RegisterUserUseCase     *usecase.RegisterUserUseCase
	AuthenticateUserUseCase *usecase.AuthenticateUserUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(
	registerUserUC *usecase.RegisterUserUseCase,
	authenticateUserUC *usecase.AuthenticateUserUseCase,
) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase:     registerUserUC,
		AuthenticateUserUseCase: authenticateUserUC,
	}
}

// RegisterUser handles user registration.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req usecase.RegisterUserRequest
	if err := c.ShouldBind(&req); err != nil { // Changed from ShouldBindJSON to ShouldBind
		HandleError(c, domain.ErrInvalidInput) // Map binding errors to InvalidInput
		return
	}

	user, err := h.RegisterUserUseCase.Execute(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

// LoginUser handles user login and generates a JWT token.
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req usecase.AuthenticateUserRequest
	if err := c.ShouldBind(&req); err != nil { // Changed from ShouldBindJSON to ShouldBind
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	user, err := h.AuthenticateUserUseCase.Execute(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// ShowRegisterPage renders the registration form page.
func (h *UserHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "Регистрация"})
}

// ShowLoginPage renders the login form page.
func (h *UserHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "Авторизация"})
}
