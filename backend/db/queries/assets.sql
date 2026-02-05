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
