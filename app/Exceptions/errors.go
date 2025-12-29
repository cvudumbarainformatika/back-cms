package exceptions

import "fmt"

// AppError represents a custom application error
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrNotFound         = &AppError{Code: 404, Message: "Resource not found"}
	ErrUnauthorized     = &AppError{Code: 401, Message: "Unauthorized"}
	ErrForbidden        = &AppError{Code: 403, Message: "Forbidden"}
	ErrBadRequest       = &AppError{Code: 400, Message: "Bad request"}
	ErrInternalServer   = &AppError{Code: 500, Message: "Internal server error"}
	ErrValidationFailed = &AppError{Code: 422, Message: "Validation failed"}
)
