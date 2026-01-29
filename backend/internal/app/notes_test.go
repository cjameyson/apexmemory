package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// createTestNotebook creates a notebook for use in note tests.
func createTestNotebook(t *testing.T, app *Application, userID uuid.UUID) uuid.UUID {
	t.Helper()
	nb, err := app.CreateNotebook(t.Context(), userID, CreateNotebookParams{
		Name: "Test Notebook",
	})
	if err != nil {
		t.Fatalf("failed to create test notebook: %v", err)
	}
	return nb.ID
}

// basicContent returns valid basic note content.
func basicContent() json.RawMessage {
	return json.RawMessage(`{
		"version": 1,
		"fields": [{"name": "front", "type": "plain_text", "value": "Q"}, {"name": "back", "type": "plain_text", "value": "A"}]
	}`)
}

// clozeContent returns valid cloze note content.
func clozeContent(text string) json.RawMessage {
	return json.RawMessage(`{
		"version": 1,
		"fields": [{"name": "text", "type": "cloze_text", "cloze_text": "` + text + `"}]
	}`)
}

// imageOcclusionContent returns valid image occlusion content with given region IDs.
func imageOcclusionContent(regionIDs ...string) json.RawMessage {
	regions := ""
	for i, id := range regionIDs {
		if i > 0 {
			regions += ","
		}
		regions += `{"id": "` + id + `", "shape": {"type": "rect", "x": 0, "y": 0, "width": 100, "height": 100}}`
	}
	return json.RawMessage(`{
		"version": 1,
		"fields": [{
			"name": "image",
			"type": "image_occlusion",
			"image": {"url": "https://example.com/img.png", "width": 800, "height": 600},
			"regions": [` + regions + `]
		}]
	}`)
}

func TestCreateNoteHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
		wantCards  int
	}{
		{
			name: "basic note",
			body: map[string]any{
				"note_type": "basic",
				"content": map[string]any{
					"version": 1,
					"fields":  []any{map[string]any{"name": "front", "type": "plain_text", "value": "Q"}},
				},
			},
			wantStatus: http.StatusCreated,
			wantCards:  1,
		},
		{
			name: "cloze note with 2 deletions",
			body: map[string]any{
				"note_type": "cloze",
				"content": map[string]any{
					"version": 1,
					"fields":  []any{map[string]any{"name": "text", "type": "cloze_text", "cloze_text": "The {{c1::mitochondria}} is the {{c2::powerhouse}}"}},
				},
			},
			wantStatus: http.StatusCreated,
			wantCards:  2,
		},
		{
			name: "image occlusion with 3 regions",
			body: map[string]any{
				"note_type": "image_occlusion",
				"content": map[string]any{
					"version": 1,
					"fields": []any{map[string]any{
						"name": "image",
						"type": "image_occlusion",
						"image": map[string]any{
							"url": "https://example.com/img.png", "width": 800, "height": 600,
						},
						"regions": []any{
							map[string]any{"id": "m_abcdef12", "shape": map[string]any{"type": "rect", "x": 0, "y": 0, "width": 100, "height": 100}},
							map[string]any{"id": "m_ghijkl34", "shape": map[string]any{"type": "rect", "x": 100, "y": 0, "width": 100, "height": 100}},
							map[string]any{"id": "m_mnopqr56", "shape": map[string]any{"type": "rect", "x": 200, "y": 0, "width": 100, "height": 100}},
						},
					}},
				},
			},
			wantStatus: http.StatusCreated,
			wantCards:  3,
		},
		{
			name: "default note type is basic",
			body: map[string]any{
				"content": map[string]any{
					"version": 1,
					"fields":  []any{map[string]any{"name": "front", "type": "plain_text", "value": "Q"}},
				},
			},
			wantStatus: http.StatusCreated,
			wantCards:  1,
		},
		{
			name: "fail missing content",
			body: map[string]any{
				"note_type": "basic",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "fail invalid content no version",
			body: map[string]any{
				"note_type": "basic",
				"content": map[string]any{
					"fields": []any{},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "fail cloze with no deletions",
			body: map[string]any{
				"note_type": "cloze",
				"content": map[string]any{
					"version": 1,
					"fields":  []any{map[string]any{"name": "text", "type": "cloze_text", "cloze_text": "no cloze here"}},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "fail image occlusion with no regions",
			body: map[string]any{
				"note_type": "image_occlusion",
				"content": map[string]any{
					"version": 1,
					"fields": []any{map[string]any{
						"name":    "image",
						"type":    "image_occlusion",
						"regions": []any{},
					}},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "fail image occlusion duplicate region ID",
			body: map[string]any{
				"note_type": "image_occlusion",
				"content": map[string]any{
					"version": 1,
					"fields": []any{map[string]any{
						"name": "image",
						"type": "image_occlusion",
						"image": map[string]any{
							"url": "https://example.com/img.png", "width": 800, "height": 600,
						},
						"regions": []any{
							map[string]any{"id": "m_abcdef12", "shape": map[string]any{"type": "rect", "x": 0, "y": 0, "width": 100, "height": 100}},
							map[string]any{"id": "m_abcdef12", "shape": map[string]any{"type": "rect", "x": 100, "y": 0, "width": 100, "height": 100}},
						},
					}},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(t, "POST", "/v1/notebooks/"+notebookID.String()+"/notes", tt.body)
			req.SetPathValue("notebook_id", notebookID.String())
			req = app.WithUser(req, user)
			rr := httptest.NewRecorder()

			app.CreateNoteHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}

			if tt.wantStatus == http.StatusCreated && tt.wantCards > 0 {
				var resp map[string]any
				decodeResponse(t, rr, &resp)
				cards, ok := resp["cards"].([]any)
				if !ok {
					t.Fatalf("expected cards array, got %v", resp["cards"])
				}
				if len(cards) != tt.wantCards {
					t.Errorf("expected %d cards, got %d", tt.wantCards, len(cards))
				}
			}
		})
	}
}

func TestCreateNoteHandler_NotebookNotFound(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	fakeID := uuid.New()

	req := jsonRequest(t, "POST", "/v1/notebooks/"+fakeID.String()+"/notes", map[string]any{
		"content": map[string]any{"version": 1, "fields": []any{}},
	})
	req.SetPathValue("notebook_id", fakeID.String())
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.CreateNoteHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestGetNoteHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	note, cards, err := app.CreateNote(t.Context(), user.ID, notebookID, "cloze",
		clozeContent("The {{c1::mitochondria}} is the {{c2::powerhouse}}"))
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	t.Run("success with cards", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/notes/"+note.ID.String(), nil)
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", note.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetNoteHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		respCards, ok := resp["cards"].([]any)
		if !ok {
			t.Fatalf("expected cards array, got %v", resp["cards"])
		}
		if len(respCards) != len(cards) {
			t.Errorf("expected %d cards, got %d", len(cards), len(respCards))
		}
	})

	t.Run("404 non-existent", func(t *testing.T) {
		fakeID := uuid.New()
		req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/notes/"+fakeID.String(), nil)
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", fakeID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetNoteHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})
}

func TestListNotesHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	// Create 3 notes
	for i := 0; i < 3; i++ {
		_, _, err := app.CreateNote(t.Context(), user.ID, notebookID, "basic", basicContent())
		if err != nil {
			t.Fatalf("failed to create note %d: %v", i, err)
		}
	}

	req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/notes", nil)
	req.SetPathValue("notebook_id", notebookID.String())
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.ListNotesHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)

	data, ok := resp["data"].([]any)
	if !ok {
		t.Fatalf("expected data array")
	}
	if len(data) != 3 {
		t.Errorf("expected 3 notes, got %d", len(data))
	}
	if resp["total"] != float64(3) {
		t.Errorf("expected total 3, got %v", resp["total"])
	}
}

func TestUpdateNoteHandler_ClozeDiff(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	// Create cloze note with c1 and c2
	note, _, err := app.CreateNote(t.Context(), user.ID, notebookID, "cloze",
		clozeContent("{{c1::a}} and {{c2::b}}"))
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	t.Run("add c3, remove c2", func(t *testing.T) {
		newContent := clozeContent("{{c1::a}} and {{c3::c}}")

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebookID.String()+"/notes/"+note.ID.String(), map[string]any{
			"content": json.RawMessage(newContent),
		})
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", note.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNoteHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["created"] != float64(1) {
			t.Errorf("expected 1 created, got %v", resp["created"])
		}
		if resp["deleted"] != float64(1) {
			t.Errorf("expected 1 deleted, got %v", resp["deleted"])
		}
		if resp["unchanged"] != float64(1) {
			t.Errorf("expected 1 unchanged, got %v", resp["unchanged"])
		}
	})

	t.Run("content only change, no card diff", func(t *testing.T) {
		// Same cloze numbers, different text
		newContent := clozeContent("{{c1::alpha}} and {{c3::gamma}}")

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebookID.String()+"/notes/"+note.ID.String(), map[string]any{
			"content": json.RawMessage(newContent),
		})
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", note.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNoteHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["created"] != float64(0) {
			t.Errorf("expected 0 created, got %v", resp["created"])
		}
		if resp["deleted"] != float64(0) {
			t.Errorf("expected 0 deleted, got %v", resp["deleted"])
		}
		if resp["unchanged"] != float64(2) {
			t.Errorf("expected 2 unchanged, got %v", resp["unchanged"])
		}
	})
}

func TestUpdateNoteHandler_ImageOcclusionDiff(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	// Create image occlusion note with 2 regions
	note, _, err := app.CreateNote(t.Context(), user.ID, notebookID, "image_occlusion",
		imageOcclusionContent("m_region_aaa", "m_region_bbb"))
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	t.Run("add region, remove region", func(t *testing.T) {
		// Keep m_region_aaa, remove m_region_bbb, add m_region_ccc
		newContent := imageOcclusionContent("m_region_aaa", "m_region_ccc")

		req := jsonRequest(t, "PATCH", "/v1/notebooks/"+notebookID.String()+"/notes/"+note.ID.String(), map[string]any{
			"content": json.RawMessage(newContent),
		})
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", note.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UpdateNoteHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)

		if resp["created"] != float64(1) {
			t.Errorf("expected 1 created, got %v", resp["created"])
		}
		if resp["deleted"] != float64(1) {
			t.Errorf("expected 1 deleted, got %v", resp["deleted"])
		}
		if resp["unchanged"] != float64(1) {
			t.Errorf("expected 1 unchanged, got %v", resp["unchanged"])
		}
	})
}

func TestDeleteNoteHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	note, _, err := app.CreateNote(t.Context(), user.ID, notebookID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/v1/notebooks/"+notebookID.String()+"/notes/"+note.ID.String(), nil)
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", note.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.DeleteNoteHandler(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("404 after delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/v1/notebooks/"+notebookID.String()+"/notes/"+note.ID.String(), nil)
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", note.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.DeleteNoteHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})
}

func TestListCardsHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	// Create a cloze note (2 cards) + basic note (1 card)
	_, _, err := app.CreateNote(t.Context(), user.ID, notebookID, "cloze",
		clozeContent("{{c1::a}} and {{c2::b}}"))
	if err != nil {
		t.Fatalf("failed to create cloze note: %v", err)
	}
	_, _, err = app.CreateNote(t.Context(), user.ID, notebookID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create basic note: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/cards", nil)
	req.SetPathValue("notebook_id", notebookID.String())
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.ListCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)

	data, ok := resp["data"].([]any)
	if !ok {
		t.Fatalf("expected data array")
	}
	if len(data) != 3 {
		t.Errorf("expected 3 cards, got %d", len(data))
	}
}

