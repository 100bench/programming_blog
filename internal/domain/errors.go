package domain

import "errors"

// Predefined errors for common domain scenarios.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
	// Add more domain-specific errors as needed
)
