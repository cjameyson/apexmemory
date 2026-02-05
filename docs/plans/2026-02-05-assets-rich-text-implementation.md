# Assets & Rich Text Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add asset upload/storage infrastructure, wire image occlusion uploads, and introduce TipTap rich text editing for basic fact fields.

**Architecture:** Go `Storage` interface with local filesystem implementation, asset DB table with user-scoped ownership, multipart upload proxied through SvelteKit to Go backend. TipTap editor stores JSON in existing JSONB `content` field with `asset_ids` denormalized at content root.

**Tech Stack:** Go (net/http, pgx, crypto/sha256, image stdlib), PostgreSQL (JSONB, GIN index), SvelteKit (proxy routes), TipTap 2.x (ProseMirror-based editor), Tailwind.

---

## Task 1: Asset Database Migration

**Files:**
- Create: `backend/db/migrations/004_assets.sql`

**Step 1: Write the migration file**

```sql
-- 004_assets.sql

CREATE TABLE app.assets (
    user_id      UUID NOT NULL REFERENCES app.users(id),
    id           UUID NOT NULL DEFAULT uuidv7(),
    content_type TEXT NOT NULL,
    filename     TEXT NOT NULL,
    size_bytes   BIGINT NOT NULL,
    sha256       TEXT NOT NULL,
    metadata     JSONB NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, id)
);

-- Index for listing assets by user
CREATE INDEX ix_assets_user ON app.assets(user_id, created_at DESC);

-- Index for SHA-256 lookups (future dedup)
CREATE INDEX ix_assets_sha256 ON app.assets(user_id, sha256);

-- GIN index on facts content->'asset_ids' for orphan detection queries
CREATE INDEX ix_facts_asset_ids ON app.facts USING GIN ((content->'asset_ids'));

-- Trigger for updated_at
CREATE TRIGGER trg_assets_set_updated_at
    BEFORE UPDATE ON app.assets
    FOR EACH ROW
    EXECUTE FUNCTION app_code.set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_assets_set_updated_at ON app.assets;
DROP INDEX IF EXISTS ix_facts_asset_ids;
DROP INDEX IF EXISTS ix_assets_sha256;
DROP INDEX IF EXISTS ix_assets_user;
DROP TABLE IF EXISTS app.assets;
```

**Step 2: Apply migration**

Run: `make tern.migrate`
Expected: Migration 004 applied successfully. Schema dump updated.

**Step 3: Verify migration**

Run: `make db.psql.claude SQL="SELECT column_name, data_type FROM information_schema.columns WHERE table_schema='app' AND table_name='assets' ORDER BY ordinal_position;"`
Expected: All columns listed (user_id, id, content_type, filename, size_bytes, sha256, metadata, created_at, updated_at).

**Step 4: Commit**

```
git add backend/db/migrations/004_assets.sql backend/db/schema-dump.sql
git commit -m "feat: add assets table migration with GIN index on facts asset_ids"
```

---

## Task 2: SQLC Asset Queries

**Files:**
- Create: `backend/db/queries/assets.sql`
- Modify: `backend/sqlc.yml` (add rename for `app_asset`)

**Step 1: Write SQLC queries**

```sql
-- name: CreateAsset :one
INSERT INTO app.assets (user_id, content_type, filename, size_bytes, sha256, metadata)
VALUES (@user_id, @content_type, @filename, @size_bytes, @sha256, @metadata)
RETURNING *;

-- name: GetAsset :one
SELECT * FROM app.assets
WHERE user_id = @user_id AND id = @id;

-- name: DeleteAsset :execrows
DELETE FROM app.assets
WHERE user_id = @user_id AND id = @id;

-- name: DeleteAssets :execrows
DELETE FROM app.assets
WHERE user_id = @user_id AND id = ANY(@ids::uuid[]);

-- name: CountFactsReferencingAsset :one
SELECT count(*) FROM app.facts
WHERE user_id = @user_id
  AND id != @exclude_fact_id
  AND content->'asset_ids' @> @asset_id_json::jsonb;
```

**Step 2: Add rename to sqlc.yml**

In `backend/sqlc.yml`, add to the `rename` map:

```yaml
app_asset: Asset
```

**Step 3: Generate Go code**

Run: `make db.sqlc`
Expected: New `assets.sql.go` generated in `backend/internal/db/`.

**Step 4: Verify generation**

Run: `cd backend && go build ./...`
Expected: Clean compile.

**Step 5: Commit**

```
git add backend/db/queries/assets.sql backend/sqlc.yml backend/internal/db/
git commit -m "feat: add sqlc queries for asset CRUD and reference counting"
```

---

## Task 3: Go Storage Interface + Local Implementation

**Files:**
- Create: `backend/internal/storage/storage.go`
- Create: `backend/internal/storage/local.go`

**Step 1: Write the Storage interface**

```go
// backend/internal/storage/storage.go
package storage

import (
    "context"
    "errors"
    "io"
)

var ErrNotFound = errors.New("object not found")

// Storage is the interface for object storage backends.
type Storage interface {
    Put(ctx context.Context, key string, reader io.Reader, contentType string) error
    Get(ctx context.Context, key string) (io.ReadCloser, string, error) // returns reader, contentType, error
    Delete(ctx context.Context, key string) error
}
```

Note: `Get` returns `contentType` as the second value so the serve handler can set the `Content-Type` header without a DB lookup.

**Step 2: Write the LocalStorage implementation**

```go
// backend/internal/storage/local.go
package storage

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// LocalStorage stores objects on the local filesystem.
// Each object is stored as two files: the data file and a .meta JSON file.
type LocalStorage struct {
    basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
    if err := os.MkdirAll(basePath, 0o750); err != nil {
        return nil, fmt.Errorf("create storage directory: %w", err)
    }
    return &LocalStorage{basePath: basePath}, nil
}

type objectMeta struct {
    ContentType string `json:"content_type"`
}

func (s *LocalStorage) Put(_ context.Context, key string, reader io.Reader, contentType string) error {
    path := filepath.Join(s.basePath, key)
    if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
        return fmt.Errorf("create directory: %w", err)
    }

    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("create file: %w", err)
    }
    defer f.Close()

    if _, err := io.Copy(f, reader); err != nil {
        os.Remove(path)
        return fmt.Errorf("write file: %w", err)
    }

    // Write metadata sidecar
    meta := objectMeta{ContentType: contentType}
    metaBytes, _ := json.Marshal(meta)
    if err := os.WriteFile(path+".meta", metaBytes, 0o640); err != nil {
        os.Remove(path)
        return fmt.Errorf("write metadata: %w", err)
    }

    return nil
}

func (s *LocalStorage) Get(_ context.Context, key string) (io.ReadCloser, string, error) {
    path := filepath.Join(s.basePath, key)

    f, err := os.Open(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, "", ErrNotFound
        }
        return nil, "", fmt.Errorf("open file: %w", err)
    }

    // Read metadata sidecar
    var meta objectMeta
    metaBytes, err := os.ReadFile(path + ".meta")
    if err == nil {
        json.Unmarshal(metaBytes, &meta)
    }

    return f, meta.ContentType, nil
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
    path := filepath.Join(s.basePath, key)
    os.Remove(path + ".meta")

    if err := os.Remove(path); err != nil {
        if os.IsNotExist(err) {
            return nil // idempotent
        }
        return fmt.Errorf("delete file: %w", err)
    }
    return nil
}
```

