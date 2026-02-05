//go:build integration

package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"apexmemory.ai/internal/storage"
	"github.com/google/uuid"
)

// testAppWithStorage creates a test app with local file storage backed by a temp directory.
func testAppWithStorage(t *testing.T) *Application {
	t.Helper()
	app := testApp(t)
	store, err := storage.NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocalStorage: %v", err)
	}
	app.Storage = store
	return app
}

// createTestPNG generates a solid red PNG image of the given dimensions.
func createTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, red)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode PNG: %v", err)
	}
	return buf.Bytes()
}

// createTestJPEG generates a solid blue JPEG image of the given dimensions.
func createTestJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, blue)
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("failed to encode JPEG: %v", err)
	}
	return buf.Bytes()
}

// multipartRequest builds an HTTP request with a multipart/form-data body containing
// one file part with the given field name, filename, content type, and data.
func multipartRequest(t *testing.T, fieldName, filename, contentType string, data []byte) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Use CreatePart with explicit Content-Type header so the server sees the correct MIME type.
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("failed to create form part: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(data)); err != nil {
		t.Fatalf("failed to write file data: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/v1/assets", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// ---------------------------------------------------------------------------
// TestUploadAssetHandler
// ---------------------------------------------------------------------------

func TestUploadAssetHandler(t *testing.T) {
	app := testAppWithStorage(t)
	user := createTestUser(t, app)

	pngData := createTestPNG(t, 100, 80)
	jpegData := createTestJPEG(t, 64, 64)

	tests := []struct {
		name        string
		fieldName   string
		filename    string
		contentType string
		data        []byte
		wantStatus  int
		wantError   bool
	}{
		{
			name:        "valid PNG upload",
			fieldName:   "file",
			filename:    "test.png",
			contentType: "image/png",
			data:        pngData,
			wantStatus:  http.StatusCreated,
			wantError:   false,
		},
		{
			name:        "valid JPEG upload",
			fieldName:   "file",
			filename:    "photo.jpg",
			contentType: "image/jpeg",
			data:        jpegData,
			wantStatus:  http.StatusCreated,
			wantError:   false,
		},
		{
			name:        "invalid content type",
			fieldName:   "file",
			filename:    "readme.txt",
			contentType: "text/plain",
			data:        []byte("hello world"),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "empty file",
			fieldName:   "file",
			filename:    "empty.png",
			contentType: "image/png",
			data:        []byte{},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := multipartRequest(t, tt.fieldName, tt.filename, tt.contentType, tt.data)
			req = app.WithUser(req, user)
			rr := httptest.NewRecorder()

			app.UploadAssetHandler(rr, req)

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
				// Verify response fields
				if _, ok := resp["id"]; !ok {
					t.Errorf("expected id in response, got: %v", resp)
				}
				if resp["content_type"] != tt.contentType {
					t.Errorf("expected content_type %q, got %v", tt.contentType, resp["content_type"])
				}
				if resp["filename"] != tt.filename {
					t.Errorf("expected filename %q, got %v", tt.filename, resp["filename"])
				}

				// SHA-256 must match the uploaded data
				hash := sha256.Sum256(tt.data)
				expectedHash := hex.EncodeToString(hash[:])
				if resp["sha256"] != expectedHash {
					t.Errorf("expected sha256 %q, got %v", expectedHash, resp["sha256"])
				}

				// size_bytes must match
				sizeFloat, ok := resp["size_bytes"].(float64)
				if !ok {
					t.Fatalf("expected size_bytes to be a number, got %T", resp["size_bytes"])
				}
				if int(sizeFloat) != len(tt.data) {
					t.Errorf("expected size_bytes %d, got %d", len(tt.data), int(sizeFloat))
				}
			}
		})
	}

	t.Run("missing file field", func(t *testing.T) {
		// Build a multipart body with the wrong field name so "file" is absent.
		req := multipartRequest(t, "wrong_field", "test.png", "image/png", pngData)
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.UploadAssetHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})
}

// ---------------------------------------------------------------------------
// TestServeAssetHandler
// ---------------------------------------------------------------------------

func TestServeAssetHandler(t *testing.T) {
	app := testAppWithStorage(t)
	user := createTestUser(t, app)

	// Upload an asset via the business-logic method so we have a known file to serve.
	pngData := createTestPNG(t, 50, 50)
	asset, err := app.UploadAsset(t.Context(), user.ID, "serve-test.png", "image/png", bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("UploadAsset: %v", err)
	}

	expectedETag := `"` + asset.Sha256 + `"`

	t.Run("serves file with correct headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/assets/"+asset.ID.String()+"/file", nil)
		req.SetPathValue("id", asset.ID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.ServeAssetHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		// Content-Type
		if got := rr.Header().Get("Content-Type"); got != "image/png" {
			t.Errorf("expected Content-Type %q, got %q", "image/png", got)
		}

		// Cache-Control
		wantCC := "private, max-age=31536000, immutable"
		if got := rr.Header().Get("Cache-Control"); got != wantCC {
			t.Errorf("expected Cache-Control %q, got %q", wantCC, got)
		}

		// ETag
		if got := rr.Header().Get("ETag"); got != expectedETag {
			t.Errorf("expected ETag %q, got %q", expectedETag, got)
		}

		// Body must match the original data
		if !bytes.Equal(rr.Body.Bytes(), pngData) {
			t.Errorf("response body does not match uploaded data (got %d bytes, want %d)", rr.Body.Len(), len(pngData))
		}
	})

	t.Run("304 with matching ETag", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/assets/"+asset.ID.String()+"/file", nil)
		req.SetPathValue("id", asset.ID.String())
		req.Header.Set("If-None-Match", expectedETag)
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.ServeAssetHandler(rr, req)

		if rr.Code != http.StatusNotModified {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNotModified, rr.Code, rr.Body.String())
		}

		// Body should be empty on 304
		if rr.Body.Len() != 0 {
			t.Errorf("expected empty body on 304, got %d bytes", rr.Body.Len())
		}
	})

	t.Run("404 for unknown asset ID", func(t *testing.T) {
		fakeID := uuid.New()
		req := httptest.NewRequest("GET", "/v1/assets/"+fakeID.String()+"/file", nil)
		req.SetPathValue("id", fakeID.String())
		req = app.WithUser(req, user)
		rr := httptest.NewRecorder()

		app.ServeAssetHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
		}
	})
}
