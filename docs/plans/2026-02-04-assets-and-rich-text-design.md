# Assets & Rich Text Design

**Date:** 2026-02-04
**Status:** Approved
**Scope:** Asset upload/storage infrastructure, rich text editing for basic facts, image upload for image occlusion editor.

## Motivation

Users need to include images in flashcards: geography maps, flags, X-rays, diagrams, annotated screenshots. This requires two capabilities that don't exist yet:

1. Asset upload and storage infrastructure (no asset table, endpoints, or storage service exist today)
2. Rich text editing for basic fact fields (currently plain text only)
3. Wiring the image occlusion editor's stub `ImageUploader` to a real upload flow

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Upload flow | Proxy through Go backend | Simpler than presigned URLs. Validate server-side. Atomic DB+storage write. Revisit when traffic demands it. |
| Dev storage | Go `Storage` interface with local filesystem impl | No extra containers. ~50 lines of Go. Same interface for R2 in prod. |
| Asset scoping | User-scoped (no notebook FK) | Allows reuse across notebooks (flag sets, anatomy images). |
| Rich text editor | TipTap (minimal: formatting + images) | Already planned in AGENTS.md. Extensible for math/code later. |
| Rich text storage | TipTap JSON in JSONB `value` field | Structured, versionable, round-trips perfectly. |
| Asset ID index | `asset_ids` array at content root | O(1) lookup for references. Avoids traversing ProseMirror doc. |
| Image serving | Stream through Go backend | Consistent with proxy approach. Add `Cache-Control` + `ETag` headers. Future: signed R2 URLs. |
| Image processing | Originals only, client-side resize | No server-side thumbnails. Resize photos > 2000px wide before upload. |
| Deletion strategy | Cleanup on fact delete/update | No user-facing asset management. Orphan check via JSONB containment query. |
| Hash algorithm | SHA-256 | Ecosystem standard (R2/S3). Serves as ETag. Enables future dedup. |

## Asset Storage Infrastructure

### Database: `app.assets`

```sql
CREATE TABLE app.assets (
    user_id     UUID NOT NULL REFERENCES app.users(id),
    id          UUID NOT NULL DEFAULT gen_random_uuid(),
    content_type TEXT NOT NULL,
    filename    TEXT NOT NULL,
    size_bytes  BIGINT NOT NULL,
    sha256      TEXT NOT NULL,            -- hex-encoded SHA-256 of file contents
    metadata    JSONB NOT NULL DEFAULT '{}', -- { width, height, exif, ... }
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, id)
);
```

`metadata` JSONB holds dimensions, EXIF data, and future processing results (OCR text, LLM analysis) without schema migrations.

### Go Storage Interface

```go
type Storage interface {
    Put(ctx context.Context, key string, reader io.Reader, contentType string) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
}
```

Two implementations:
- `LocalStorage` - writes to configurable directory (e.g., `./data/assets/`), serves via static file handler
- `R2Storage` - S3 SDK against Cloudflare R2

Selected by environment config.

### Storage Key Format

`assets/{user_id}/{asset_id}` - flat, no nesting by notebook.

### Endpoints

**Upload:** `POST /v1/assets`
- Accepts `multipart/form-data`
- Validates: file type (JPEG, PNG, WebP, GIF), max size (10MB)
- Computes SHA-256 during streaming read
- Extracts image dimensions, stores in `metadata`
- Writes to storage, creates DB record
- Returns asset record with ID

**Serve:** `GET /v1/assets/{id}/file`
- Streams file from storage
- Headers: `Cache-Control: private, max-age=31536000, immutable`
- `ETag` based on SHA-256 hash
- Supports `If-None-Match` for 304 responses

## Fact Content Schema

### Basic fact with rich text

```json
{
  "version": 1,
  "asset_ids": ["019...abc"],
  "fields": [
    {
      "name": "front",
      "type": "rich_text",
      "value": {
        "type": "doc",
        "content": [
          { "type": "paragraph", "content": [{ "type": "text", "text": "What flag is this?" }] },
          { "type": "image", "attrs": { "asset_id": "019...abc", "width": 400, "height": 300 } }
        ]
      }
    },
    {
      "name": "back",
      "type": "plain_text",
      "value": "France"
    }
  ]
}
```