**Step 3: Verify compilation**

Run: `cd backend && go build ./...`
Expected: Clean compile.

**Step 4: Write tests for LocalStorage**

Create `backend/internal/storage/local_test.go`:

```go
package storage

import (
    "bytes"
    "context"
    "io"
    "os"
    "testing"
)

func TestLocalStorage_PutGetDelete(t *testing.T) {
    dir := t.TempDir()
    s, err := NewLocalStorage(dir)
    if err != nil {
        t.Fatalf("NewLocalStorage: %v", err)
    }

    ctx := context.Background()
    key := "assets/test-user/test-file"
    data := []byte("hello world")
    contentType := "text/plain"

    // Put
    if err := s.Put(ctx, key, bytes.NewReader(data), contentType); err != nil {
        t.Fatalf("Put: %v", err)
    }

    // Get
    reader, ct, err := s.Get(ctx, key)
    if err != nil {
        t.Fatalf("Get: %v", err)
    }
    defer reader.Close()

    got, _ := io.ReadAll(reader)
    if !bytes.Equal(got, data) {
        t.Errorf("data mismatch: got %q, want %q", got, data)
    }
    if ct != contentType {
        t.Errorf("content type: got %q, want %q", ct, contentType)
    }

    // Delete
    if err := s.Delete(ctx, key); err != nil {
        t.Fatalf("Delete: %v", err)
    }

    // Get after delete
    _, _, err = s.Get(ctx, key)
    if err != ErrNotFound {
        t.Errorf("expected ErrNotFound after delete, got %v", err)
    }
}

func TestLocalStorage_GetNotFound(t *testing.T) {
    dir := t.TempDir()
    s, err := NewLocalStorage(dir)
    if err != nil {
        t.Fatalf("NewLocalStorage: %v", err)
    }

    _, _, err = s.Get(context.Background(), "nonexistent")
    if err != ErrNotFound {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}

func TestLocalStorage_DeleteIdempotent(t *testing.T) {
    dir := t.TempDir()
    s, err := NewLocalStorage(dir)
    if err != nil {
        t.Fatalf("NewLocalStorage: %v", err)
    }

    // Delete non-existent key should not error
    if err := s.Delete(context.Background(), "nonexistent"); err != nil {
        t.Errorf("expected nil error for idempotent delete, got %v", err)
    }
}
```

**Step 5: Run tests**

Run: `cd backend && go test ./internal/storage/...`
Expected: All 3 tests pass.

**Step 6: Commit**

```
git add backend/internal/storage/
git commit -m "feat: add Storage interface with local filesystem implementation"
```

---

## Task 4: Wire Storage Into Application

**Files:**
- Modify: `backend/internal/app/app.go` (add Storage field to Application)
- Modify: `backend/internal/app/app.go` (add StoragePath to Config)
- Modify: `backend/cmd/api/main.go` (initialize LocalStorage)

**Step 1: Add Storage field to Application struct**

In `backend/internal/app/app.go`, add import for `storage` package and a `Storage` field:

```go
import "apexmemory.ai/internal/storage"

type Application struct {
    Config         Config
    Logger         *slog.Logger
    DB             *pgxpool.Pool
    Queries        *db.Queries
    Storage        storage.Storage  // add this line
    RateLimiters   *RateLimiters
    BackgroundJobs *BackgroundJobs
    trustedProxies *TrustedProxyChecker
}
```

Add `StoragePath` to Config:

```go
type Config struct {
    Port int
    Env  string
    DB   struct {
        DSN          string
        MaxOpenConns int
        MinIdleConns int
        MaxIdleTime  time.Duration
    }
    StoragePath    string  // add this line
    TrustedProxies []string
}
```

**Step 2: Initialize LocalStorage in main.go**

In `backend/cmd/api/main.go`, after config parsing and before `app.New()`, initialize storage:

```go
storagePath := os.Getenv("STORAGE_PATH")
if storagePath == "" {
    storagePath = "./data/assets"
}
cfg.StoragePath = storagePath

store, err := storage.NewLocalStorage(cfg.StoragePath)
if err != nil {
    logger.Error("failed to initialize storage", "error", err)
    os.Exit(1)
}
```

Then set `app.Storage = store` after creating the Application.

**Step 3: Add `data/` to .gitignore**

Add `data/` to `backend/.gitignore` (or project root `.gitignore`) so local asset files aren't committed.

**Step 4: Verify compilation**

Run: `cd backend && go build ./...`
Expected: Clean compile.

**Step 5: Run existing tests**

Run: `cd backend && go test ./internal/app/...`
Expected: All existing tests still pass. Tests use `testApp()` which doesn't set Storage - that's fine for now. Asset handlers will check `app.Storage != nil`.

**Step 6: Commit**

```
git add backend/internal/app/app.go backend/cmd/api/main.go .gitignore
git commit -m "feat: wire Storage interface into Application struct"
```

---

## Task 5: Asset Upload Handler

**Files:**
- Create: `backend/internal/app/assets.go` (business logic)
- Create: `backend/internal/app/assets_handlers.go` (HTTP handlers)
- Modify: `backend/internal/app/routes.go` (register routes)

**Step 1: Write asset business logic**

Create `backend/internal/app/assets.go`:

