package app

import (
	"net/http"

	"github.com/google/uuid"
)

type AppUser struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

// AnonymousUser represents an unauthenticated user.
// Use IsAnonymous() to check - never compare pointers directly.
var AnonymousUser = &AppUser{}

func (user *AppUser) IsAnonymous() bool {
	return user.ID == uuid.Nil
}

func (app *Application) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	app.RespondJSON(w, r, http.StatusOK, map[string]any{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}
