package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"apexmemory.ai/internal/db"
	"apexmemory.ai/internal/storage"
	"github.com/jackc/pgx/v5"
)

// AssetResponse is the API response for an asset.
type AssetResponse struct {
	ID          string          `json:"id"`
	ContentType string          `json:"content_type"`
	Filename    string          `json:"filename"`
	SizeBytes   int64           `json:"size_bytes"`
	Sha256      string          `json:"sha256"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
}

// UploadAssetHandler handles POST /v1/assets
func (app *Application) UploadAssetHandler(w http.ResponseWriter, r *http.Request) {
	user := app.MustUser(r)

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+1024) // extra for multipart overhead

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid multipart form or file too large")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Missing or invalid file field")
		return
	}
	defer file.Close()

	// Determine content type from the multipart file header
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	asset, err := app.UploadAsset(r.Context(), user.ID, header.Filename, contentType, file)
	if err != nil {
		var validationErr *AssetValidationError
		if errors.As(err, &validationErr) {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrInternal("asset upload failed", err))
		return
	}

	resp := AssetResponse{
		ID:          asset.ID.String(),
		ContentType: asset.ContentType,
		Filename:    asset.Filename,
		SizeBytes:   asset.SizeBytes,
		Sha256:      asset.Sha256,
		Metadata:    json.RawMessage(asset.Metadata),
		CreatedAt:   asset.CreatedAt,
	}

	app.RespondJSON(w, r, http.StatusCreated, resp)
}

// ServeAssetHandler handles GET /v1/assets/{id}/file
func (app *Application) ServeAssetHandler(w http.ResponseWriter, r *http.Request) {
	user := app.MustUser(r)

	assetID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	// Verify ownership via DB lookup
	asset, err := app.Queries.GetAsset(r.Context(), db.GetAssetParams{
		UserID: user.ID,
		ID:     assetID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.RespondError(w, r, http.StatusNotFound, "Asset not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("get asset", err))
		return
	}

	// ETag-based caching: if client already has this version, return 304
	etag := `"` + asset.Sha256 + `"`
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Fetch file from storage
	key := storageKey(user.ID, assetID)
	reader, _, err := app.Storage.Get(r.Context(), key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Asset file not found")
			return
		}
		app.RespondServerError(w, r, ErrInternal("failed to retrieve asset file", err))
		return
	}
	defer reader.Close()

	// Set response headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Content-Type", asset.ContentType)
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	w.Header().Set("ETag", etag)

	// Stream file to response
	if _, err := io.Copy(w, reader); err != nil {
		// Client likely disconnected; log but don't try to write error response
		app.logError(r, ErrInternal("failed to stream asset", err))
	}
}