```go
package app

import (
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
    "strings"

    "apexmemory.ai/internal/db"
    "github.com/google/uuid"
)

const (
    maxUploadSize = 10 << 20 // 10MB
)

var allowedContentTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/webp": true,
    "image/gif":  true,
}

var errAssetNotFound = fmt.Errorf("asset not found")

func storageKey(userID, assetID uuid.UUID) string {
    return fmt.Sprintf("assets/%s/%s", userID, assetID)
}

type AssetMetadata struct {
    Width  int `json:"width,omitempty"`
    Height int `json:"height,omitempty"`
}

// UploadAsset reads the file, computes hash, extracts dimensions, stores in storage + DB.
func (app *Application) UploadAsset(ctx context.Context, userID uuid.UUID, file io.Reader, filename, contentType string, size int64) (db.Asset, error) {
    if !allowedContentTypes[contentType] {
        return db.Asset{}, &FactValidationError{Message: fmt.Sprintf("unsupported content type: %s", contentType)}
    }
    if size > maxUploadSize {
        return db.Asset{}, &FactValidationError{Message: fmt.Sprintf("file too large: %d bytes (max %d)", size, maxUploadSize)}
    }

    // Read file into memory for hashing + dimension extraction
    data, err := io.ReadAll(io.LimitReader(file, maxUploadSize+1))
    if err != nil {
        return db.Asset{}, fmt.Errorf("read upload: %w", err)
    }
    if int64(len(data)) > maxUploadSize {
        return db.Asset{}, &FactValidationError{Message: "file too large"}
    }

    // Compute SHA-256
    hash := sha256.Sum256(data)
    hashHex := hex.EncodeToString(hash[:])

    // Extract image dimensions
    var meta AssetMetadata
    if strings.HasPrefix(contentType, "image/") {
        cfg, _, err := image.DecodeConfig(strings.NewReader(string(data)))
        if err == nil {
            meta.Width = cfg.Width
            meta.Height = cfg.Height
        }
    }

    metaJSON, err := json.Marshal(meta)
    if err != nil {
        return db.Asset{}, fmt.Errorf("marshal metadata: %w", err)
    }

    // Create DB record first to get the ID
    asset, err := app.Queries.CreateAsset(ctx, db.CreateAssetParams{
        UserID:      userID,
        ContentType: contentType,
        Filename:    filename,
        SizeBytes:   int64(len(data)),
        Sha256:      hashHex,
        Metadata:    metaJSON,
    })
    if err != nil {
        return db.Asset{}, fmt.Errorf("create asset record: %w", err)
    }

    // Store file
    key := storageKey(userID, asset.ID)
    if err := app.Storage.Put(ctx, key, strings.NewReader(string(data)), contentType); err != nil {
        // Best-effort cleanup of DB record
        app.Queries.DeleteAsset(ctx, db.DeleteAssetParams{UserID: userID, ID: asset.ID})
        return db.Asset{}, fmt.Errorf("store file: %w", err)
    }

    return asset, nil
}
```

Note: Using `strings.NewReader(string(data))` for the dimension decode is a workaround. A cleaner approach uses `bytes.NewReader(data)` - replace `strings.NewReader(string(data))` with `bytes.NewReader(data)` throughout. Import `bytes` instead of `strings` where applicable.

**Step 2: Write asset handlers**

Create `backend/internal/app/assets_handlers.go`:

```go
package app

import (
    "errors"
    "io"
    "net/http"
    "time"

    "apexmemory.ai/internal/storage"
)

type AssetResponse struct {
    ID          string          `json:"id"`
    ContentType string          `json:"content_type"`
    Filename    string          `json:"filename"`
    SizeBytes   int64           `json:"size_bytes"`
    Sha256      string          `json:"sha256"`
    Metadata    json.RawMessage `json:"metadata"`
    CreatedAt   time.Time       `json:"created_at"`
}

func (app *Application) UploadAssetHandler(w http.ResponseWriter, r *http.Request) {
    user := app.GetUser(r.Context())
    if user.IsAnonymous() {
        app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
        return
    }

    // Limit request body size
    r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+1024) // extra for multipart headers

    if err := r.ParseMultipartForm(maxUploadSize); err != nil {
        app.RespondError(w, r, http.StatusBadRequest, "Invalid multipart form or file too large")
        return
    }
    defer r.MultipartForm.RemoveAll()

    file, header, err := r.FormFile("file")
    if err != nil {
        app.RespondError(w, r, http.StatusBadRequest, "Missing file field")
        return
    }
    defer file.Close()

    contentType := header.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    asset, err := app.UploadAsset(r.Context(), user.ID, file, header.Filename, contentType, header.Size)
    if err != nil {
        var validationErr *FactValidationError
        if errors.As(err, &validationErr) {
            app.RespondError(w, r, http.StatusBadRequest, err.Error())
            return
        }
        app.RespondServerError(w, r, ErrInternal("upload asset", err))
        return
    }

    resp := AssetResponse{
        ID:          asset.ID.String(),
        ContentType: asset.ContentType,
        Filename:    asset.Filename,
        SizeBytes:   asset.SizeBytes,
        Sha256:      asset.Sha256,
        Metadata:    asset.Metadata,
        CreatedAt:   asset.CreatedAt,
    }

    app.RespondJSON(w, r, http.StatusCreated, resp)
}

func (app *Application) ServeAssetHandler(w http.ResponseWriter, r *http.Request) {
    user := app.GetUser(r.Context())
    if user.IsAnonymous() {
        app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
        return
    }

    assetID, ok := app.PathUUID(w, r, "id")
    if !ok {
        return
    }

    // Check asset exists and belongs to user
    asset, err := app.Queries.GetAsset(r.Context(), db.GetAssetParams{
        UserID: user.ID,
        ID:     assetID,
    })
    if err != nil {
        app.RespondError(w, r, http.StatusNotFound, "Asset not found")
        return
    }

    // Check If-None-Match for 304
    if r.Header.Get("If-None-Match") == `"`+asset.Sha256+`"` {
        w.WriteHeader(http.StatusNotModified)
        return
    }

    key := storageKey(user.ID, assetID)
    reader, contentType, err := app.Storage.Get(r.Context(), key)
    if err != nil {
        if errors.Is(err, storage.ErrNotFound) {
            app.RespondError(w, r, http.StatusNotFound, "Asset file not found")
            return
        }
        app.RespondServerError(w, r, ErrInternal("get asset file", err))
        return
    }
    defer reader.Close()

    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
    w.Header().Set("ETag", `"`+asset.Sha256+`"`)
    w.WriteHeader(http.StatusOK)
    io.Copy(w, reader)
}
```

**Step 3: Register routes**

In `backend/internal/app/routes.go`, add to the protected routes section:

```go
// Assets
mux.Handle("POST /v1/assets", protected.ThenFunc(app.UploadAssetHandler))
mux.Handle("GET /v1/assets/{id}/file", protected.ThenFunc(app.ServeAssetHandler))
```

**Step 4: Verify compilation**

Run: `cd backend && go build ./...`
Expected: Clean compile.

**Step 5: Commit**

```
git add backend/internal/app/assets.go backend/internal/app/assets_handlers.go backend/internal/app/routes.go
git commit -m "feat: add asset upload and serve endpoints"
```

---

## Task 6: Asset Handler Integration Tests

**Files:**
- Create: `backend/internal/app/assets_test.go`

**Step 1: Write integration tests**

