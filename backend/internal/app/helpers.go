package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// maxJSONBodySize is the maximum allowed size for JSON request bodies (1MB).
// Prevents memory exhaustion from large payloads.
const maxJSONBodySize = 1 * 1024 * 1024

func (app *Application) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
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

	if dec.More() {
		return errors.New("body must contain a single JSON value")
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
