package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
)

const maxUploadSize = 10 << 20 // 10MB

var allowedContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

// AssetValidationError represents a validation error for asset uploads.
type AssetValidationError struct {
	Message string
}

func (e *AssetValidationError) Error() string {
	return e.Message
}

// storageKey returns the object storage key for an asset.
func storageKey(userID, assetID uuid.UUID) string {
	return fmt.Sprintf("assets/%s/%s", userID, assetID)
}

// assetMetadata is the JSON structure stored in the metadata column.
type assetMetadata struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// UploadAsset validates, stores, and records an asset upload.
func (app *Application) UploadAsset(ctx context.Context, userID uuid.UUID, filename, contentType string, reader io.Reader) (db.Asset, error) {
	// Validate content type
	if !allowedContentTypes[contentType] {
		return db.Asset{}, &AssetValidationError{
			Message: fmt.Sprintf("unsupported content type: %s", contentType),
		}
	}

	// Read file into memory (limited to maxUploadSize)
	data, err := io.ReadAll(io.LimitReader(reader, maxUploadSize+1))
	if err != nil {
		return db.Asset{}, fmt.Errorf("failed to read upload: %w", err)
	}
	if len(data) > maxUploadSize {
		return db.Asset{}, &AssetValidationError{
			Message: fmt.Sprintf("file exceeds maximum size of %d bytes", maxUploadSize),
		}
	}
	if len(data) == 0 {
		return db.Asset{}, &AssetValidationError{Message: "file is empty"}
	}

	// Verify the declared content type matches actual file bytes.
	detected := http.DetectContentType(data)
	if !allowedContentTypes[detected] {
		return db.Asset{}, &AssetValidationError{
			Message: fmt.Sprintf("file content does not match an allowed image type (detected: %s)", detected),
		}
	}
	contentType = detected

	// Compute SHA-256 hash
	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])

	// Extract image dimensions (best-effort)
	var meta assetMetadata
	cfg, _, decodeErr := image.DecodeConfig(bytes.NewReader(data))
	if decodeErr == nil {
		meta.Width = cfg.Width
		meta.Height = cfg.Height
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return db.Asset{}, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create DB record
	asset, err := app.Queries.CreateAsset(ctx, db.CreateAssetParams{
		UserID:      userID,
		ContentType: contentType,
		Filename:    filename,
		SizeBytes:   int64(len(data)),
		Sha256:      hashHex,
		Metadata:    metaJSON,
	})
	if err != nil {
		return db.Asset{}, fmt.Errorf("failed to create asset record: %w", err)
	}

	// Store file in object storage
	key := storageKey(userID, asset.ID)
	if err := app.Storage.Put(ctx, key, bytes.NewReader(data), contentType); err != nil {
		// Best-effort cleanup of DB record
		if _, delErr := app.Queries.DeleteAsset(ctx, db.DeleteAssetParams{
			UserID: userID,
			ID:     asset.ID,
		}); delErr != nil {
			GetLogger(ctx).Error("failed to clean up asset record after storage failure",
				"error", delErr,
				"asset_id", asset.ID,
				"user_id", userID,
			)
		}
		return db.Asset{}, fmt.Errorf("failed to store asset: %w", err)
	}

	return asset, nil
}