func TestListCardsHandler_StateFilter(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	_, _, err := app.CreateNote(t.Context(), user.ID, notebookID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	// All cards are 'new' initially
	req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/cards?state=new", nil)
	req.SetPathValue("notebook_id", notebookID.String())
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.ListCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)
	data := resp["data"].([]any)
	if len(data) != 1 {
		t.Errorf("expected 1 new card, got %d", len(data))
	}

	// No review cards
	req = httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/cards?state=review", nil)
	req.SetPathValue("notebook_id", notebookID.String())
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.ListCardsHandler(rr, req)

	decodeResponse(t, rr, &resp)
	data = resp["data"].([]any)
	if len(data) != 0 {
		t.Errorf("expected 0 review cards, got %d", len(data))
	}
}

func TestGetCardHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	notebookID := createTestNotebook(t, app, user.ID)

	_, cards, err := app.CreateNote(t.Context(), user.ID, notebookID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		cardID := cards[0].ID
		req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/cards/"+cardID.String(), nil)
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", cardID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetCardHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]any
		decodeResponse(t, rr, &resp)
		if resp["state"] != "new" {
			t.Errorf("expected state 'new', got %v", resp["state"])
		}
	})

	t.Run("404 non-existent", func(t *testing.T) {
		fakeID := uuid.New()
		req := httptest.NewRequest("GET", "/v1/notebooks/"+notebookID.String()+"/cards/"+fakeID.String(), nil)
		req.SetPathValue("notebook_id", notebookID.String())
		req.SetPathValue("id", fakeID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.GetCardHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})
}

func TestContentValidation(t *testing.T) {
	tests := []struct {
		name     string
		noteType string
		content  json.RawMessage
		wantErr  bool
	}{
		{
			name:     "basic valid",
			noteType: "basic",
			content:  basicContent(),
			wantErr:  false,
		},
		{
			name:     "cloze valid",
			noteType: "cloze",
			content:  clozeContent("{{c1::hello}} {{c2::world}}"),
			wantErr:  false,
		},
		{
			name:     "cloze gaps allowed",
			noteType: "cloze",
			content:  clozeContent("{{c1::a}} {{c3::b}}"),
			wantErr:  false,
		},
		{
			name:     "cloze duplicates collapsed",
			noteType: "cloze",
			content:  clozeContent("{{c1::a}} {{c1::b}}"),
			wantErr:  false,
		},
		{
			name:     "cloze c0 invalid",
			noteType: "cloze",
			content:  clozeContent("{{c0::nope}}"),
			wantErr:  true,
		},
		{
			name:     "image occlusion valid",
			noteType: "image_occlusion",
			content:  imageOcclusionContent("m_abcdef12"),
			wantErr:  false,
		},
		{
			name:     "unsupported type",
			noteType: "unknown",
			content:  basicContent(),
			wantErr:  true,
		},
		{
			name:     "missing version",
			noteType: "basic",
			content:  json.RawMessage(`{"fields": []}`),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateNoteContent(tt.noteType, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}
