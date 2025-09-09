package main

import (
	"fmt"
	"html/template"
	"log"

	// "os" // Removed as it's no longer directly used

	"programming_blog_go/config"
	"programming_blog_go/internal/adapter/handler"
	"programming_blog_go/internal/adapter/persistence/postgres"
	"programming_blog_go/internal/adapter/service"
	"programming_blog_go/internal/middleware"
	"programming_blog_go/internal/usecase"

	"github.com/gin-gonic/gin"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the database schema
	// NOTE: In a production environment, you would use a dedicated migration tool
	// like `golang-migrate/migrate` or `gorm.io/gorm/migrator` with versioned migrations.
	// For simplicity in this project, we are using GORM's AutoMigrate.
	// err = db.AutoMigrate(&domain.Blog{}, &domain.Category{}, &domain.User{})
	// if err != nil {
	// 	log.Fatalf("Failed to auto migrate database: %v", err)
	// }

	// Initialize repositories
	categoryRepo := postgres.NewCategoryRepository(db)
	blogRepo := postgres.NewBlogRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Initialize mailer service
	mailer := service.NewSMTPSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		cfg.SMTPFrom,
	)

	// Initialize use cases
	getBlogPostsUC := &usecase.GetBlogPostsUseCase{BlogRepository: blogRepo}
	getBlogPostsByCategoryUC := &usecase.GetBlogPostsByCategoryUseCase{BlogRepository: blogRepo, CategoryRepository: categoryRepo}
	getBlogPostBySlugUC := &usecase.GetBlogPostBySlugUseCase{BlogRepository: blogRepo}
	createBlogPostUC := &usecase.CreateBlogPostUseCase{BlogRepository: blogRepo, CategoryRepository: categoryRepo}
	registerUserUC := &usecase.RegisterUserUseCase{UserRepository: userRepo}
	authenticateUserUC := &usecase.AuthenticateUserUseCase{UserRepository: userRepo}
	sendContactMessageUC := &usecase.SendContactMessageUseCase{MailerService: mailer}
	getAllCategoriesUC := &usecase.GetAllCategoriesUseCase{CategoryRepository: categoryRepo}

	// Initialize handlers
	blogHandler := handler.NewBlogHandler(getBlogPostsUC, getBlogPostsByCategoryUC, getBlogPostBySlugUC, createBlogPostUC)
	userHandler := handler.NewUserHandler(registerUserUC, authenticateUserUC)
	contactHandler := handler.NewContactHandler(sendContactMessageUC)

	// Set up Gin router
	r := gin.Default()

	// Load templates
	r.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.html")))

	// Serve static files
	r.Static("/static", "./web/static") // Assuming static files are in web/static

	// Apply CategoryContextMiddleware to all routes that render HTML
	htmlRoutes := r.Group("/")
	htmlRoutes.Use(middleware.CategoryContextMiddleware(getAllCategoriesUC))
	{
		htmlRoutes.GET("/", blogHandler.GetBlogPosts)
		htmlRoutes.GET("/post/:post_slug", blogHandler.GetBlogPost)
		htmlRoutes.GET("/category/:cat_slug", blogHandler.GetBlogPostsByCategory)

		// Pages that serve HTML forms (for now, these are simple renders)
		htmlRoutes.GET("/addpage", blogHandler.AddPostPage)
		htmlRoutes.GET("/register", userHandler.ShowRegisterPage)
		htmlRoutes.GET("/login", userHandler.ShowLoginPage)
		htmlRoutes.GET("/contact", contactHandler.ShowContactPage)
	}

	// API endpoints
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.RegisterUser)
		api.POST("/login", userHandler.LoginUser)
		api.POST("/contact", contactHandler.SendContactMessage)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware([]byte(cfg.JWTSecret)))
		{
			protected.POST("/posts", blogHandler.CreateBlogPost)
			// TODO: Add other protected routes here (e.g., update/delete posts)
		}
	}

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server listening on :%s", port)
	log.Fatal(r.Run(":" + port))
}
