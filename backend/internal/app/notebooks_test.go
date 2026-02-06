//go:build integration

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestCreateNotebookHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
		wantError  bool
	}{
		{
			name: "success minimal",
			body: map[string]any{
				"name": "My Notebook",
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "success full",
			body: map[string]any{
				"name":        "Full Notebook",
				"description": "A detailed description",
				"emoji":       "📚",
				"color":       "#FF5733",
				"position":    5,
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "success with description",
			body: map[string]any{
				"name":        "With Description",
				"description": "This is a test description",
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "fail empty name",
			body: map[string]any{
				"name": "",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "fail whitespace name",
			body: map[string]any{
				"name": "   ",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "fail missing name",
			body:       map[string]any{},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "fail description too long",
			body: map[string]any{
				"name":        "Valid Name",
				"description": string(make([]byte, 10001)),
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(t, "POST", "/v1/notebooks", tt.body)
			req = app.WithUser(req, user)
			rr := httptest.NewRecorder()

			app.CreateNotebookHandler(rr, req)

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
				if _, ok := resp["id"]; !ok {
					t.Errorf("expected id in response, got: %v", resp)
				}
				if resp["name"] != tt.body["name"] {
					t.Errorf("expected name %v, got %v", tt.body["name"], resp["name"])
				}
				// Verify description if provided
				if desc, ok := tt.body["description"]; ok {
					if resp["description"] != desc {
						t.Errorf("expected description %v, got %v", desc, resp["description"])
					}
				}
				if emoji, ok := tt.body["emoji"]; ok {
					if resp["emoji"] != emoji {
						t.Errorf("expected emoji %v, got %v", emoji, resp["emoji"])
					}
				}
				if color, ok := tt.body["color"]; ok {
					if resp["color"] != color {
						t.Errorf("expected color %v, got %v", color, resp["color"])
					}
				}
			}
		})
	}
}

func TestCreateNotebookHandler_Unauthenticated(t *testing.T) {
	app := testApp(t)

	handler := app.RequireAuth(http.HandlerFunc(app.CreateNotebookHandler))
	req := jsonRequest(t, "POST", "/v1/notebooks", map[string]any{
		"name": "My Notebook",
	})
	req = app.WithUser(req, AnonymousUser)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestGetNotebookHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create a notebook with description
	desc := "Test description"
	notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
		Name:        "Test Notebook",
		Description: &desc,
	})
	if err != nil {
		t.Fatalf("failed to create test notebook: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/notebooks/"+notebook.ID.String(), nil)
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["name"] != "Test Notebook" {
			t.Errorf("expected name 'Test Notebook', got %v", resp["name"])
		}
		if resp["description"] != "Test description" {
			t.Errorf("expected description 'Test description', got %v", resp["description"])
		}
	})

	t.Run("404 non-existent", func(t *testing.T) {
		fakeID := uuid.New()
		req := httptest.NewRequest("GET", "/v1/notebooks/"+fakeID.String(), nil)
		req.SetPathValue("id", fakeID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetNotebookHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
		}
	})

	t.Run("404 other user's notebook", func(t *testing.T) {
		otherUser := createTestUserWithEmail(t, app, "other@example.com", "otheruser")

		req := httptest.NewRequest("GET", "/v1/notebooks/"+notebook.ID.String(), nil)
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, otherUser)
		rr := httptest.NewRecorder()

		app.GetNotebookHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
		}
	})
}

func TestListNotebooksHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	t.Run("empty array", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/notebooks", nil)
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.ListNotebooksHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp []map[string]any
		decodeResponse(t, rr, &resp)

		if len(resp) != 0 {
			t.Errorf("expected empty array, got %d items", len(resp))
		}
	})

	t.Run("correct order and excludes archived", func(t *testing.T) {
		// Create notebooks with different positions
		pos1 := int32(1)
		pos2 := int32(2)
		pos3 := int32(3)

		_, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name:     "Notebook C",
			Position: &pos3,
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}
		_, err = app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name:     "Notebook A",
			Position: &pos1,
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}
		notebookB, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name:     "Notebook B",
			Position: &pos2,
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		// Archive one notebook
		err = app.ArchiveNotebook(t.Context(), user.ID, notebookB.ID)
		if err != nil {
			t.Fatalf("failed to archive notebook: %v", err)
		}

		req := httptest.NewRequest("GET", "/v1/notebooks", nil)
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.ListNotebooksHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp []map[string]any
		decodeResponse(t, rr, &resp)

		if len(resp) != 2 {
			t.Fatalf("expected 2 notebooks (excluding archived), got %d", len(resp))
		}

		// Should be ordered by position ASC
		if resp[0]["name"] != "Notebook A" {
			t.Errorf("expected first notebook to be 'Notebook A', got %v", resp[0]["name"])
		}
		if resp[1]["name"] != "Notebook C" {
			t.Errorf("expected second notebook to be 'Notebook C', got %v", resp[1]["name"])
		}
	})
}

func TestUpdateNotebookHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	t.Run("success name", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "Original Name",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"name": "Updated Name",
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["name"] != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %v", resp["name"])
		}
	})

	t.Run("success update description", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "Desc Test",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"description": "New description",
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["description"] != "New description" {
			t.Errorf("expected description 'New description', got %v", resp["description"])
		}
	})

	t.Run("success clear description", func(t *testing.T) {
		desc := "Original description"
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name:        "Clear Desc Test",
			Description: &desc,
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		// Clear description by setting to null
		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"description": nil,
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["description"] != nil {
			t.Errorf("expected description to be nil, got %v", resp["description"])
		}
	})

	t.Run("success desired_retention", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "Retention Test",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"desired_retention": 0.85,
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		fsrs, ok := resp["fsrs_settings"].(map[string]any)
		if !ok {
			t.Fatalf("expected fsrs_settings object, got %v", resp["fsrs_settings"])
		}
		if fsrs["desired_retention"] != 0.85 {
			t.Errorf("expected desired_retention 0.85, got %v", fsrs["desired_retention"])
		}
	})

	t.Run("success set and clear emoji and color", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "Emoji Color Test",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		// Set emoji and color
		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"emoji": "📚",
			"color": "#FF5733",
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["emoji"] != "📚" {
			t.Errorf("expected emoji 📚, got %v", resp["emoji"])
		}
		if resp["color"] != "#FF5733" {
			t.Errorf("expected color #FF5733, got %v", resp["color"])
		}

		// Clear emoji by setting to null
		req = jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"emoji": nil,
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr = httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		decodeResponse(t, rr, &resp)

		if resp["emoji"] != nil {
			t.Errorf("expected emoji to be nil, got %v", resp["emoji"])
		}
		if resp["color"] != "#FF5733" {
			t.Errorf("expected color preserved as #FF5733, got %v", resp["color"])
		}
	})

	t.Run("partial update preserves fields", func(t *testing.T) {
		desc := "Original description"
		pos := int32(5)
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name:        "Original",
			Description: &desc,
			Position:    &pos,
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		// Update only name
		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"name": "New Name",
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["name"] != "New Name" {
			t.Errorf("expected name 'New Name', got %v", resp["name"])
		}
		if resp["description"] != "Original description" {
			t.Errorf("expected description preserved, got %v", resp["description"])
		}
		if int32(resp["position"].(float64)) != 5 {
			t.Errorf("expected position preserved as 5, got %v", resp["position"])
		}
	})

	t.Run("fail retention equals 0", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "Zero Retention",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"desired_retention": 0.0,
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})

	t.Run("fail retention equals 1", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "One Retention",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"desired_retention": 1.0,
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})

	t.Run("fail empty name", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "Empty Name Test",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebook.ID.String(), map[string]any{
			"name": "",
		})
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNotebookHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})
}

func TestDeleteNotebookHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	t.Run("success", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "To Delete",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		req := httptest.NewRequest("DELETE", "/v1/notebooks/"+notebook.ID.String(), nil)
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.DeleteNotebookHandler(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNoContent, rr.Code, rr.Body.String())
		}

		// Verify excluded from list
		listReq := httptest.NewRequest("GET", "/v1/notebooks", nil)
		listReq = app.WithUser(listReq, user)
		listRR := httptest.NewRecorder()

		app.ListNotebooksHandler(listRR, listReq)

		var notebooks []map[string]any
		decodeResponse(t, listRR, &notebooks)

		for _, n := range notebooks {
			if n["id"] == notebook.ID.String() {
				t.Errorf("deleted notebook should not appear in list")
			}
		}
	})

	t.Run("404 non-existent", func(t *testing.T) {
		fakeID := uuid.New()
		req := httptest.NewRequest("DELETE", "/v1/notebooks/"+fakeID.String(), nil)
		req.SetPathValue("id", fakeID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.DeleteNotebookHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
		}
	})

	t.Run("idempotent archive already archived", func(t *testing.T) {
		notebook, err := app.CreateNotebook(t.Context(), user.ID, CreateNotebookParams{
			Name: "To Archive Twice",
		})
		if err != nil {
			t.Fatalf("failed to create notebook: %v", err)
		}

		// First archive
		req := httptest.NewRequest("DELETE", "/v1/notebooks/"+notebook.ID.String(), nil)
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.DeleteNotebookHandler(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("first archive: expected status %d, got %d", http.StatusNoContent, rr.Code)
		}

		// Second archive (idempotent)
		req = httptest.NewRequest("DELETE", "/v1/notebooks/"+notebook.ID.String(), nil)
		req.SetPathValue("id", notebook.ID.String())
		req = app.WithUser(req, user)
		rr = httptest.NewRecorder()

		app.DeleteNotebookHandler(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("second archive (idempotent): expected status %d, got %d. Body: %s", http.StatusNoContent, rr.Code, rr.Body.String())
		}
	})
}

// Helper functions for tests

func createTestUser(t *testing.T, app *Application) *AppUser {
	t.Helper()
	return createTestUserWithEmail(t, app, "test@example.com", "testuser")
}

func createTestUserWithEmail(t *testing.T, app *Application, email, username string) *AppUser {
	t.Helper()

	user, err := app.RegisterUser(t.Context(), email, username, "password123")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}