```go
//go:build integration

package app

import (
    "bytes"
    "encoding/json"
    "image"
    "image/color"
    "image/png"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "testing"

    "apexmemory.ai/internal/storage"
)

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

func createTestPNG(t *testing.T, width, height int) []byte {
    t.Helper()
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
        }
    }
    var buf bytes.Buffer
    if err := png.Encode(&buf, img); err != nil {
        t.Fatalf("encode png: %v", err)
    }
    return buf.Bytes()
}

func multipartRequest(t *testing.T, fieldName, filename, contentType string, data []byte) (*http.Request, string) {
    t.Helper()
    var body bytes.Buffer
    writer := multipart.NewWriter(&body)
    part, err := writer.CreateFormFile(fieldName, filename)
    if err != nil {
        t.Fatalf("CreateFormFile: %v", err)
    }
    part.Write(data)
    writer.Close()

    req := httptest.NewRequest("POST", "/v1/assets", &body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    return req, writer.FormDataContentType()
}

func TestUploadAssetHandler(t *testing.T) {
    app := testAppWithStorage(t)
    user := createTestUser(t, app)

    pngData := createTestPNG(t, 100, 50)

    tests := []struct {
        name       string
        filename   string
        data       []byte
        wantStatus int
    }{
        {
            name:       "valid PNG upload",
            filename:   "test.png",
            data:       pngData,
            wantStatus: http.StatusCreated,
        },
        {
            name:       "missing file",
            filename:   "",
            data:       nil,
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.data == nil {
                // Send empty multipart form
                req := httptest.NewRequest("POST", "/v1/assets", nil)
                req.Header.Set("Content-Type", "multipart/form-data; boundary=xxx")
                req = app.WithUser(req, user)
                rr := httptest.NewRecorder()
                app.UploadAssetHandler(rr, req)
                if rr.Code != tt.wantStatus {
                    t.Errorf("status: got %d, want %d. Body: %s", rr.Code, tt.wantStatus, rr.Body.String())
                }
                return
            }

            req, _ := multipartRequest(t, "file", tt.filename, "image/png", tt.data)
            req = app.WithUser(req, user)
            rr := httptest.NewRecorder()

            app.UploadAssetHandler(rr, req)

            if rr.Code != tt.wantStatus {
                t.Errorf("status: got %d, want %d. Body: %s", rr.Code, tt.wantStatus, rr.Body.String())
            }

            if tt.wantStatus == http.StatusCreated {
                var resp AssetResponse
                json.Unmarshal(rr.Body.Bytes(), &resp)
                if resp.ID == "" {
                    t.Error("expected non-empty ID")
                }
                if resp.Sha256 == "" {
                    t.Error("expected non-empty sha256")
                }
                if resp.ContentType != "image/png" {
                    t.Errorf("content_type: got %q, want %q", resp.ContentType, "image/png")
                }
            }
        })
    }
}

func TestServeAssetHandler(t *testing.T) {
    app := testAppWithStorage(t)
    user := createTestUser(t, app)

    // Upload an asset first
    pngData := createTestPNG(t, 100, 50)
    asset, err := app.UploadAsset(t.Context(), user.ID, bytes.NewReader(pngData), "test.png", "image/png", int64(len(pngData)))
    if err != nil {
        t.Fatalf("UploadAsset: %v", err)
    }

    t.Run("serves file with correct headers", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/assets/"+asset.ID.String()+"/file", nil)
        req.SetPathValue("id", asset.ID.String())
        req = app.WithUser(req, user)
        rr := httptest.NewRecorder()

        app.ServeAssetHandler(rr, req)

        if rr.Code != http.StatusOK {
            t.Fatalf("status: got %d, want %d", rr.Code, http.StatusOK)
        }
        if rr.Header().Get("Content-Type") != "image/png" {
            t.Errorf("Content-Type: got %q", rr.Header().Get("Content-Type"))
        }
        if rr.Header().Get("Cache-Control") == "" {
            t.Error("expected Cache-Control header")
        }
        if rr.Header().Get("ETag") == "" {
            t.Error("expected ETag header")
        }

        body, _ := io.ReadAll(rr.Body)
        if !bytes.Equal(body, pngData) {
            t.Errorf("body length mismatch: got %d, want %d", len(body), len(pngData))
        }
    })

    t.Run("returns 304 with matching ETag", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/assets/"+asset.ID.String()+"/file", nil)
        req.SetPathValue("id", asset.ID.String())
        req.Header.Set("If-None-Match", `"`+asset.Sha256+`"`)
        req = app.WithUser(req, user)
        rr := httptest.NewRecorder()

        app.ServeAssetHandler(rr, req)

        if rr.Code != http.StatusNotModified {
            t.Errorf("status: got %d, want %d", rr.Code, http.StatusNotModified)
        }
    })

    t.Run("returns 404 for unknown asset", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/assets/019501a0-0000-7000-8000-000000000099/file", nil)
        req.SetPathValue("id", "019501a0-0000-7000-8000-000000000099")
        req = app.WithUser(req, user)
        rr := httptest.NewRecorder()

        app.ServeAssetHandler(rr, req)

        if rr.Code != http.StatusNotFound {
            t.Errorf("status: got %d, want %d", rr.Code, http.StatusNotFound)
        }
    })
}
```

**Step 2: Run tests**

Run: `cd backend && go test -tags=integration ./internal/app/ -run TestUploadAsset -v && go test -tags=integration ./internal/app/ -run TestServeAsset -v`
Expected: All tests pass.

**Step 3: Commit**

```
git add backend/internal/app/assets_test.go
git commit -m "test: add integration tests for asset upload and serve handlers"
```

---

## Task 7: SvelteKit Asset Proxy Routes + Client Helper

**Files:**
- Create: `frontend/src/routes/api/assets/+server.ts`
- Create: `frontend/src/routes/api/assets/[id]/file/+server.ts`
- Create: `frontend/src/lib/api/client.ts`

**Step 1: Create upload proxy route**

```typescript
// frontend/src/routes/api/assets/+server.ts
import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_BASE_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ request, locals }) => {
    const token = locals.sessionToken;
    if (!token) {
        error(401, { message: 'Unauthorized' });
    }

    const formData = await request.formData();

    const response = await fetch(`${API_BASE_URL}/v1/assets`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });

    if (!response.ok) {
        const err = await response.json().catch(() => ({ error: 'Upload failed' }));
        error(response.status, { message: err.error });
    }

    const data = await response.json();
    return json(data, { status: 201 });
};
```

**Step 2: Create serve proxy route**

```typescript
// frontend/src/routes/api/assets/[id]/file/+server.ts
import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_BASE_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ params, request, locals }) => {
    const token = locals.sessionToken;
    if (!token) {
        error(401, { message: 'Unauthorized' });
    }

    const headers: Record<string, string> = {
        'Authorization': `Bearer ${token}`
    };

    // Forward If-None-Match for 304 support
    const ifNoneMatch = request.headers.get('If-None-Match');
    if (ifNoneMatch) {
        headers['If-None-Match'] = ifNoneMatch;
    }

    const response = await fetch(`${API_BASE_URL}/v1/assets/${params.id}/file`, { headers });

    if (response.status === 304) {
        return new Response(null, { status: 304 });
    }

    if (!response.ok) {
        error(response.status, { message: 'Asset not found' });
    }

    return new Response(response.body, {
        status: 200,
        headers: {
            'Content-Type': response.headers.get('Content-Type') ?? 'application/octet-stream',
            'Cache-Control': response.headers.get('Cache-Control') ?? '',
            'ETag': response.headers.get('ETag') ?? ''
        }
    });
};
```

**Step 3: Create client helper**