### Key rules

- `plain_text` fields: `value` is a string. No changes to existing facts.
- `rich_text` fields: `value` is a TipTap JSON document object.
- `asset_ids` at root: denormalized index of all asset IDs referenced in the fact. Frontend keeps in sync on save.
- Backend validates: `asset_ids` match image nodes in doc, referenced assets exist and belong to user.
- Image occlusion facts: `asset_ids` populated with the occlusion image's asset ID.
- Existing plain text facts untouched. No content migration needed.

### GIN index for orphan queries

```sql
CREATE INDEX idx_facts_asset_ids ON app.facts USING GIN ((content->'asset_ids'));
```

## Frontend: Rich Text Editor

### RichTextEditor.svelte

TipTap wrapper with minimal extensions:
- `StarterKit` (bold, italic, underline, lists, headings)
- Custom `Image` extension - renders by asset ID, resolves URL to `/api/assets/{id}/file`
- Minimal floating toolbar

Accepts and emits TipTap JSON. Used only for `rich_text` field types.

### Image insertion flow

1. User clicks image button, pastes, or drops an image
2. Client-side resize if > 2000px wide (for phone camera photos)
3. Upload via `POST /api/assets` (SvelteKit proxy)
4. On success, insert TipTap image node with `asset_id` attribute
5. On save, collect all `asset_id` values into root `asset_ids` array

### RichTextContent.svelte

Read-only renderer for card review. Renders TipTap JSON without editor overhead. Same image URL resolution.

### ImageUploader.svelte (image occlusion)

Existing stub wired to the same `uploadAsset()` helper. Returns asset ID for the occlusion data's `image.assetId` field.

## SvelteKit Proxy Layer

### Routes

- `POST /api/assets` - proxies multipart upload to Go `POST /v1/assets`
- `GET /api/assets/[id]/file` - proxies to Go `GET /v1/assets/{id}/file`, passes through cache headers

### Client helper

`uploadAsset(file: File): Promise<Asset>` in `$lib/api/client.ts`. Constructs `FormData`, calls proxy route. Used by both `RichTextEditor` and `ImageUploader`.

## Asset Deletion

### On fact delete

1. Read fact's `asset_ids` from content
2. For each asset ID, check if any other fact by this user references it: `content->'asset_ids' @> '["019...abc"]'::jsonb`
3. Delete unreferenced assets from storage + DB
4. DB deletes in transaction, storage deletes after commit

### On fact update

- Compare old `asset_ids` with new
- Removed IDs go through same orphan check

No background jobs. No user-facing asset management.

## Out of Scope

- Presigned URLs / direct-to-R2 uploads
- Server-side image thumbnails/resizing
- Image processing pipeline (OCR, LLM analysis) - asset model supports it, not building it
- Math/KaTeX TipTap extension - added later, no architectural impact
- Audio/video assets - same pipeline, different content type validation
- User-facing asset management UI
- Asset deduplication (hash enables it later)

## Implementation Order

1. Asset DB migration + Go storage interface + local filesystem implementation
2. Upload endpoint (`POST /v1/assets`) + serve endpoint (`GET /v1/assets/{id}/file`)
3. SvelteKit proxy routes + `uploadAsset()` client helper
4. Wire `ImageUploader.svelte` in image occlusion editor to real upload flow
5. TipTap `RichTextEditor.svelte` component with image extension
6. Update `basic-fact-editor.svelte` to use rich text editor for front/back fields
7. `RichTextContent.svelte` for card review rendering
8. Backend validation updates for `rich_text` fields + `asset_ids`
9. Asset cleanup on fact delete/update
10. Client-side image resize before upload

Steps 1-4 unblock the image occlusion editor. Steps 5-7 deliver rich text for basic facts. Steps 8-10 are integrity and polish.
