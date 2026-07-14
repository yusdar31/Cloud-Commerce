package errors

import (
	"errors"
	"fmt"
)

// Common domain errors that can be used across services.
var (
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrValidation     = errors.New("validation error")
	ErrInternalServer = errors.New("internal server error")
	ErrConflict       = errors.New("conflict")
)

// NotFound creates a not-found error with context.
func NotFound(entity, id string) error {
	return fmt.Errorf("%s with id %s: %w", entity, id, ErrNotFound)
}

// AlreadyExists creates a conflict error with context.
func AlreadyExists(entity, key string) error {
	return fmt.Errorf("%s with %s already exists: %w", entity, key, ErrAlreadyExists)
}

// IsNotFound checks if the error is a not-found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists checks if the error is an already-exists error.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsUnauthorized checks if the error is an unauthorized error.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden checks if the error is a forbidden error.
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsValidation checks if the error is a validation error.
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsConflict checks if the error is a conflict error.
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}