```typescript
// frontend/src/lib/api/client.ts

export interface ApiAsset {
    id: string;
    content_type: string;
    filename: string;
    size_bytes: number;
    sha256: string;
    metadata: {
        width?: number;
        height?: number;
    };
    created_at: string;
}

export async function uploadAsset(file: File): Promise<ApiAsset> {
    const formData = new FormData();
    formData.append('file', file);

    const res = await fetch('/api/assets', {
        method: 'POST',
        body: formData
    });

    if (!res.ok) {
        const err = await res.json().catch(() => ({ message: 'Upload failed' }));
        throw new Error(err.message ?? 'Upload failed');
    }

    return res.json();
}

export function assetUrl(assetId: string): string {
    return `/api/assets/${assetId}/file`;
}
```

**Step 4: Add ApiAsset type to existing types file**

In `frontend/src/lib/api/types.ts`, add the `ApiAsset` export or re-export from `client.ts`. Depending on existing patterns, you may want to keep all API types in `types.ts` and import them in `client.ts`.

**Step 5: Verify frontend compiles**

Run: `cd frontend && npm run check`
Expected: No type errors.

**Step 6: Commit**

```
git add frontend/src/routes/api/assets/ frontend/src/lib/api/client.ts
git commit -m "feat: add SvelteKit proxy routes and client helper for asset upload/serve"
```

---

## Task 8: Wire ImageUploader to Real Upload Flow

**Files:**
- Modify: `frontend/src/lib/components/image-occlusion/ImageUploader.svelte`
- Modify: `frontend/src/lib/components/image-occlusion/types.ts` (if needed)

**Step 1: Implement ImageUploader with real upload**

Replace the stub in `ImageUploader.svelte` with a working implementation:

```svelte
<script lang="ts">
    import { Upload } from '@lucide/svelte';
    import { uploadAsset, assetUrl } from '$lib/api/client';

    interface Props {
        onImageLoad?: (url: string, width: number, height: number, assetId: string) => void;
    }

    let { onImageLoad }: Props = $props();

    let isDragging = $state(false);
    let uploading = $state(false);
    let error = $state<string | null>(null);

    async function handleFile(file: File) {
        if (!file.type.startsWith('image/')) {
            error = 'Please select an image file (JPEG, PNG, WebP, or GIF)';
            return;
        }

        error = null;
        uploading = true;

        try {
            const asset = await uploadAsset(file);
            const url = assetUrl(asset.id);
            const width = asset.metadata?.width ?? 0;
            const height = asset.metadata?.height ?? 0;
            onImageLoad?.(url, width, height, asset.id);
        } catch (err) {
            error = err instanceof Error ? err.message : 'Upload failed';
        } finally {
            uploading = false;
        }
    }

    function handleDrop(e: DragEvent) {
        e.preventDefault();
        isDragging = false;
        const file = e.dataTransfer?.files[0];
        if (file) handleFile(file);
    }

    function handleFileInput(e: Event) {
        const input = e.target as HTMLInputElement;
        const file = input.files?.[0];
        if (file) handleFile(file);
        input.value = '';
    }

    let fileInput: HTMLInputElement;
</script>

<div
    class="flex h-full w-full flex-col items-center justify-center p-8 {isDragging ? 'bg-primary/5' : ''}"
    ondragover={(e) => { e.preventDefault(); isDragging = true; }}
    ondragleave={() => (isDragging = false)}
    ondrop={handleDrop}
    role="button"
    tabindex="0"
    onclick={() => fileInput?.click()}
    onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') fileInput?.click(); }}
>
    <input
        bind:this={fileInput}
        type="file"
        accept="image/jpeg,image/png,image/webp,image/gif"
        class="hidden"
        onchange={handleFileInput}
    />

    <div class="flex max-w-md flex-col items-center rounded-xl border-2 border-dashed border-border p-12 text-center">
        {#if uploading}
            <div class="text-muted-foreground">Uploading...</div>
        {:else}
            <Upload class="mb-4 h-12 w-12 text-muted-foreground" />
            <p class="text-foreground mb-2 font-medium">Drop an image here or click to browse</p>
            <p class="text-muted-foreground text-sm">JPEG, PNG, WebP, or GIF up to 10MB</p>
        {/if}

        {#if error}
            <p class="text-destructive mt-4 text-sm">{error}</p>
        {/if}
    </div>
</div>
```

**Step 2: Update onImageLoad callback signature**

The `onImageLoad` callback now passes `assetId` as a fourth parameter. Update the `ImageOcclusionEditor.svelte` to handle this - set `image.assetId` in the editor state when an image is loaded.

Check `ImageOcclusionEditor.svelte` for where `onImageLoad` is handled and update accordingly. The `ImageData` type in `types.ts` already has `assetId?: string` so no type change needed.

**Step 3: Verify frontend compiles**

Run: `cd frontend && npm run check`
Expected: No type errors.

**Step 4: Commit**

```
git add frontend/src/lib/components/image-occlusion/ImageUploader.svelte
git commit -m "feat: wire ImageUploader to real asset upload flow"
```

---

## Task 9: Install TipTap Dependencies

**Files:**
- Modify: `frontend/package.json`

**Step 1: Install TipTap packages**

Run:
```bash
cd frontend && npm install @tiptap/core @tiptap/pm @tiptap/starter-kit @tiptap/extension-image @tiptap/extension-underline
```

Note: `@tiptap/pm` provides ProseMirror dependencies. `@tiptap/starter-kit` bundles bold, italic, lists, headings, etc. `@tiptap/extension-underline` is separate from starter-kit.

**Step 2: Verify build**

Run: `cd frontend && npm run check`
Expected: No errors.

**Step 3: Commit**

```
git add frontend/package.json frontend/package-lock.json
git commit -m "chore: install TipTap editor dependencies"
```

---

## Task 10: RichTextEditor Component

**Files:**
- Create: `frontend/src/lib/components/rich-text/RichTextEditor.svelte`
- Create: `frontend/src/lib/components/rich-text/EditorToolbar.svelte`
- Create: `frontend/src/lib/components/rich-text/index.ts`

**Step 1: Create the TipTap editor component**

