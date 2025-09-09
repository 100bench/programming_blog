package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration settings.
type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
	SMTPFrom   string
	AppPort    string
}

// LoadConfig loads configuration from .env file or environment variables.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file, assuming environment variables are set: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "blogdb"),
		DBPort:     getEnv("DB_PORT", "5432"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretjwtkey"), // Default for development
		SMTPHost:   getEnv("SMTP_HOST", "localhost"),
		SMTPPort:   getEnv("SMTP_PORT", "1025"), // Default Mailhog/Mailtrap local port
		SMTPUser:   getEnv("SMTP_USERNAME", ""),
		SMTPPass:   getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:   getEnv("SMTP_FROM", "noreply@example.com"),
		AppPort:    getEnv("PORT", "8080"),
	}
}

// getEnv retrieves environment variables or provides a fallback default.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
