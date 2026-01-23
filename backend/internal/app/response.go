package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// RespondJSON writes a JSON response with the given status code.
// No wrapper envelope - data is written directly.
func (app *Application) RespondJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	err := app.WriteJSON(w, status, data, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RespondError writes an error response with the given status code.
// Response format: {"error": "message", "code": "ERROR_CODE"}
func (app *Application) RespondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	response := map[string]string{
		"error": message,
	}

	err := app.WriteJSON(w, status, response, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RespondErrorWithCode writes an error response with a specific error code.
func (app *Application) RespondErrorWithCode(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	response := map[string]string{
		"error": message,
		"code":  code,
	}

	err := app.WriteJSON(w, status, response, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RespondFieldErrors writes a validation error response with field-specific errors.
func (app *Application) RespondFieldErrors(w http.ResponseWriter, r *http.Request, fieldErrors map[string]string) {
	response := map[string]any{
		"error":       "Validation failed",
		"code":        "VALIDATION_FAILED",
		"fieldErrors": fieldErrors,
	}

	err := app.WriteJSON(w, http.StatusBadRequest, response, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RespondServerError logs the error and writes a generic server error response.
func (app *Application) RespondServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.RespondError(w, r, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

// logError logs an error with full context from the request.
func (app *Application) logError(r *http.Request, err error) {
	GetLogger(r.Context()).Error("operation_failed",
		slog.String("event", "operation_failed"),
		errorAttr(err),
	)
}

func errorAttr(err error) slog.Attr {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return slog.Attr{Key: "error", Value: appErr.LogValue()}
	}

	return slog.Attr{
		Key: "error",
		Value: slog.GroupValue(
			slog.String("error_code", "UNKNOWN"),
			slog.String("error_message", err.Error()),
			slog.String("error_type", fmt.Sprintf("%T", err)),
			slog.String("severity", string(SeverityMedium)),
			slog.Bool("retryable", false),
		),
	}
}