```svelte
<!-- frontend/src/lib/components/rich-text/RichTextEditor.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { Editor } from '@tiptap/core';
    import StarterKit from '@tiptap/starter-kit';
    import Underline from '@tiptap/extension-underline';
    import Image from '@tiptap/extension-image';
    import EditorToolbar from './EditorToolbar.svelte';
    import { uploadAsset, assetUrl } from '$lib/api/client';

    interface Props {
        content: Record<string, unknown> | null;
        onchange: (content: Record<string, unknown>) => void;
        placeholder?: string;
        class?: string;
    }

    let {
        content,
        onchange,
        placeholder = '',
        class: className = ''
    }: Props = $props();

    let element: HTMLDivElement;
    let editor: Editor | null = $state(null);

    const CustomImage = Image.extend({
        addAttributes() {
            return {
                ...this.parent?.(),
                'asset_id': { default: null },
            };
        },
    });

    onMount(() => {
        editor = new Editor({
            element,
            extensions: [
                StarterKit,
                Underline,
                CustomImage.configure({
                    inline: false,
                    allowBase64: false,
                }),
            ],
            content: content ?? { type: 'doc', content: [{ type: 'paragraph' }] },
            onUpdate: ({ editor: e }) => {
                onchange(e.getJSON() as Record<string, unknown>);
            },
            editorProps: {
                attributes: {
                    class: 'prose prose-sm dark:prose-invert max-w-none focus:outline-none min-h-[80px] px-3 py-2',
                },
            },
        });
    });

    onDestroy(() => {
        editor?.destroy();
    });

    async function handleImageInsert() {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'image/jpeg,image/png,image/webp,image/gif';
        input.onchange = async () => {
            const file = input.files?.[0];
            if (!file || !editor) return;

            try {
                const asset = await uploadAsset(file);
                editor.chain().focus().setImage({
                    src: assetUrl(asset.id),
                    alt: asset.filename,
                    'asset_id': asset.id,
                } as any).run();
            } catch {
                // Error handling is minimal for now
            }
        };
        input.click();
    }

    async function handleDrop(e: DragEvent) {
        if (!e.dataTransfer?.files.length || !editor) return;

        const file = e.dataTransfer.files[0];
        if (!file.type.startsWith('image/')) return;

        e.preventDefault();

        try {
            const asset = await uploadAsset(file);
            editor.chain().focus().setImage({
                src: assetUrl(asset.id),
                alt: asset.filename,
                'asset_id': asset.id,
            } as any).run();
        } catch {
            // Error handling minimal for now
        }
    }

    async function handlePaste(e: ClipboardEvent) {
        if (!e.clipboardData?.files.length || !editor) return;

        const file = e.clipboardData.files[0];
        if (!file.type.startsWith('image/')) return;

        e.preventDefault();

        try {
            const asset = await uploadAsset(file);
            editor.chain().focus().setImage({
                src: assetUrl(asset.id),
                alt: asset.filename,
                'asset_id': asset.id,
            } as any).run();
        } catch {
            // Error handling minimal for now
        }
    }
</script>

<div class="rounded-md border border-border bg-card {className}">
    {#if editor}
        <EditorToolbar {editor} onImageInsert={handleImageInsert} />
    {/if}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        bind:this={element}
        ondrop={handleDrop}
        onpaste={handlePaste}
    ></div>
</div>
```

**Step 2: Create the toolbar component**

```svelte
<!-- frontend/src/lib/components/rich-text/EditorToolbar.svelte -->
<script lang="ts">
    import type { Editor } from '@tiptap/core';
    import {
        Bold, Italic, Underline, List, ListOrdered,
        Heading2, Image, Undo2, Redo2
    } from '@lucide/svelte';

    interface Props {
        editor: Editor;
        onImageInsert: () => void;
    }

    let { editor, onImageInsert }: Props = $props();
</script>

<div class="flex flex-wrap gap-0.5 border-b border-border p-1">
    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted {editor.isActive('bold') ? 'bg-muted text-foreground' : 'text-muted-foreground'}"
        onclick={() => editor.chain().focus().toggleBold().run()}
        title="Bold"
    >
        <Bold class="h-4 w-4" />
    </button>
    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted {editor.isActive('italic') ? 'bg-muted text-foreground' : 'text-muted-foreground'}"
        onclick={() => editor.chain().focus().toggleItalic().run()}
        title="Italic"
    >
        <Italic class="h-4 w-4" />
    </button>
    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted {editor.isActive('underline') ? 'bg-muted text-foreground' : 'text-muted-foreground'}"
        onclick={() => editor.chain().focus().toggleUnderline().run()}
        title="Underline"
    >
        <Underline class="h-4 w-4" />
    </button>

    <div class="mx-1 w-px bg-border"></div>

    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted {editor.isActive('heading', { level: 2 }) ? 'bg-muted text-foreground' : 'text-muted-foreground'}"
        onclick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        title="Heading"
    >
        <Heading2 class="h-4 w-4" />
    </button>
    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted {editor.isActive('bulletList') ? 'bg-muted text-foreground' : 'text-muted-foreground'}"
        onclick={() => editor.chain().focus().toggleBulletList().run()}
        title="Bullet List"
    >
        <List class="h-4 w-4" />
    </button>
    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted {editor.isActive('orderedList') ? 'bg-muted text-foreground' : 'text-muted-foreground'}"
        onclick={() => editor.chain().focus().toggleOrderedList().run()}
        title="Numbered List"
    >
        <ListOrdered class="h-4 w-4" />
    </button>

    <div class="mx-1 w-px bg-border"></div>

    <button
        type="button"
        class="rounded p-1.5 hover:bg-muted text-muted-foreground"
        onclick={onImageInsert}
        title="Insert Image"
    >
        <Image class="h-4 w-4" />
    </button>

    <div class="ml-auto flex gap-0.5">
        <button
            type="button"
            class="rounded p-1.5 hover:bg-muted text-muted-foreground disabled:opacity-30"
            onclick={() => editor.chain().focus().undo().run()}
            disabled={!editor.can().undo()}
            title="Undo"
        >
            <Undo2 class="h-4 w-4" />
        </button>
        <button
            type="button"
            class="rounded p-1.5 hover:bg-muted text-muted-foreground disabled:opacity-30"
            onclick={() => editor.chain().focus().redo().run()}
            disabled={!editor.can().redo()}
            title="Redo"
        >
            <Redo2 class="h-4 w-4" />
        </button>
    </div>
</div>
```

**Step 3: Create barrel export**

```typescript
// frontend/src/lib/components/rich-text/index.ts
export { default as RichTextEditor } from './RichTextEditor.svelte';
```

**Step 4: Verify frontend compiles**

Run: `cd frontend && npm run check`
Expected: No type errors.

**Step 5: Commit**

```
git add frontend/src/lib/components/rich-text/
git commit -m "feat: add RichTextEditor component with TipTap and image support"
```

---

## Task 11: Update Basic Fact Editor to Use Rich Text

**Files:**
- Modify: `frontend/src/lib/components/facts/basic-fact-editor.svelte`

**Step 1: Read the current basic-fact-editor.svelte**

Read the file to understand current structure before modifying.

**Step 2: Replace text inputs with RichTextEditor**

Update the `front` and `back` fields to use `RichTextEditor` instead of plain text inputs. Keep `hint` and `backExtra` as plain text for now.

Key changes:
- Import `RichTextEditor` from `$lib/components/rich-text`
- Change `front` and `back` from `string` to TipTap JSON (`Record<string, unknown> | null`)
- Update the `BasicFactData` interface to reflect the new types
- Wire `onchange` to emit the updated data

The `notify()` function should still call the parent's `onchange` with the updated data. The parent (`create-fact-modal.svelte`) will need to build content with `type: "rich_text"` for front/back fields instead of `type: "plain_text"`.

