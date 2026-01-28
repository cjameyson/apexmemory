package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// Severity levels for errors.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// AppError is a structured error with metadata for logging and debugging.
type AppError struct {
	Err             error    // underlying cause
	Message         string   // log message
	Code            string   // e.g., NOTEBOOK_NOT_FOUND, DB_QUERY_FAILED
	HTTPStatus      int      // maps to HTTP response code
	Severity        Severity // low, medium, high, critical
	Retryable       bool
	RemediationHint string
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error for error chain traversal.
func (e *AppError) Unwrap() error {
	return e.Err
}

// Status returns the HTTP status code for this error.
func (e *AppError) Status() int {
	if e.HTTPStatus == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTPStatus
}

// LogValue returns structured log attributes for slog integration.
func (e *AppError) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("error_code", e.Code),
		slog.String("error_message", e.Message),
		slog.String("error_type", errorType(e)),
		slog.String("severity", string(e.Severity)),
		slog.Bool("retryable", e.Retryable),
	}

	if e.RemediationHint != "" {
		attrs = append(attrs, slog.String("remediation_hint", e.RemediationHint))
	}

	if e.Err != nil {
		attrs = append(attrs, slog.String("cause", e.Err.Error()))
	}

	return slog.GroupValue(attrs...)
}

func errorType(err *AppError) string {
	if err.Err != nil {
		return fmt.Sprintf("%T", err.Err)
	}
	return "AppError"
}

// Error constructors - keep codes short and consistent.

// ErrNotFound creates a not found error for a resource.
func ErrNotFound(resource string, err error) *AppError {
	upper := strings.ToUpper(resource)
	return &AppError{
		Err:             err,
		Message:         fmt.Sprintf("%s not found", resource),
		Code:            fmt.Sprintf("%s_NOT_FOUND", upper),
		HTTPStatus:      http.StatusNotFound,
		Severity:        SeverityLow,
		Retryable:       false,
		RemediationHint: fmt.Sprintf("Verify the %s ID exists and belongs to the user", resource),
	}
}

// ErrDBQuery creates an error for database query failures.
func ErrDBQuery(operation string, err error) *AppError {
	return &AppError{
		Err:             err,
		Message:         fmt.Sprintf("database %s failed", operation),
		Code:            "DB_QUERY_FAILED",
		HTTPStatus:      http.StatusInternalServerError,
		Severity:        SeverityHigh,
		Retryable:       true,
		RemediationHint: "Check database connectivity and query syntax",
	}
}

// ErrDBTransaction creates an error for transaction failures.
func ErrDBTransaction(phase string, err error) *AppError {
	return &AppError{
		Err:             err,
		Message:         fmt.Sprintf("transaction %s failed", phase),
		Code:            "DB_TX_FAILED",
		HTTPStatus:      http.StatusInternalServerError,
		Severity:        SeverityHigh,
		Retryable:       true,
		RemediationHint: "Retry the operation; check for deadlocks if persistent",
	}
}

// ErrValidation creates a validation error.
func ErrValidation(field, reason string) *AppError {
	return &AppError{
		Message:         fmt.Sprintf("validation failed: %s - %s", field, reason),
		Code:            "VALIDATION_FAILED",
		HTTPStatus:      http.StatusBadRequest,
		Severity:        SeverityLow,
		Retryable:       false,
		RemediationHint: "Check request payload matches API specification",
	}
}

// ErrUnauthorized creates an unauthorized error.
func ErrUnauthorized(reason string) *AppError {
	return &AppError{
		Message:         reason,
		Code:            "UNAUTHORIZED",
		HTTPStatus:      http.StatusUnauthorized,
		Severity:        SeverityLow,
		Retryable:       false,
		RemediationHint: "Provide valid authentication credentials",
	}
}

// ErrForbidden creates a forbidden error.
func ErrForbidden(reason string) *AppError {
	return &AppError{
		Message:         reason,
		Code:            "FORBIDDEN",
		HTTPStatus:      http.StatusForbidden,
		Severity:        SeverityLow,
		Retryable:       false,
		RemediationHint: "Verify you have permission to access this resource",
	}
}

// ErrInternal creates an internal server error.
func ErrInternal(message string, err error) *AppError {
	return &AppError{
		Err:             err,
		Message:         message,
		Code:            "INTERNAL_ERROR",
		HTTPStatus:      http.StatusInternalServerError,
		Severity:        SeverityCritical,
		Retryable:       false,
		RemediationHint: "Check server logs for details",
	}
}

// ErrConflict creates a conflict error (e.g., duplicate resource).
func ErrConflict(resource string, err error) *AppError {
	upper := strings.ToUpper(resource)
	return &AppError{
		Err:             err,
		Message:         fmt.Sprintf("%s already exists", resource),
		Code:            fmt.Sprintf("%s_CONFLICT", upper),
		HTTPStatus:      http.StatusConflict,
		Severity:        SeverityLow,
		Retryable:       false,
		RemediationHint: fmt.Sprintf("Use a different %s identifier or update existing", resource),
	}
}

// ErrBadRequest creates a generic bad request error.
func ErrBadRequest(message string, err error) *AppError {
	return &AppError{
		Err:             err,
		Message:         message,
		Code:            "BAD_REQUEST",
		HTTPStatus:      http.StatusBadRequest,
		Severity:        SeverityLow,
		Retryable:       false,
		RemediationHint: "Check request format and parameters",
	}
}

// ErrRateLimited creates a rate limit exceeded error.
func ErrRateLimited() *AppError {
	return &AppError{
		Message:         "rate limit exceeded",
		Code:            "RATE_LIMITED",
		HTTPStatus:      http.StatusTooManyRequests,
		Severity:        SeverityLow,
		Retryable:       true,
		RemediationHint: "Wait and retry after the Retry-After interval",
	}
}
