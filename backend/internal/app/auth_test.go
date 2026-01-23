package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	app := testApp(t)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
		wantError  bool
	}{
		{
			name: "valid registration",
			body: map[string]string{
				"email":    "test-register@example.com",
				"username": "testuser123",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "missing email",
			body: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "short password",
			body: map[string]string{
				"email":    "short-password@example.com",
				"username": "shortpw",
				"password": "short",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "short username",
			body: map[string]string{
				"email":    "short-username@example.com",
				"username": "ab",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "invalid email format",
			body: map[string]string{
				"email":    "not-an-email",
				"username": "invalidemail",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "invalid username format",
			body: map[string]string{
				"email":    "invalid-username@example.com",
				"username": "123start",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(t, "POST", "/v1/auth/register", tt.body)
			rr := httptest.NewRecorder()

			app.RegisterHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}

			var resp map[string]any
			decodeResponse(t, rr, &resp)

			if tt.wantError {
				if _, ok := resp["error"]; !ok {
					t.Errorf("expected error in response, got: %v", resp)
				}
			} else {
				if _, ok := resp["session_token"]; !ok {
					t.Errorf("expected session_token in response, got: %v", resp)
				}
			}
		})
	}
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	app := testApp(t)

	email := "duplicate@example.com"

	// First registration should succeed
	req := jsonRequest(t, "POST", "/v1/auth/register", map[string]string{
		"email":    email,
		"username": "firstuser",
		"password": "password123",
	})
	rr := httptest.NewRecorder()
	app.RegisterHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("first registration failed: %s", rr.Body.String())
	}

	// Second registration with same email should fail
	req = jsonRequest(t, "POST", "/v1/auth/register", map[string]string{
		"email":    email,
		"username": "seconduser",
		"password": "password123",
	})
	rr = httptest.NewRecorder()
	app.RegisterHandler(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected status %d for duplicate email, got %d. Body: %s",
			http.StatusConflict, rr.Code, rr.Body.String())
	}
}

func TestLoginHandler(t *testing.T) {
	app := testApp(t)

	testEmail := "test-login@example.com"

	// First register a user
	_, err := app.RegisterUser(t.Context(), testEmail, "loginuser", "password123")
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
		wantError  bool
	}{
		{
			name: "valid login",
			body: map[string]string{
				"email":    testEmail,
				"password": "password123",
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "wrong password",
			body: map[string]string{
				"email":    testEmail,
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  true,
		},
		{
			name: "nonexistent user",
			body: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  true,
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": testEmail,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(t, "POST", "/v1/auth/login", tt.body)
			rr := httptest.NewRecorder()

			app.LoginHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}

			var resp map[string]any
			decodeResponse(t, rr, &resp)

			if tt.wantError {
				if _, ok := resp["error"]; !ok {
					t.Errorf("expected error in response, got: %v", resp)
				}
			} else {
				if _, ok := resp["session_token"]; !ok {
					t.Errorf("expected session_token in response, got: %v", resp)
				}
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	app := testApp(t)

	testEmail := "test-logout@example.com"

	// Register and login
	user, err := app.RegisterUser(t.Context(), testEmail, "logoutuser", "password123")
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	token, err := app.CreateSession(t.Context(), user.ID, "test-agent", "127.0.0.1", SessionDuration)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	t.Run("logout with valid token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.LogoutHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}
	})

	t.Run("logout without auth", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/auth/logout", nil)
		req = app.WithUser(req, AnonymousUser)
		rr := httptest.NewRecorder()

		app.LogoutHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})
}

func TestGetCurrentUserHandler(t *testing.T) {
	app := testApp(t)

	testEmail := "test-me@example.com"

	// Register user
	user, err := app.RegisterUser(t.Context(), testEmail, "meuser", "password123")
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	t.Run("get current user authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/auth/me", nil)
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetCurrentUserHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["email"] != testEmail {
			t.Errorf("expected email %s, got %v", testEmail, resp["email"])
		}
	})

	t.Run("get current user unauthenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/auth/me", nil)
		req = app.WithUser(req, AnonymousUser)
		rr := httptest.NewRecorder()

		app.GetCurrentUserHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})
}

func TestPasswordHashing(t *testing.T) {
	app := testApp(t)

	password := "test-password-123"

	hash, err := app.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	t.Run("verify correct password", func(t *testing.T) {
		err := app.VerifyPassword(password, hash)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("verify wrong password", func(t *testing.T) {
		err := app.VerifyPassword("wrong-password", hash)
		if err != ErrInvalidPassword {
			t.Errorf("expected ErrInvalidPassword, got: %v", err)
		}
	})
}