**Step 3: Update create-fact-modal to build rich_text content**

When building the fact content for submission, the modal should:
- Use `type: "rich_text"` for front/back fields
- Set the `value` to the TipTap JSON doc
- Collect `asset_ids` by walking the TipTap JSON for image nodes with `asset_id` attributes
- Include `asset_ids` at the content root

Helper function for extracting asset IDs from TipTap JSON:

```typescript
function extractAssetIds(doc: Record<string, unknown>): string[] {
    const ids: string[] = [];
    function walk(node: any) {
        if (node.type === 'image' && node.attrs?.asset_id) {
            ids.push(node.attrs.asset_id);
        }
        if (Array.isArray(node.content)) {
            node.content.forEach(walk);
        }
    }
    walk(doc);
    return [...new Set(ids)];
}
```

**Step 4: Verify frontend compiles**

Run: `cd frontend && npm run check`
Expected: No type errors.

**Step 5: Commit**

```
git add frontend/src/lib/components/facts/basic-fact-editor.svelte frontend/src/lib/components/facts/create-fact-modal.svelte
git commit -m "feat: use RichTextEditor for basic fact front/back fields"
```

---

## Task 12: RichTextContent (Read-Only Renderer)

**Files:**
- Create: `frontend/src/lib/components/rich-text/RichTextContent.svelte`
- Modify: `frontend/src/lib/components/rich-text/index.ts`

**Step 1: Create the read-only renderer**

This component renders TipTap JSON to HTML without loading the full editor. Uses `generateHTML` from `@tiptap/core`:

```svelte
<!-- frontend/src/lib/components/rich-text/RichTextContent.svelte -->
<script lang="ts">
    import { generateHTML } from '@tiptap/core';
    import StarterKit from '@tiptap/starter-kit';
    import Underline from '@tiptap/extension-underline';
    import Image from '@tiptap/extension-image';
    import { assetUrl } from '$lib/api/client';

    interface Props {
        content: Record<string, unknown>;
        class?: string;
    }

    let { content, class: className = '' }: Props = $props();

    const CustomImage = Image.extend({
        addAttributes() {
            return {
                ...this.parent?.(),
                'asset_id': { default: null },
            };
        },
        renderHTML({ HTMLAttributes }) {
            const assetId = HTMLAttributes['asset_id'];
            if (assetId) {
                HTMLAttributes.src = assetUrl(assetId);
            }
            return ['img', HTMLAttributes];
        },
    });

    let html = $derived(
        generateHTML(content as any, [StarterKit, Underline, CustomImage])
    );
</script>

<div class="prose prose-sm dark:prose-invert max-w-none {className}">
    {@html html}
</div>
```

**Step 2: Update barrel export**

```typescript
// frontend/src/lib/components/rich-text/index.ts
export { default as RichTextEditor } from './RichTextEditor.svelte';
export { default as RichTextContent } from './RichTextContent.svelte';
```

**Step 3: Verify frontend compiles**

Run: `cd frontend && npm run check`
Expected: No type errors.

**Step 4: Commit**

```
git add frontend/src/lib/components/rich-text/
git commit -m "feat: add RichTextContent component for read-only rendering"
```

---

## Task 13: Backend Validation for rich_text Fields + asset_ids

**Files:**
- Modify: `backend/internal/app/facts.go`

**Step 1: Read current validation logic**

Read `facts.go` to understand the current `validateFactContent` and `extractElementIDs` functions.

**Step 2: Update validation to handle asset_ids**

Add validation for the `asset_ids` array in the content root. When `asset_ids` is present:
- Validate it's an array of strings
- Validate each string is a valid UUID
- Verify each asset exists and belongs to the user (requires passing `ctx`, `userID`, and `queries` to validation - this changes the `validateFactContent` signature)

For `rich_text` fields:
- Validate that `value` is an object (not a string)
- Validate it has `type: "doc"` at minimum

Update `validateFactContent` to accept context and queries for asset validation:

```go
func (app *Application) validateFactContent(ctx context.Context, userID uuid.UUID, factType string, content json.RawMessage) ([]string, error) {
    // ... existing validation ...

    // Validate asset_ids if present
    assetIDs, err := extractAssetIDs(content)
    if err != nil {
        return nil, err
    }

    for _, assetID := range assetIDs {
        _, err := app.Queries.GetAsset(ctx, db.GetAssetParams{
            UserID: userID,
            ID:     assetID,
        })
        if err != nil {
            return nil, &FactValidationError{
                Message: fmt.Sprintf("asset %s not found", assetID),
            }
        }
    }

    return elementIDs, nil
}
```

Add `extractAssetIDs` helper:

```go
func extractAssetIDs(content json.RawMessage) ([]uuid.UUID, error) {
    var parsed struct {
        AssetIDs []string `json:"asset_ids"`
    }
    if err := json.Unmarshal(content, &parsed); err != nil {
        return nil, &FactValidationError{Message: "invalid content JSON"}
    }

    if len(parsed.AssetIDs) == 0 {
        return nil, nil
    }

    ids := make([]uuid.UUID, 0, len(parsed.AssetIDs))
    for _, raw := range parsed.AssetIDs {
        id, err := uuid.Parse(raw)
        if err != nil {
            return nil, &FactValidationError{Message: fmt.Sprintf("invalid asset ID: %s", raw)}
        }
        ids = append(ids, id)
    }
    return ids, nil
}
```

**Step 3: Update callers of validateFactContent**

Update `CreateFact` and `UpdateFact` to pass `ctx` and `userID` to the new validation signature.

**Step 4: Run existing tests**

