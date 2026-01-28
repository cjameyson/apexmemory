package app

import (
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	SessionDuration = 30 * 24 * time.Hour // 30 days
)

// Username validation rules
var (
	usernameRegex  = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	usernameMinLen = 3
	usernameMaxLen = 30
)

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &input)
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.Email == "" || input.Username == "" || input.Password == "" {
		app.RespondError(w, r, http.StatusBadRequest, "Email, username, and password are required")
		return
	}

	// Normalize and validate email format
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	addr, err := mail.ParseAddress(input.Email)
	if err != nil || addr.Address != input.Email {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate username
	input.Username = strings.TrimSpace(input.Username)
	if len(input.Username) < usernameMinLen {
		app.RespondError(w, r, http.StatusBadRequest, "Username must be at least 3 characters")
		return
	}
	if len(input.Username) > usernameMaxLen {
		app.RespondError(w, r, http.StatusBadRequest, "Username must not exceed 30 characters")
		return
	}
	if !usernameRegex.MatchString(input.Username) {
		app.RespondError(w, r, http.StatusBadRequest, "Username must start with a letter and contain only letters, numbers, underscores, and hyphens")
		return
	}

	// Password length validation (min 8, max 128 to prevent DoS via expensive hashing)
	if len(input.Password) < 8 {
		app.RespondError(w, r, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}
	if len(input.Password) > 128 {
		app.RespondError(w, r, http.StatusBadRequest, "Password must not exceed 128 characters")
		return
	}

	user, err := app.RegisterUser(r.Context(), input.Email, input.Username, input.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_email_key":
				app.RespondErrorWithCode(w, r, http.StatusConflict, "EMAIL_EXISTS", "Email already registered")
				return
			case "users_username_key":
				app.RespondErrorWithCode(w, r, http.StatusConflict, "USERNAME_EXISTS", "Username already taken")
				return
			}
		}
		app.RespondServerError(w, r, ErrDBQuery("insert user", err))
		return
	}

	userAgent := r.Header.Get("User-Agent")
	ipAddress := app.GetClientIP(r)

	sessionToken, err := app.CreateSession(r.Context(), user.ID, userAgent, ipAddress, SessionDuration)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("create session", err))
		return
	}

	app.RespondJSON(w, r, http.StatusCreated, sessionResponse(user, sessionToken))
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &input)
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.Email == "" || input.Password == "" {
		app.RespondError(w, r, http.StatusBadRequest, "Email and password are required")
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	user, err := app.AuthenticateUser(r.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			app.RespondErrorWithCode(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("authenticate", err))
		return
	}

	userAgent := r.Header.Get("User-Agent")
	ipAddress := app.GetClientIP(r)

	sessionToken, err := app.CreateSession(r.Context(), user.ID, userAgent, ipAddress, SessionDuration)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("create session", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, sessionResponse(user, sessionToken))
}

func sessionResponse(user *AppUser, token string) map[string]interface{} {
	return map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
		"session_token": token,
		"expires_at":    time.Now().Add(SessionDuration),
	}
}

// LogoutHandler invalidates the current session.
// Note: The IsAnonymous check is defense-in-depth; RequireAuth middleware already
// rejects anonymous users, but this guard protects against misconfigured routes.
func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	token := app.extractTokenFromHeader(r)
	if err := app.DeleteSession(r.Context(), token); err != nil {
		GetLogger(r.Context()).Error("failed to delete session", "error", ErrDBQuery("delete session", err))
	}

	response := map[string]string{
		"message": "Logged out successfully",
	}

	app.RespondJSON(w, r, http.StatusOK, response)
}

// LogoutAllHandler invalidates all sessions for the current user.
// Note: The IsAnonymous check is defense-in-depth; RequireAuth middleware already
// rejects anonymous users, but this guard protects against misconfigured routes.
func (app *Application) LogoutAllHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	if err := app.DeleteUserSessions(r.Context(), user.ID); err != nil {
		app.RespondServerError(w, r, ErrDBQuery("delete sessions", err))
		return
	}

	response := map[string]string{
		"message": "Logged out from all devices successfully",
	}

	app.RespondJSON(w, r, http.StatusOK, response)
}
