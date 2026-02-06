package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	defaultPageLimit = 50
	maxPageLimit     = 100
)

// parsePagination extracts limit and offset from query params with defaults.
func parsePagination(r *http.Request) (limit, offset int32) {
	limit = defaultPageLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = int32(n)
		}
	}
	if limit > maxPageLimit {
		limit = maxPageLimit
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = int32(n)
		}
	}
	return
}

// PageResponse is a typed wrapper for paginated list responses.
type PageResponse[T any] struct {
	Data    []T   `json:"data"`
	Total   int64 `json:"total"`
	HasMore bool  `json:"has_more"`
}

// NewPageResponse constructs a PageResponse with has_more calculated automatically.
func NewPageResponse[T any](data []T, total int64, limit, offset int32) PageResponse[T] {
	return PageResponse[T]{
		Data:    data,
		Total:   total,
		HasMore: int64(offset)+int64(len(data)) < total,
	}
}

// maxJSONBodySize is the maximum allowed size for JSON request bodies (1MB).
// Prevents memory exhaustion from large payloads.
const maxJSONBodySize = 1 * 1024 * 1024

func (app *Application) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	return nil
}

func (app *Application) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	var (
		js  []byte
		err error
	)

	if app.Config.Env == "production" {
		js, err = json.Marshal(data)
	} else {
		js, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *Application) PathUUID(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool) {
	idStr := r.PathValue(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "invalid ID format")
		return uuid.Nil, false
	}
	return id, true
}

// MustUser returns the authenticated user from the request context.
// Panics if called for an anonymous user - this should only be used in
// handlers behind RequireAuth middleware.
func (app *Application) MustUser(r *http.Request) *AppUser {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		panic("MustUser called for anonymous user - check route middleware")
	}
	return user
}

// SortParam holds a parsed sort field and direction.
type SortParam struct {
	Field string // empty = use query default
	Asc   bool
}

// parseSort parses `?sort=field` or `?sort=-field` from the request.
// Returns error if the field is not in allowedFields.
// Empty/missing sort param returns zero SortParam (no error).
func parseSort(r *http.Request, allowedFields ...string) (SortParam, error) {
	raw := r.URL.Query().Get("sort")
	if raw == "" {
		return SortParam{}, nil
	}

	var sp SortParam
	if raw[0] == '-' {
		sp.Field = raw[1:]
	} else {
		sp.Field = raw
		sp.Asc = true
	}

	for _, f := range allowedFields {
		if sp.Field == f {
			return sp, nil
		}
	}
	return SortParam{}, fmt.Errorf("invalid sort field: %q", sp.Field)
}

// Optional types for distinguishing "not provided" from "explicitly null" in JSON

// OptionalString distinguishes between missing, null, and present string values.
type OptionalString struct {
	Value *string
	Set   bool
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.Value = &s
	return nil
}

// OptionalUUID distinguishes between missing, null, and present UUID values.
type OptionalUUID struct {
	Value *uuid.UUID
	Set   bool
}

func (o *OptionalUUID) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var id uuid.UUID
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	o.Value = &id
	return nil
}