Run: `cd backend && go test -tags=integration ./internal/app/ -v`
Expected: All existing tests still pass (they don't use `asset_ids`).

**Step 5: Add tests for asset_ids validation**

Add test cases to `facts_test.go`:
- Basic fact with `asset_ids` referencing a valid uploaded asset
- Basic fact with `asset_ids` referencing a non-existent asset (should fail)
- Basic fact with invalid UUID in `asset_ids` (should fail)

**Step 6: Run tests**

Run: `cd backend && go test -tags=integration ./internal/app/ -v`
Expected: All tests pass.

**Step 7: Commit**

```
git add backend/internal/app/facts.go backend/internal/app/facts_test.go
git commit -m "feat: validate asset_ids in fact content and verify asset ownership"
```

---

## Task 14: Asset Cleanup on Fact Delete/Update

**Files:**
- Modify: `backend/internal/app/facts.go` (DeleteFact, UpdateFact functions)

**Step 1: Read current DeleteFact and UpdateFact**

Read `facts.go` to understand the current implementation.

**Step 2: Add cleanup to DeleteFact**

After deleting the fact, check each of its `asset_ids` for orphans:

```go
func (app *Application) DeleteFact(ctx context.Context, userID, notebookID, factID uuid.UUID) error {
    // Read fact first to get asset_ids
    fact, err := app.Queries.GetFact(ctx, db.GetFactParams{
        UserID: userID, ID: factID, NotebookID: notebookID,
    })
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return errFactNotFound
        }
        return err
    }

    assetIDs, _ := extractAssetIDs(fact.Content)

    // Delete the fact
    rows, err := app.Queries.DeleteFact(ctx, db.DeleteFactParams{
        UserID: userID, ID: factID, NotebookID: notebookID,
    })
    if err != nil {
        return err
    }
    if rows == 0 {
        return errFactNotFound
    }

    // Cleanup orphaned assets (best-effort, after fact delete)
    app.cleanupOrphanedAssets(ctx, userID, factID, assetIDs)

    return nil
}
```

**Step 3: Add cleanup helper**

```go
func (app *Application) cleanupOrphanedAssets(ctx context.Context, userID uuid.UUID, excludeFactID uuid.UUID, assetIDs []uuid.UUID) {
    for _, assetID := range assetIDs {
        assetIDJSON, _ := json.Marshal([]string{assetID.String()})
        count, err := app.Queries.CountFactsReferencingAsset(ctx, db.CountFactsReferencingAssetParams{
            UserID:        userID,
            ExcludeFactID: excludeFactID,
            AssetIdJson:   assetIDJSON,
        })
        if err != nil || count > 0 {
            continue // still referenced or query failed, skip
        }

        // Delete from storage + DB
        key := storageKey(userID, assetID)
        app.Storage.Delete(ctx, key)
        app.Queries.DeleteAsset(ctx, db.DeleteAssetParams{UserID: userID, ID: assetID})
    }
}
```

**Step 4: Add cleanup to UpdateFact**

In UpdateFact, compare old and new `asset_ids`. For any IDs removed, run the orphan check:

```go
// In UpdateFact, after successful update:
oldAssetIDs, _ := extractAssetIDs(oldFact.Content)
newAssetIDs, _ := extractAssetIDs(input.Content)

removedIDs := diffAssetIDs(oldAssetIDs, newAssetIDs)
if len(removedIDs) > 0 {
    app.cleanupOrphanedAssets(ctx, userID, factID, removedIDs)
}
```

Add helper:

```go
func diffAssetIDs(old, new []uuid.UUID) []uuid.UUID {
    newSet := make(map[uuid.UUID]bool, len(new))
    for _, id := range new {
        newSet[id] = true
    }
    var removed []uuid.UUID
    for _, id := range old {
        if !newSet[id] {
            removed = append(removed, id)
        }
    }
    return removed
}
```

**Step 5: Run tests**

Run: `cd backend && go test -tags=integration ./internal/app/ -v`
Expected: All tests pass.

**Step 6: Write test for asset cleanup**

Add test to `assets_test.go`:
- Upload an asset
- Create a fact referencing it
- Delete the fact
- Verify the asset is also deleted from DB

**Step 7: Run tests**

Run: `cd backend && go test -tags=integration ./internal/app/ -run TestAssetCleanup -v`
Expected: Test passes.

**Step 8: Commit**

```
git add backend/internal/app/facts.go backend/internal/app/assets_test.go
git commit -m "feat: cleanup orphaned assets on fact delete/update"
```

---

## Task 15: Client-Side Image Resize Before Upload

**Files:**
- Create: `frontend/src/lib/utils/image-resize.ts`
- Modify: `frontend/src/lib/api/client.ts` (use resize before upload)

**Step 1: Create image resize utility**

```typescript
// frontend/src/lib/utils/image-resize.ts

const MAX_WIDTH = 2000;
const MAX_HEIGHT = 2000;

export async function resizeImageIfNeeded(file: File): Promise<File> {
    if (!file.type.startsWith('image/') || file.type === 'image/gif') {
        return file; // Don't resize GIFs (animated) or non-images
    }

    return new Promise((resolve) => {
        const img = new window.Image();
        const url = URL.createObjectURL(file);

        img.onload = () => {
            URL.revokeObjectURL(url);

            if (img.width <= MAX_WIDTH && img.height <= MAX_HEIGHT) {
                resolve(file); // No resize needed
                return;
            }

            const ratio = Math.min(MAX_WIDTH / img.width, MAX_HEIGHT / img.height);
            const width = Math.round(img.width * ratio);
            const height = Math.round(img.height * ratio);

            const canvas = document.createElement('canvas');
            canvas.width = width;
            canvas.height = height;

            const ctx = canvas.getContext('2d');
            if (!ctx) {
                resolve(file);
                return;
            }

            ctx.drawImage(img, 0, 0, width, height);

            canvas.toBlob(
                (blob) => {
                    if (!blob) {
                        resolve(file);
                        return;
                    }
                    resolve(new File([blob], file.name, { type: file.type }));
                },
                file.type,
                0.9 // quality for JPEG/WebP
            );
        };

        img.onerror = () => {
            URL.revokeObjectURL(url);
            resolve(file); // On error, use original
        };

        img.src = url;
    });
}
```

**Step 2: Use resize in uploadAsset**

In `frontend/src/lib/api/client.ts`, update `uploadAsset`:

```typescript
import { resizeImageIfNeeded } from '$lib/utils/image-resize';

export async function uploadAsset(file: File): Promise<ApiAsset> {
    const resized = await resizeImageIfNeeded(file);

    const formData = new FormData();
    formData.append('file', resized);

    const res = await fetch('/api/assets', {
        method: 'POST',
        body: formData
    });

    if (!res.ok) {
        const err = await res.json().catch(() => ({ message: 'Upload failed' }));
        throw new Error(err.message ?? 'Upload failed');
    }

    return res.json();
}
```

**Step 3: Verify frontend compiles**

Run: `cd frontend && npm run check`
Expected: No type errors.

**Step 4: Commit**

```
git add frontend/src/lib/utils/image-resize.ts frontend/src/lib/api/client.ts
git commit -m "feat: add client-side image resize before upload for large photos"
```

---

## Summary

| Task | Description | Layer |
|------|-------------|-------|
| 1 | Asset DB migration | Backend DB |
| 2 | SQLC asset queries | Backend DB |
| 3 | Storage interface + local impl | Backend |
| 4 | Wire storage into Application | Backend |
| 5 | Asset upload + serve handlers | Backend API |
| 6 | Asset handler integration tests | Backend Tests |
| 7 | SvelteKit proxy routes + client helper | Frontend |
| 8 | Wire ImageUploader to upload flow | Frontend |
| 9 | Install TipTap dependencies | Frontend |
| 10 | RichTextEditor component | Frontend |
| 11 | Update basic fact editor to rich text | Frontend |
| 12 | RichTextContent read-only renderer | Frontend |
| 13 | Backend validation for rich_text + asset_ids | Backend |
| 14 | Asset cleanup on fact delete/update | Backend |
| 15 | Client-side image resize | Frontend |

Tasks 1-8 unblock image occlusion. Tasks 9-12 deliver rich text for basic facts. Tasks 13-15 are integrity and polish.
