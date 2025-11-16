package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrNotFound            = NewAppError(http.StatusNotFound, "Resource not found", nil)
	ErrBadRequest          = NewAppError(http.StatusBadRequest, "Invalid request", nil)
	ErrUnauthorized        = NewAppError(http.StatusUnauthorized, "Unauthorized", nil)
	ErrForbidden           = NewAppError(http.StatusForbidden, "Forbidden", nil)
	ErrConflict            = NewAppError(http.StatusConflict, "Resource already exists", nil)
	ErrInternalServerError = NewAppError(http.StatusInternalServerError, "Internal server error", nil)
	ErrEmailExists         = NewAppError(http.StatusConflict, "Email already exists", nil)
	ErrInvalidCredentials  = NewAppError(http.StatusUnauthorized, "Invalid email or password", nil)
	ErrUserNotFound        = NewAppError(http.StatusNotFound, "User not found", nil)
	ErrBookNotFound        = NewAppError(http.StatusNotFound, "Book not found", nil)
)

// Wrap wraps an error with additional context
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	// Check if it's a predefined error
	if err == ErrNotFound || err == ErrBadRequest || err == ErrUnauthorized ||
		err == ErrForbidden || err == ErrConflict || err == ErrInternalServerError ||
		err == ErrEmailExists || err == ErrInvalidCredentials || err == ErrUserNotFound ||
		err == ErrBookNotFound {
		return err.(*AppError)
	}

	// Try to unwrap if it's already an AppError
	if unwrapped, ok := err.(*AppError); ok {
		return NewAppError(unwrapped.Code, message, err)
	}

	return NewAppError(http.StatusInternalServerError, message, err)
}
